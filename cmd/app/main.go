package main

import (
	"context"
	"fmt"
	"lang/internal/handler"
	"lang/internal/repository"
	"lang/internal/server"
	"lang/internal/service"
	"lang/pkg/database"
	"os"

	_ "github.com/lib/pq"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {

	fmt.Println(os.Getenv("DB_PASSWORD"))
	logrus.SetFormatter(new(logrus.JSONFormatter))

	if err := initConfig(); err != nil {
		logrus.Fatalf("error initializing config: %s", err.Error())
	}

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

	repo := repository.NewRepository(db)
	service := service.NewService(repo)
	handler := handler.NewHandler(service)

	srv := new(server.Server)
	err = srv.Run(viper.GetString("port"), handler.InitRoute())
	if err != nil {
		logrus.Fatalf("Error occurred while running http server: %s", err.Error())
	}

	if err := srv.Shutdown(context.Background()); err != nil {
		logrus.Printf("error occured on server shutting down: %s", err.Error())
	}
	if err := db.Close(); err != nil {
		logrus.Printf("error occured on db close: %s", err.Error())

	}

}

func initConfig() error {
	viper.AddConfigPath("config")
	viper.SetConfigName("configs")
	return viper.ReadInConfig()
}
