package utils

import (
	"crypto/rand"

	"github.com/goccy/go-json"
	"github.com/rs/zerolog/log"
)

func RandStr(len int) string {
	buffer := make([]byte, len)
	_, err := rand.Read(buffer)
	if err != nil {
		log.Panic().Err(err).Msg("failed to generate random string")
	}

	return string(buffer)
}

func EncodeElasticSearchID(str string) string {
	return Fnv1a32SumBytes([]byte(str))
}

func JsonRaw(data interface{}) json.RawMessage {
	raw, err := json.Marshal(data)
	if err != nil {
		log.Panic().Err(err).Msg("failed to marshal json")
	}
	return raw
}
