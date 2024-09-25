package utils

import "fmt"

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
