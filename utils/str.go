package utils

import (
	"crypto/rand"
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
