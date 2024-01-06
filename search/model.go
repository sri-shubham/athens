package search

import (
	"context"
	"fmt"

	elastic "github.com/olivere/elastic/v7"
	"go.uber.org/zap"
)

// ProjectDocument represents the Elasticsearch document structure for projects
type ProjectDocument struct {
	ProjectID   int       `json:"project_id"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Description string    `json:"description"`
	User        User      `json:"users"`
	Hashtags    []Hashtag `json:"hashtags"`
}

type User struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type Hashtag struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

const indexName = "projects"

func createIndex(ctx context.Context, client *elastic.Client) error {
	// Create the index with a custom analyzer for full-text search
	createIndex, err := client.CreateIndex(indexName).
		Body(`{
			"settings": {
				"analysis": {
					"analyzer": {
						"custom_analyzer": {
							"tokenizer": "standard",
							"filter": ["lowercase", "asciifolding"]
						}
					}
				}
			},
			"mappings": {
				"properties": {
					"project_id": {
						"type": "integer"
					},
					"name": {
						"type": "text",
						"analyzer": "custom_analyzer"
					},
					"slug": {
						"type": "text",
						"analyzer": "custom_analyzer"
					},
					"description": {
						"type": "text",
						"analyzer": "custom_analyzer"
					},
					"users": {
						"type": "nested",
						"properties": {
							"id": {
								"type": "integer"
							},
							"name": {
								"type": "text",
								"analyzer": "custom_analyzer"
							},
							"created_at": {
								"type": "date"
							}
						}
					},
					"hashtags": {
						"type": "nested",
						"properties": {
							"id": {
								"type": "integer"
							},
							"name": {
								"type": "text",
								"analyzer": "custom_analyzer"
							},
							"created_at": {
								"type": "date"
							}
						}
					}
				}
			}
		}`).
		Do(ctx)
	if err != nil {
		return err
	}

	// Check if the index was created successfully
	if !createIndex.Acknowledged {
		return fmt.Errorf("index creation not acknowledged")
	}

	return nil
}

func InitIndex(ctx context.Context, client *elastic.Client) error {
	// Check if the index already exists
	exists, err := client.IndexExists(indexName).Do(ctx)
	if err != nil {
		zap.L().Error("Elastic search index check failed", zap.Error(err))
		return err
	}

	if !exists {
		err = createIndex(ctx, client)
		if err != nil {
			zap.L().Error("Elastic search index create failed", zap.Error(err))
			return err
		}
	}
	return nil
}

// SyncData : Add documents to index, optionally can clear index and add documents
func SyncData(ctx context.Context, client *elastic.Client, projects []ProjectDocument, clearData bool) error {
	// Check if the index already exists
	exists, err := client.IndexExists(indexName).Do(ctx)
	if err != nil {
		zap.L().Error("Elastic search index check failed", zap.Error(err))
		return err
	}

	if !exists {
		err = createIndex(ctx, client)
		if err != nil {
			zap.L().Error("Elastic search index create failed", zap.Error(err))
			return err
		}
	} else if clearData {
		// Delete all documents in the index
		deleteByQueryService := client.DeleteByQuery(indexName).Query(elastic.NewMatchAllQuery())
		_, err = deleteByQueryService.Do(ctx)
		if err != nil {
			return err
		}
	}

	// Bulk indexing projects
	bulkRequest := client.Bulk()
	for _, p := range projects {
		indexReq := elastic.NewBulkIndexRequest().Index(indexName).Type("_doc").Doc(p)
		bulkRequest.Add(indexReq)
	}

	// Execute the bulk request
	_, err = bulkRequest.Do(ctx)
	if err != nil {
		return err
	}

	return nil
}
