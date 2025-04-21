package main

import (
	"context"
	"go.uber.org/zap"
	"lang/internal/handler"
	"lang/internal/repository"
	"lang/internal/server"
	"lang/internal/service"
	loggers2 "lang/logger"
	"lang/pkg/database"
	"lang/pkg/utils"
	"log"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/lib/pq"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// TODO rating
// TODO quiz submit
// TODO
func main() {

	logrus.SetFormatter(new(logrus.JSONFormatter))

	initLogger, err := loggers2.InitLogger()
	if err != nil {
		log.Fatal(err)
		return
	}
	log.Println("Init logger")
	defer func(logger *zap.Logger) {
		err = logger.Sync()
		if err != nil {
			log.Println(err)
		}
	}(initLogger)

	if err := utils.InitConfig(); err != nil {
		logrus.Fatalf("error initializing config: %s", err.Error())
	}

	log.Println("Init configs")
	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("Error loading .env file: %s", err.Error())
	}

	db, err := database.InitPostgres(database.ConfigsDb{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		DbName:   viper.GetString("db.dbname"),
		Username: viper.GetString("db.username"),
		Password: os.Getenv("DB_PASSWORD"),
		SLLmode:  viper.GetString("db.sslmode"),
	})

	if err != nil {
		logrus.Fatalf("error to initialize db: %s", err.Error())
	}

	log.Println("Connection with DB")

	repo := repository.NewRepository(db)
	services := service.NewService(repo)
	handlers := handler.NewHandler(services, initLogger)

	log.Println("Connection dependencies")

	srv := new(server.Server)
	go func() {
		if err := srv.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil {
			logrus.Fatal("Error occurred while running http Server:", err.Error())
		}
	}()

	log.Println("TodoApp started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	if err := srv.Shutdown(context.Background()); err != nil {
		logrus.Fatal("Error occurred while shutting down http server:", err.Error())
	}

}
