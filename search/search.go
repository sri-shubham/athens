package search

import (
	"context"
	"encoding/json"

	elastic "github.com/olivere/elastic/v7"
)

type ProjectSearcher interface {
	ByUser(context.Context, string) ([]*ProjectDocument, error)
	ByHashTag(context.Context, string) ([]*ProjectDocument, error)
	Fuzzy(context.Context, string) ([]*ProjectDocument, error)
}

type ProjectSearchHelper struct {
	elastic *elastic.Client
}

func NewProjectSearcher(elastic *elastic.Client) ProjectSearcher {
	return &ProjectSearchHelper{
		elastic: elastic,
	}
}

// ByHashTag implements ProjectSearcher.
func (h *ProjectSearchHelper) ByHashTag(ctx context.Context, searchTerm string) ([]*ProjectDocument, error) {
	query := elastic.NewNestedQuery("hashtags",
		elastic.NewMatchQuery("hashtags.name", searchTerm),
	)

	searchResult, err := h.elastic.Search().
		Index(indexName).
		Query(query).
		Do(context.Background())
	if err != nil {
		return nil, err
	}

	return getProjectDocuments(searchResult)
}

// ByUser implements ProjectSearcher.
func (h *ProjectSearchHelper) ByUser(ctx context.Context, searchTerm string) ([]*ProjectDocument, error) {
	query := elastic.NewNestedQuery("users",
		elastic.NewMatchQuery("users.name", searchTerm),
	)

	searchResult, err := h.elastic.Search().
		Index(indexName).
		Query(query).
		Do(context.Background())
	if err != nil {
		return nil, err
	}

	return getProjectDocuments(searchResult)
}

// Fuzzy implements ProjectSearcher.
func (h *ProjectSearchHelper) Fuzzy(ctx context.Context, searchTerm string) ([]*ProjectDocument, error) {
	query := elastic.NewMultiMatchQuery(searchTerm, "name^2", "slug^2", "description^1").Fuzziness("AUTO")

	searchResult, err := h.elastic.Search().
		Index(indexName).
		Query(query).
		Do(context.Background())
	if err != nil {
		return nil, err
	}

	return getProjectDocuments(searchResult)
}

func getProjectDocuments(searchResult *elastic.SearchResult) ([]*ProjectDocument, error) {
	var projects []*ProjectDocument
	for _, hit := range searchResult.Hits.Hits {
		project := &ProjectDocument{}

		err := json.Unmarshal(hit.Source, project)
		if err != nil {
			return nil, err
		}

		projects = append(projects, project)
	}

	return projects, nil
}
