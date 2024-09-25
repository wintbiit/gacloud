package internal

import (
	"context"
	"github.com/rs/zerolog/log"
	"github.com/wintbiit/gacloud/config"
	"github.com/wintbiit/gacloud/utils"
	"time"
	"xorm.io/xorm"
)

func init() {
	utils.AddShutdownHook(func() {
		if server != nil {
			server.Close()
			log.Info().Msg("Server closed")
		}
	})
}

var engine *xorm.Engine

func getEngine() *xorm.Engine {
	if engine != nil {
		return engine
	}

	dbType, ok := config.Get("db.type")
	if !ok {
		return nil
	}

	dbDsn, ok := config.Get("db.dsn")
	if !ok {
		return nil
	}

	var err error
	engine, err = xorm.NewEngine(dbType, dbDsn)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create database engine")
		engine = nil
		return nil
	}
	engine.SetLogger(utils.NewXormLogger())

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = engine.PingContext(ctx)
	if err == nil {
		return engine
	}

	engine = nil
	log.Error().Err(err).Msg("Failed to connect to database")
	return nil
}

func TryDbOptions(engineType string, params map[string]interface{}) error {
	dbDsn := utils.ParseDb(engineType, params)
	if dbDsn == "" {
		return nil
	}

	engine, err := xorm.NewEngine(engineType, dbDsn)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create database engine")
		return err
	}
	defer engine.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = engine.PingContext(ctx)
	if err != nil {
		return err

	}

	config.Set("db.type", engineType)
	config.Set("db.dsn", dbDsn)
	return nil
}

type GaCloudServer struct {
	Db *xorm.Engine
}

var server *GaCloudServer

func GetGaCloudServer() *GaCloudServer {
	if server != nil {
		return server
	}

	engine := getEngine()
	if engine == nil {
		return nil
	}

	server = &GaCloudServer{
		Db: engine,
	}

	return server
}

func (s *GaCloudServer) HealthCheck() bool {
	if s.Db == nil {
		return false
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := s.Db.PingContext(ctx)
	return err == nil
}

func (s *GaCloudServer) Close() {
	if s.Db != nil {
		s.Db.Close()
	}
}
