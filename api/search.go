package api

import (
	"net/http"

	"github.com/sri-shubham/athens/search"
)

type Search struct {
	search search.ProjectSearcher
}

type Searcher interface {
	ByUser(w http.ResponseWriter, r *http.Request)
	ByHashTag(w http.ResponseWriter, r *http.Request)
	Fuzzy(w http.ResponseWriter, r *http.Request)
}

func NewSearchService(search search.ProjectSearcher) Searcher {
	return &Search{
		search: search,
	}
}

// ByHashTag implements Searcher.
func (s *Search) ByHashTag(w http.ResponseWriter, r *http.Request) {
	searchTerm := r.URL.Query().Get("search")
	if searchTerm == "" {
		http.Error(w, "searchTerm is required", http.StatusBadRequest)
		return
	}

	projectDocuments, err := s.search.ByHashTag(r.Context(), searchTerm)
	if err != nil {
		http.Error(w, "failed to search projects "+err.Error(), http.StatusInternalServerError)
		return
	}

	if projectDocuments == nil {
		projectDocuments = []*search.ProjectDocument{}
	}

	sendJSONResponse(w, projectDocuments, http.StatusOK)
}

// ByUser implements Searcher.
func (s *Search) ByUser(w http.ResponseWriter, r *http.Request) {
	searchTerm := r.URL.Query().Get("search")
	if searchTerm == "" {
		http.Error(w, "searchTerm is required", http.StatusBadRequest)
		return
	}

	projectDocuments, err := s.search.ByUser(r.Context(), searchTerm)
	if err != nil {
		http.Error(w, "failed to search projects", http.StatusInternalServerError)
		return
	}

	if projectDocuments == nil {
		projectDocuments = []*search.ProjectDocument{}
	}

	sendJSONResponse(w, projectDocuments, http.StatusOK)
}

// Fuzzy implements Searcher.
func (s *Search) Fuzzy(w http.ResponseWriter, r *http.Request) {
	searchTerm := r.URL.Query().Get("search")
	if searchTerm == "" {
		http.Error(w, "searchTerm is required", http.StatusBadRequest)
		return
	}

	if len(searchTerm) < 5 {
		http.Error(w, "searchTerm too small", http.StatusBadRequest)
		return
	}

	projectDocuments, err := s.search.Fuzzy(r.Context(), searchTerm)
	if err != nil {
		http.Error(w, "failed to search projects"+err.Error(), http.StatusInternalServerError)
		return
	}

	if projectDocuments == nil {
		projectDocuments = []*search.ProjectDocument{}
	}

	sendJSONResponse(w, projectDocuments, http.StatusOK)
}
