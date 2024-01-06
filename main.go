package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/sri-shubham/athens/api"
	"github.com/sri-shubham/athens/models"
	"github.com/sri-shubham/athens/search"
	"github.com/sri-shubham/athens/util"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
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

	// Init Queues
	userUpdates := util.NewQueue()
	projectUpdates := util.NewQueue()
	userProjectUpdates := util.NewQueue()
	hashtagUpdates := util.NewQueue()
	projectHashtagUpdates := util.NewQueue()

	// Init DB Models
	users := models.NewPgUserHelper(util.GetDb(), userUpdates)
	userProjects := models.NewPgUserProjectHelper(util.GetDb(), userProjectUpdates)
	projects := models.NewPgProjectHelper(util.GetDb(), projectUpdates)
	projectHashtags := models.NewPgProjectHashtagHelper(util.GetDb(), projectHashtagUpdates)
	hashtags := models.NewPgHashtagHelper(util.GetDb(), hashtagUpdates)

	// Init search indexes
	searchModel := search.NewProjectSearcher(util.GetElasticClient())
	searchService := api.NewSearchService(searchModel)

	// Init Search Index
	syncHelper := search.NewSyncHelper(util.GetElasticClient(),
		users,
		projects,
		projectHashtags,
		hashtags,
		userProjects,
		userUpdates,
		projectUpdates,
		projectHashtagUpdates,
		hashtagUpdates,
		userProjectUpdates)

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

	wg, wgCtx := errgroup.WithContext(context.Background())

	// Create a new HTTP server
	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	wg.Go(func() error {
		return server.ListenAndServe()
	})

	wg.Go(func() error {
		<-wgCtx.Done()

		// Create a context with a timeout for server shutdown
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Attempt to gracefully shutdown the server
		if err := server.Shutdown(shutdownCtx); err != nil {
			zap.L().Error("Failed to shutdown server", zap.Error(err))
			panic(err)
		}

		return nil
	})

	wg.Go(func() error {
		return syncHelper.BackgroundSync(wgCtx)
	})

	log.Fatal(wg.Wait())
}
