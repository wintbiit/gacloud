package utils

import (
	"fmt"

	"github.com/glebarez/sqlite"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func OpenDB(dbType, dsn string) (*gorm.DB, error) {
	switch dbType {
	case "mysql":
		return gorm.Open(mysql.Open(dsn), &gorm.Config{})
	case "postgres":
		return gorm.Open(postgres.Open(dsn), &gorm.Config{})
	case "sqlite":
		return gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	default:
		return nil, fmt.Errorf("unsupported db type")
	}
}

func ParseDb(dbType string, params map[string]interface{}) string {
	switch dbType {
	case "mysql":
		return parseMysql(params)
	case "postgres":
		return parsePostgres(params)
	case "sqlite":
		return parseSqlite(params)
	default:
		return ""
	}
}

func parseMysql(params map[string]interface{}) string {
	host, ok := params["host"]
	if !ok {
		return ""
	}

	port, ok := params["port"]
	if !ok {
		return ""
	}

	user, ok := params["user"]
	if !ok {
		return ""
	}

	password, ok := params["password"]
	if !ok {
		return ""
	}

	dbname, ok := params["dbname"]
	if !ok {
		return ""
	}

	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, password, host, port, dbname)
}

func parsePostgres(params map[string]interface{}) string {
	host, ok := params["host"]
	if !ok {
		return ""
	}

	port, ok := params["port"]
	if !ok {
		return ""
	}

	user, ok := params["user"]
	if !ok {
		return ""
	}

	password, ok := params["password"]
	if !ok {
		return ""
	}

	dbname, ok := params["dbname"]
	if !ok {
		return ""
	}

	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
}

func parseSqlite(params map[string]interface{}) string {
	path, ok := params["path"]
	if !ok {
		return ""
	}

	pathStr, ok := path.(string)
	if !ok {
		return ""
	}

	return pathStr
}
