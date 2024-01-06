package api

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sri-shubham/athens/models"
)

type routeConf struct {
	prefix string
}

func SetupRoutes(
	users models.Users,
	projects models.Projects,
	projectHashtags models.ProjectHashtags,
	hashtags models.Hashtags,
	userProjects models.UserProjects,
	search Searcher,
) *mux.Router {
	router := mux.NewRouter()
	routesMap := map[GenericCrud]routeConf{
		NewUsersCrud(users):                     {prefix: "/users"},
		NewProjectsCrud(projects):               {prefix: "/projects"},
		NewUserProjectsCrud(userProjects):       {prefix: "/user/projects"},
		NewHashtagsCrud(hashtags):               {prefix: "/hashtags"},
		NewProjectHashtagsCrud(projectHashtags): {prefix: "/project/hashtags"},
	}

	for handlers, route := range routesMap {
		router.HandleFunc(route.prefix, handlers.Create).Methods(http.MethodPost)
		router.HandleFunc(route.prefix+"/{id}", handlers.Read).Methods(http.MethodGet)
		router.HandleFunc(route.prefix, handlers.Update).Methods(http.MethodPut)
		router.HandleFunc(route.prefix+"/{id}", handlers.Delete).Methods(http.MethodDelete)
	}

	router.HandleFunc("/search", search.Fuzzy).Methods(http.MethodGet)
	router.HandleFunc("/search/user", search.ByUser).Methods(http.MethodGet)
	router.HandleFunc("/search/hashtag", search.ByHashTag).Methods(http.MethodGet)

	return router
}
