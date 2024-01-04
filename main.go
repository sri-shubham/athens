package main

import (
	"context"
	"fmt"
	"log"

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

	// Init Search Index
	syncHelper := search.NewSyncHelper(util.GetElasticClient(),
		users,
		projects,
		projectHashtags,
		hashtags,
		userProjects)

	user, err := users.Create(&models.User{
		Name: "Shubham",
	})
	if err != nil {
		panic(err)
	}

	ht, err := hashtags.Create(&models.Hashtag{
		Name: "Go",
	})
	if err != nil {
		panic(err)
	}

	project, err := projects.Create(&models.Project{
		Name:        "Project 1",
		Slug:        "T12345",
		Description: "Test description",
	})
	if err != nil {
		panic(err)
	}

	projectHashtag, err := projectHashtags.Create(&models.ProjectHashtag{
		HashtagID: ht,
		ProjectID: project,
	})
	if err != nil {
		panic(err)
	}

	userProject, err := userProjects.Create(&models.UserProject{
		ProjectID: project,
		UserID:    user,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(user, project, ht, projectHashtag, userProject)

	err = syncHelper.SyncAll(context.Background())
	if err != nil {
		panic(err)
	}
}
