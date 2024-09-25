package config

import (
	"context"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/zerolog/log"
	"github.com/wintbiit/gacloud/model"
	"strconv"
	"time"
	"xorm.io/xorm"
	"xorm.io/xorm/caches"
)

type config struct {
	Id    int64  `xorm:"pk autoincr"`
	Key   string `xorm:"unique,index"`
	Value string `xorm:"text"`
	model.TimeModel
}

var engine *xorm.Engine

func init() {
	var err error
	engine, err = xorm.NewEngine("sqlite3", "./gacloud.config.db")
	if err != nil {
		log.Panic().Err(err).Msg("failed to load config database. Please check ./gacloud.config.db")
	}

	cache := caches.NewLRUCacher2(caches.NewMemoryStore(), 30*time.Second, 1000)
	engine.SetDefaultCacher(cache)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err = engine.PingContext(ctx); err != nil {
		log.Panic().Err(err).Msg("failed to connect to config database")
	}

	if err = engine.Sync2(new(config)); err != nil {
		log.Panic().Err(err).Msg("failed to sync config database")
	}

	log.Info().Msg("config database is ready")
}

func Get(key string) (string, bool) {
	var c config
	has, err := engine.Where("key = ?", key).Get(&c)
	if err != nil {
		log.Error().Err(err).Str("key", key).Msg("failed to get config")
	}

	return c.Value, has
}

func GetWithDefault(key, defaultValue string) string {
	value, has := Get(key)
	if !has {
		return defaultValue
	}

	return value
}

func Set(key, value string) {
	c := config{
		Key:   key,
		Value: value,
	}

	_, err := engine.Insert(c)
	if err != nil {
		log.Error().Err(err).Str("key", key).Str("value", value).Msg("failed to set config")
	}
}

func Delete(key string) {
	_, err := engine.Where("key = ?", key).Delete(new(config))
	if err != nil {
		log.Error().Err(err).Str("key", key).Msg("failed to delete config")
	}
}

func GetInt(key string) (int, bool) {
	value, has := Get(key)
	if !has {
		return 0, false
	}

	parse, err := strconv.Atoi(value)
	if err != nil {
		log.Error().Err(err).Str("key", key).Str("value", value).Msg("failed to parse int")
		return 0, false
	}

	return parse, true
}

func GetIntWithDefault(key string, defaultValue int) int {
	value, has := GetInt(key)
	if !has {
		return defaultValue
	}

	return value
}

func SetInt(key string, value int) {
	Set(key, strconv.Itoa(value))
}

func GetBool(key string) (bool, bool) {
	value, has := Get(key)
	if !has {
		return false, false
	}

	parse, err := strconv.ParseBool(value)
	if err != nil {
		log.Error().Err(err).Str("key", key).Str("value", value).Msg("failed to parse bool")
		return false, false
	}

	return parse, true
}

func GetBoolWithDefault(key string, defaultValue bool) bool {
	value, has := GetBool(key)
	if !has {
		return defaultValue
	}

	return value
}

func SetBool(key string, value bool) {
	Set(key, strconv.FormatBool(value))
}

func MustGet(key string) string {
	value, has := Get(key)
	if !has {
		log.Panic().Str("key", key).Msg("missing config")
	}

	return value
}

func GetAll(prefix string) map[string]string {
	var cs []config
	err := engine.Where("key like ?", prefix+"%").Find(&cs)
	if err != nil {
		log.Error().Err(err).Str("prefix", prefix).Msg("failed to get all config")
	}

	m := make(map[string]string)
	for _, c := range cs {
		m[c.Key] = c.Value
	}

	return m
}
