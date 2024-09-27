package server

import (
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/rs/zerolog"
	"github.com/wintbiit/gacloud/config"
	"github.com/wintbiit/gacloud/model"
	"github.com/wintbiit/gacloud/utils"
	"gorm.io/gorm"
)

type GaCloudServer struct {
	db          *gorm.DB
	es          *elasticsearch.Client
	logger      *zerolog.Logger
	Info        *model.AppInfo
	Maintenance bool
}

func NewLocalGaCloudServer() (*GaCloudServer, error) {
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

	serverLogger := utils.NewLogger("server")
	serverInfo := &model.AppInfo{
		SiteName:    config.GetWithDefault("site.name", "GaCloud"),
		ExternalUrl: config.GetWithDefault("site.external_url", "http://localhost:8080"),
		SiteLogo:    config.GetWithDefault("site.logo", ""),
	}

	return &GaCloudServer{
		db:     db,
		logger: &serverLogger,
		Info:   serverInfo,
	}, nil
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
