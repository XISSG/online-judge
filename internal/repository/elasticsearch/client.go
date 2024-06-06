package elasticsearch

import (
	elasticsearch "github.com/elastic/go-elasticsearch/v8"
	"github.com/xissg/online-judge/internal/config"
)

type ElasticsearchClient struct {
	Client *elasticsearch.TypedClient
}

func NewElasticSearchClient(cfg config.ElasticsearchConfig) *ElasticsearchClient {
	esCfg := elasticsearch.Config{
		Addresses: cfg.Addresses,
		Username:  cfg.Username,
		Password:  cfg.Password,
	}
	client, err := elasticsearch.NewTypedClient(esCfg)
	if err != nil {
		panic(err)
	}

	return &ElasticsearchClient{
		Client: client,
	}
}
