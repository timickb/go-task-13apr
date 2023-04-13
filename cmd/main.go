package main

import (
	"github.com/sirupsen/logrus"
	"os"
	"refactoring/internal/config"
	"refactoring/internal/server"
	"refactoring/internal/service"
	"refactoring/internal/store/filestore"
	"strconv"
)

func main() {
	logger := logrus.New()
	cfg := config.NewDefault()
	fillConfigFromEnv(cfg)

	store, err := filestore.New(cfg.StoreFile)
	if err != nil {
		logger.Fatal(err)
	}

	user := service.NewUser(store, logger)
	srv := server.New(cfg, logger, user)

	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}

func fillConfigFromEnv(cfg *config.Config) {
	if os.Getenv("APP_PORT") != "" {
		cfg.AppPort, _ = strconv.Atoi(os.Getenv("APP_PORT"))
	}
	if os.Getenv("STORE_FILE") != "" {
		cfg.StoreFile = os.Getenv("STORE_FILE")
	}
}
