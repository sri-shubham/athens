package util

import (
	"github.com/olivere/elastic"
	"github.com/pkg/errors"
)

var elasticClient *elastic.Client

func ConnectToElastic(config *Config) (err error) {
	elasticClient, err = elastic.NewClient(elastic.SetURL(config.ElasticSearch.URL))
	if err != nil {
		return errors.Wrap(err, "Error connecting to Elasticsearch")
	}

	return nil
}
