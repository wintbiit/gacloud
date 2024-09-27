package server

import (
	"context"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/scriptlanguage"
	"github.com/rs/zerolog"
	"github.com/wintbiit/gacloud/config"
	"github.com/wintbiit/gacloud/model"
	"github.com/wintbiit/gacloud/utils"
	"gorm.io/gorm"
)

type GaCloudServer struct {
	db          *gorm.DB
	es          *elasticsearch.TypedClient
	logger      *zerolog.Logger
	Info        *model.AppInfo
	Maintenance bool
}

const elasticSearchIndex = "gacloud.files.v1"

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

	es, err := utils.OpenElasticSearch(esHost, esUser, esPassword, elasticSearchIndex)
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

	serverLogger := utils.NewLogger("server")
	serverInfo := &model.AppInfo{
		SiteName:    config.GetWithDefault("site.name", "GaCloud"),
		ExternalUrl: config.GetWithDefault("site.external_url", "http://localhost:8080"),
		SiteLogo:    config.GetWithDefault("site.logo", ""),
	}

	return &GaCloudServer{
		db:     db,
		es:     es,
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
