package routes

import (
	"github.com/kataras/iris/v12"
	"github.com/wintbiit/gacloud/config"
	"github.com/wintbiit/gacloud/fs"
	"github.com/wintbiit/gacloud/model"
	"github.com/wintbiit/gacloud/utils"
	"gorm.io/gorm"
)

func init() {
	addHookFront("/api/v1/setup", func(app iris.Party) {
		app.Use(SetupPreroute)
		app.Get("/", GetSetupStatus)
		app.Post("/database", SetDataBaseOptions)
		app.Post("/database/test", TestDataBase)
		app.Post("/elasticsearch", SetElasticSearchOptions)
		app.Post("/elasticsearch/test", TestElasticSearch)
		app.Get("/storage/providers", GetStorageProviders)
		app.Post("/storage", SetStorageProvider)
		app.Post("/storage/test", TestStorageProvider)
		app.Post("/admin", SetSuperAdmin)
		app.Post("/finish", FinishSetup)
	})
}

// getSetupStatus returns the current setup status
// 0: Database not configured
// 1: ElasticSearch not configured
// 2: Initial storage not configured
// 3: Admin user not configured
// 4: Setup completed
func getSetupStatus() int {
	if !isDbConfigured() {
		return 0
	}

	if !isElasticSearchConfigured() {
		return 1
	}

	if !isStorageConfigured() {
		return 2
	}

	if !isSuperAdminConfigured() {
		return 3
	}

	return 4
}

func isDbConfigured() bool {
	_, ok := config.Get("db.type")
	if !ok {
		return false
	}

	_, ok = config.Get("db.dsn")
	if !ok {
		return false
	}

	return true
}

func isElasticSearchConfigured() bool {
	_, ok := config.Get("es.host")
	if !ok {
		return false
	}

	_, ok = config.Get("es.user")
	if !ok {
		return false
	}

	_, ok = config.Get("es.password")
	if !ok {
		return false
	}

	return true
}

func isStorageConfigured() bool {
	_, ok := config.Get("storage0.type")
	if !ok {
		return false
	}

	_, ok = config.Get("storage0.credential")
	if !ok {
		return false
	}

	return true
}

func isSuperAdminConfigured() bool {
	dbType, _ := config.Get("db.type")
	dbDsn, _ := config.Get("db.dsn")

	db, err := utils.OpenDB(dbType, dbDsn)
	if err != nil {
		return false
	}

	// 1. Check if admin group (gid = 1) exists
	var adminGroup model.Group
	if err := db.Where("id = 1").First(&adminGroup).Error; err != nil {
		return false
	}

	// 2. Check if admin user (uid = 1) exists
	var adminUser model.User
	if err := db.Where("id = 1").First(&adminUser).Error; err != nil {
		return false
	}

	// 3. Check if admin user is in admin group
	var count int64
	if err := db.Table("user_groups").Where("user_id = ? AND group_id = ?", adminUser.ID, adminGroup.ID).Count(&count).Error; err != nil {
		return false
	}

	return count > 0
}

func SetupPreroute(ctx iris.Context) {
	if config.GetBoolWithDefault("setup", false) {
		ctx.StopWithStatus(iris.StatusNotFound)
		return
	}

	if getSetupStatus() > 4 {
		ctx.StopWithStatus(iris.StatusNotFound)
		return
	}

	ctx.Next()
}

type SetupStatusResponse struct {
	CurrentStep int `json:"currentStep"`
}

func GetSetupStatus(ctx iris.Context) {
	ctx.JSON(SetupStatusResponse{
		CurrentStep: getSetupStatus(),
	})
}

type DataBaseOption struct {
	Type   string                 `json:"type"`
	Params map[string]interface{} `json:"params"`
}

func SetDataBaseOptions(ctx iris.Context) {
	var dbOption DataBaseOption
	if err := ctx.ReadJSON(&dbOption); err != nil {
		ctx.StopWithError(iris.StatusBadRequest, err)
		return
	}

	dsn := utils.ParseDb(dbOption.Type, dbOption.Params)

	config.Set("db.type", dbOption.Type)
	config.Set("db.dsn", dsn)

	ctx.StatusCode(iris.StatusOK)
}

type TestResponse struct {
	Success bool   `json:"success"`
	Reason  string `json:"reason"`
}

func TestDataBase(ctx iris.Context) {
	var dbOption DataBaseOption
	if err := ctx.ReadJSON(&dbOption); err != nil {
		ctx.JSON(TestResponse{
			Success: false,
			Reason:  err.Error(),
		})
		return
	}

	db, err := utils.OpenDB(dbOption.Type, utils.ParseDb(dbOption.Type, dbOption.Params))
	if err != nil {
		ctx.JSON(TestResponse{
			Success: false,
			Reason:  err.Error(),
		})
		return
	}

	if err := db.Exec("SELECT 1").Error; err != nil {
		ctx.JSON(TestResponse{
			Success: false,
			Reason:  err.Error(),
		})
		return
	}

	ctx.JSON(TestResponse{
		Success: true,
	})
}

type ElasticSearchOption struct {
	Host     string `json:"host"`
	User     string `json:"user"`
	Password string `json:"password"`
}

func SetElasticSearchOptions(ctx iris.Context) {
	var esOption ElasticSearchOption
	if err := ctx.ReadJSON(&esOption); err != nil {
		ctx.StopWithError(iris.StatusBadRequest, err)
		return
	}

	config.Set("es.host", esOption.Host)
	config.Set("es.user", esOption.User)
	config.Set("es.password", esOption.Password)

	ctx.StatusCode(iris.StatusOK)
}

func TestElasticSearch(ctx iris.Context) {
	var esOption ElasticSearchOption
	if err := ctx.ReadJSON(&esOption); err != nil {
		ctx.JSON(TestResponse{
			Success: false,
			Reason:  err.Error(),
		})
		return
	}

	es, err := utils.OpenElasticSearch(esOption.Host, esOption.User, esOption.Password)
	if err != nil {
		ctx.JSON(TestResponse{
			Success: false,
			Reason:  err.Error(),
		})
		return
	}

	if _, err := es.Info().Do(ctx); err != nil {
		ctx.JSON(TestResponse{
			Success: false,
			Reason:  err.Error(),
		})
		return
	}

	ctx.JSON(TestResponse{
		Success: true,
	})
}

type StorageProviderConfig struct {
	RootPath  string   `json:"rootPath"`
	Providers []string `json:"providers"`
}

func GetStorageProviders(ctx iris.Context) {
	ctx.JSON(StorageProviderConfig{
		RootPath:  utils.ServerInfo.DataDir,
		Providers: fs.ListFileProviderTypes(),
	})
}

func SetStorageProvider(ctx iris.Context) {
	var provider model.FileProvider
	if err := ctx.ReadJSON(&provider); err != nil {
		ctx.StopWithError(iris.StatusBadRequest, err)
		return
	}

	factory := fs.GetFileProviderFactory(provider.Type)
	if factory == nil {
		ctx.StopWithStatus(iris.StatusBadRequest)
		return
	}

	_, err := factory([]byte(provider.Credential))
	if err != nil {
		ctx.StopWithError(iris.StatusBadRequest, err)
		return
	}

	config.Set("storage0.name", provider.Name)
	config.Set("storage0.type", provider.Type)
	config.Set("storage0.credential", provider.Credential)

	ctx.StatusCode(iris.StatusOK)
}

func TestStorageProvider(ctx iris.Context) {
	var provider model.FileProvider
	if err := ctx.ReadJSON(&provider); err != nil {
		ctx.JSON(TestResponse{
			Success: false,
			Reason:  err.Error(),
		})
		return
	}

	factory := fs.GetFileProviderFactory(provider.Type)
	if factory == nil {
		ctx.JSON(TestResponse{
			Success: false,
			Reason:  "Invalid provider type",
		})
		return
	}

	_, err := factory([]byte(provider.Credential))
	if err != nil {
		ctx.JSON(TestResponse{
			Success: false,
			Reason:  err.Error(),
		})
		return
	}

	ctx.JSON(TestResponse{
		Success: true,
	})
}

type SetSuperAdminRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

func SetSuperAdmin(ctx iris.Context) {
	var request SetSuperAdminRequest
	if err := ctx.ReadJSON(&request); err != nil {
		ctx.StopWithError(iris.StatusBadRequest, err)
		return
	}

	dbType, _ := config.Get("db.type")
	dbDsn, _ := config.Get("db.dsn")

	db, err := utils.OpenDB(dbType, dbDsn)
	if err != nil {
		ctx.StopWithError(iris.StatusInternalServerError, err)
		return
	}

	// 1. Create admin group (gid = 1)
	adminGroup := &model.Group{
		Model: gorm.Model{
			ID: 1,
		},
		Name: "admin",
	}

	if err := db.WithContext(ctx).Create(adminGroup).Error; err != nil {
		ctx.StopWithError(iris.StatusInternalServerError, err)
		return
	}

	// 2. Create admin user (uid = 1)
	password := utils.Sha256SumBytes([]byte(request.Password))
	adminUser := &model.User{
		Model: gorm.Model{
			ID: 1,
		},
		Name:     request.Username,
		Password: password,
		Email:    request.Email,
	}

	if err := db.WithContext(ctx).Create(adminUser).Error; err != nil {
		ctx.StopWithError(iris.StatusInternalServerError, err)
		return
	}

	// 3. Add admin user to admin group
	ug := &model.UserGroup{
		User:  *adminUser,
		Group: *adminGroup,
	}

	if err := db.WithContext(ctx).Create(ug).Error; err != nil {
		ctx.StopWithError(iris.StatusInternalServerError, err)
		return
	}

	ctx.StatusCode(iris.StatusOK)
}

func FinishSetup(ctx iris.Context) {
	config.SetBool("setup", true)

	ctx.StatusCode(iris.StatusOK)
}
