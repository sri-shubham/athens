package main

import (
	"log"

	"github.com/sri-shubham/athens/models"
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
	_ = users
	userProjects := models.NewPgUserProjectHelper(util.GetDb())
	_ = userProjects
	projects := models.NewPgUserProjectHelper(util.GetDb())
	_ = projects
	projectHashtags := models.NewPgProjectHashtagHelper(util.GetDb())
	_ = projectHashtags
	hashtags := models.NewPgHashtagHelper(util.GetDb())
	_ = hashtags

}
