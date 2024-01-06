package main

import (
	"context"
	"log"
	"net/http"

	"github.com/sri-shubham/athens/api"
	"github.com/sri-shubham/athens/models"
	"github.com/sri-shubham/athens/search"
	"github.com/sri-shubham/athens/util"
	"go.uber.org/zap"
)

func main() {
	// Set up pre requisites
	logger, err := zap.NewProduction()
	if err != nil {
		log.Panic("failed to init zap logger")
	}

	zap.ReplaceGlobals(logger)

	conf, err := util.ParseNewConfig("config/config.yaml")
	if err != nil {
		zap.L().Panic("Failed to parse config", zap.Error(err))
	}

	err = util.ConnectToPostgres(conf)
	if err != nil {
		zap.L().Panic("Failed to connect to postgres", zap.Error(err))
	}

	err = util.InitPostgresDB()
	if err != nil {
		zap.L().Panic("Failed to initialize to postgres", zap.Error(err))
	}

	err = util.ConnectToElastic(conf)
	if err != nil {
		zap.L().Panic("Failed to connect to elastic search", zap.Error(err))
	}

	// Init DB Models
	users := models.NewPgUserHelper(util.GetDb())
	userProjects := models.NewPgUserProjectHelper(util.GetDb())
	projects := models.NewPgProjectHelper(util.GetDb())
	projectHashtags := models.NewPgProjectHashtagHelper(util.GetDb())
	hashtags := models.NewPgHashtagHelper(util.GetDb())

	// Init search indexes
	searchModel := search.NewProjectSearcher(util.GetElasticClient())
	searchService := api.NewSearchService(searchModel)

	// Init Search Index
	syncHelper := search.NewSyncHelper(util.GetElasticClient(),
		users,
		projects,
		projectHashtags,
		hashtags,
		userProjects)

	err = syncHelper.SyncAll(context.Background())
	if err != nil {
		panic(err)
	}

	// Setup Routes
	router := api.SetupRoutes(users,
		projects,
		projectHashtags,
		hashtags,
		userProjects,
		searchService)

	log.Fatal(http.ListenAndServe(":8080", router))
}
