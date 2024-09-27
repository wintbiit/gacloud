package config

import (
	"context"
	"os"
	"path"
	"strconv"
	"sync"
	"time"

	"github.com/glebarez/sqlite"
	"github.com/rs/zerolog/log"
	"github.com/wintbiit/gacloud/utils"
	"gorm.io/gorm"
)

type config struct {
	Key   string `gorm:"primaryKey,unique,not null,index"`
	Value string `gorm:"not null"`
	gorm.Model
}

var (
	engine  *gorm.DB
	cfgPool = sync.Pool{
		New: func() interface{} {
			return new(config)
		},
	}
)

func init() {
	var err error
	dir := path.Join(utils.ServerInfo.DataDir, "config")
	os.MkdirAll(dir, 0o755)
	engine, err = gorm.Open(sqlite.Open(path.Join(dir, "gacloud.config.db")), &gorm.Config{})
	if err != nil {
		log.Panic().Err(err).Msg("failed to load config database. Please check ./gacloud.config.db")
	}
	engine.Logger = utils.NewGormLogger()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = engine.WithContext(ctx).AutoMigrate(&config{})
	if err != nil {
		log.Panic().Err(err).Msg("failed to migrate config database")
	}

	log.Info().Msg("config database is ready")
}

func Get(key string) (string, bool) {
	c := cfgPool.Get().(*config)
	defer func() {
		c.ID = 0
		cfgPool.Put(c)
	}()

	engine.Where("key = ?", key).First(c)
	if c.ID == 0 {
		return "", false
	}

	return c.Value, true
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

	// insert or update
	err := engine.Where("key = ?", key).Assign(c).FirstOrCreate(&c).Error
	if err != nil {
		log.Error().Err(err).Str("key", key).Str("value", value).Msg("failed to set config")
	}
}

func Delete(key string) {
	err := engine.Where("key = ?", key).Delete(&config{}).Error
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
	rows, err := engine.Where("key like ?", prefix+"%").Rows()
	if err != nil {
		log.Error().Err(err).Str("prefix", prefix).Msg("failed to get all config")
	}

	defer rows.Close()
	cs := make([]config, 0)
	for rows.Next() {
		c := cfgPool.Get().(*config)
		err = engine.ScanRows(rows, c)
		if err != nil {
			log.Error().Err(err).Msg("failed to scan config")
			continue
		}

		cs = append(cs, *c)
	}

	m := make(map[string]string)
	for _, c := range cs {
		m[c.Key] = c.Value
	}

	return m
}
