package utils

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/elastic/go-elasticsearch/v8"
)

func OpenElasticSearch(host, user, password string) (*elasticsearch.TypedClient, error) {
	cfg := elasticsearch.Config{
		Addresses: []string{
			host,
		},
		Username: user,
		Password: password,
	}

	es, err := elasticsearch.NewTypedClient(cfg)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	info, err := es.Info().Do(ctx)
	if err != nil {
		return nil, err
	}

	log.Info().Interface("info", info).Msg("Elasticsearch client created")

	return es, nil
}
