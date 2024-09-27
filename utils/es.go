package utils

import (
	"context"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/rs/zerolog/log"
)

func OpenElasticSearch(host, user, password, indexName string) (*elasticsearch.TypedClient, error) {
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

	exist, err := es.Indices.Exists(indexName).Do(ctx)
	if err != nil {
		return nil, err
	}

	if !exist {
		log.Info().Str("index", indexName).Str("framework", "elasticsearch").Msg("Index not found, creating new index")
		_, err = es.Indices.Create(indexName).Do(ctx)
		if err != nil {
			return nil, err
		}
	} else {
		log.Info().Str("index", indexName).Str("framework", "elasticsearch").Msg("Index found")
	}

	return es, nil
}
