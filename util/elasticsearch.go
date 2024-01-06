package util

import (
	elastic "github.com/olivere/elastic/v7"
	"github.com/pkg/errors"
)

var elasticClient *elastic.Client

func ConnectToElastic(config *Config) (err error) {
	elasticClient, err = elastic.NewClient(elastic.SetURL(config.ElasticSearch.URL), elastic.SetSniff(false))
	if err != nil {
		return errors.Wrap(err, "Error connecting to Elasticsearch")
	}

	return nil
}

func GetElasticClient() *elastic.Client {
	return elasticClient
}
