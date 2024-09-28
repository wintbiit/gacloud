package utils

import "github.com/goccy/go-json"

func ToJsonRaw(v interface{}) string {
	b, _ := json.Marshal(v)
	return string(b)
}
