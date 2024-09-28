package server

import (
	"context"
	"path"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/scriptlanguage"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/wintbiit/gacloud/config"
	"github.com/wintbiit/gacloud/fs"
	"github.com/wintbiit/gacloud/model"
	"github.com/wintbiit/gacloud/utils"
	"gorm.io/gorm"
)

type GaCloudServer struct {
	db            *gorm.DB
	es            *elasticsearch.TypedClient
	esIndex       string
	fileProviders map[uint]fs.FileProvider
	logger        *zerolog.Logger
	Info          *model.AppInfo
	Maintenance   bool
}

const elasticSearchIndex = "gacloud.files.v1"

func init() {
	utils.AddShutdownHook(func() {
		if s != nil {
			s.Close()
		}
	})
}

func NewLocalGaCloudServer() (*GaCloudServer, error) {
	db, err := setupDb()
	if err != nil {
		return nil, err
	}

	es, err := setupElasticSearch()
	if err != nil {
		return nil, err
	}

	providers, err := setupFileProviders(db)
	if err != nil {
		return nil, err
	}

	serverLogger := utils.NewLogger("server")
	serverInfo := &model.AppInfo{
		SiteName:    config.GetWithDefault("site.name", "GaCloud"),
		ExternalUrl: config.GetWithDefault("site.external_url", "http://localhost:8080"),
		SiteLogo:    config.GetWithDefault("site.logo", ""),
	}

	return &GaCloudServer{
		db:            db,
		es:            es,
		esIndex:       elasticSearchIndex,
		fileProviders: providers,
		logger:        &serverLogger,
		Info:          serverInfo,
	}, nil
}

func (s *GaCloudServer) Close() {
}

var s *GaCloudServer

func GetServer() *GaCloudServer {
	if s != nil {
		return s
	}

	s, err := NewLocalGaCloudServer()
	if err != nil {
		return nil
	}

	return s
}

func setupDb() (*gorm.DB, error) {
	dbType, ok := config.Get("db.type")
	if !ok {
		return nil, utils.ErrorSetupNotCompleted
	}

	dbDsn, ok := config.Get("db.dsn")
	if !ok {
		return nil, utils.ErrorSetupNotCompleted
	}

	db, err := utils.OpenDB(dbType, dbDsn)
	if err != nil {
		return nil, err
	}

	logger := utils.NewGormLogger()
	db.Logger = logger

	err = model.MigrateModels(db)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func setupElasticSearch() (*elasticsearch.TypedClient, error) {
	esHost, ok := config.Get("es.host")
	if !ok {
		return nil, utils.ErrorSetupNotCompleted
	}

	esUser, ok := config.Get("es.user")
	if !ok {
		return nil, utils.ErrorSetupNotCompleted
	}

	esPassword, ok := config.Get("es.password")
	if !ok {
		return nil, utils.ErrorSetupNotCompleted
	}

	es, err := utils.OpenElasticSearch(esHost, esUser, esPassword)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := es.PutScript(listFileScriptId).Script(&types.StoredScript{
		Lang: scriptlanguage.ScriptLanguage{
			Name: "painless",
		},
		Source: listFileScript,
	}).Do(ctx)
	if err != nil {
		return nil, err
	}

	if !resp.Acknowledged {
		return nil, utils.ErrorElasticSearchScriptNotAcknowledged
	}

	resp, err = es.PutScript(permissionScriptId).Script(&types.StoredScript{
		Lang: scriptlanguage.ScriptLanguage{
			Name: "painless",
		},
		Source: permissionScript,
	}).Do(ctx)
	if err != nil {
		return nil, err
	}

	if !resp.Acknowledged {
		return nil, utils.ErrorElasticSearchScriptNotAcknowledged
	}

	return es, nil
}

func setupFileProviders(db *gorm.DB) (map[uint]fs.FileProvider, error) {
	providers := make(map[uint]fs.FileProvider)

	var fileProviders []model.FileProvider
	err := db.Find(&fileProviders).Error
	if err != nil {
		return nil, err
	}

	providers[0] = fs.NewLocalFileProviderWithConfig(fs.LocalFileProviderConfig{
		MountDir: path.Join(utils.ServerInfo.DataDir, "files"),
	})

	for _, provider := range fileProviders {
		factory := fs.GetFileProviderFactory(provider.Type)
		if factory == nil {
			log.Error().Str("provider", provider.Name).Msg("file provider not found")
			return nil, utils.ErrorFileProviderNotFound
		}

		fp, err := factory([]byte(provider.Credential))
		if err != nil {
			log.Error().Str("provider", provider.Name).Err(err).Msg("failed to create file provider")
			return nil, err
		}

		providers[provider.ID] = fp
	}

	return providers, nil
}

func (s *GaCloudServer) RefreshFileProvider() error {
	providers, err := setupFileProviders(s.db)
	if err != nil {
		return err
	}

	s.fileProviders = providers
	return nil
}
