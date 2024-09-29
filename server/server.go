package server

import (
	"context"
	"github.com/kataras/iris/v12/middleware/jwt"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
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
	signer        *jwt.Signer
	verifier      *jwt.Verifier
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
	setupCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	db, err := setupDb(setupCtx)
	if err != nil {
		return nil, err
	}

	es, err := setupElasticSearch(setupCtx)
	if err != nil {
		return nil, err
	}

	providers, err := setupFileProviders(setupCtx, db)
	if err != nil {
		return nil, err
	}

	serverLogger := utils.NewLogger("server")
	serverInfo := &model.AppInfo{
		SiteName:    config.GetWithDefault("site.name", "GaCloud"),
		ExternalUrl: config.GetWithDefault("site.external_url", "http://localhost:8080"),
		SiteLogo:    config.GetWithDefault("site.logo", ""),
	}

	jwtSecret, ok := config.Get("jwt.secret")
	if !ok {
		jwtSecret = utils.RandStr(32)
		config.Set("jwt.secret", jwtSecret)
	}

	jwtExpiration, ok := config.Get("jwt.expiration")
	if !ok {
		jwtExpiration = "1h"
		config.Set("jwt.expiration", jwtExpiration)
	}

	jwtExpirationDuration, err := time.ParseDuration(jwtExpiration)
	if err != nil {
		return nil, err
	}

	signer := jwt.NewSigner(jwt.HS256, []byte(jwtSecret), jwtExpirationDuration)
	verifier := jwt.NewVerifier(jwt.HS256, []byte(jwtSecret))

	return &GaCloudServer{
		db:            db,
		es:            es,
		esIndex:       elasticSearchIndex,
		fileProviders: providers,
		signer:        signer,
		verifier:      verifier,
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
		log.Error().Err(err).Msg("failed to setup server")
		return nil
	}

	return s
}

func setupDb(ctx context.Context) (*gorm.DB, error) {
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

	err = model.MigrateModels(ctx, db)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func setupElasticSearch(ctx context.Context) (*elasticsearch.TypedClient, error) {
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

	exists, err := es.Indices.Exists(elasticSearchIndex).Do(ctx)
	if err != nil {
		return nil, err
	}

	if !exists {
		_, err = es.Indices.Create(elasticSearchIndex).Mappings(&types.TypeMapping{
			Properties: map[string]types.Property{
				"sum":         types.NewTextProperty(),
				"path":        types.NewKeywordProperty(),
				"size":        types.NewIntegerNumberProperty(),
				"mime":        types.NewTextProperty(),
				"owner_type":  types.NewIntegerNumberProperty(),
				"owner_id":    types.NewIntegerNumberProperty(),
				"provider_id": types.NewIntegerNumberProperty(),
			},
		}).Do(ctx)
		if err != nil {
			return nil, err
		}
	}

	return es, nil
}

func setupFileProviders(ctx context.Context, db *gorm.DB) (map[uint]fs.FileProvider, error) {
	providers := make(map[uint]fs.FileProvider)

	var fileProviders []model.FileProvider

	err := db.WithContext(ctx).Find(&fileProviders).Error
	if err != nil {
		return nil, err
	}

	storage0Type, ok := config.Get("storage0.type")
	if !ok {
		return nil, utils.ErrorSetupNotCompleted
	}

	storage0Credential, ok := config.Get("storage0.credential")
	if !ok {
		return nil, utils.ErrorSetupNotCompleted
	}

	storage0Factory := fs.GetFileProviderFactory(storage0Type)
	if storage0Factory == nil {
		log.Error().Str("provider", storage0Type).Msg("file provider not found")
		return nil, utils.ErrorFileProviderNotFound
	}

	storage0, err := storage0Factory([]byte(storage0Credential))
	if err != nil {
		log.Error().Str("provider", storage0Type).Err(err).Msg("failed to create file provider")
		return nil, err
	}

	providers[0] = storage0

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
