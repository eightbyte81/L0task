package main

import (
	"L0task"
	"L0task/pkg/handler"
	"L0task/pkg/repository"
	"L0task/pkg/repository/postgres"
	"L0task/pkg/service"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
)

func main() {
	logrus.Print("L0task app starting...")

	logrus.Print("Initializing configs...")
	if configErr := initConfig(); configErr != nil {
		logrus.Fatalf("Error initializing configs: %s", configErr.Error())
	}

	logrus.Print("Loading env variables...")
	if envErr := godotenv.Load(); envErr != nil {
		logrus.Fatalf("Error loading env variables: %s", envErr.Error())
	}

	logrus.Print("Initializing database...")
	db, dbErr := postgres.NewPostgresDB(postgres.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.ssl_mode"),
		Password: os.Getenv("DB_PASSWORD"),
	})
	if dbErr != nil {
		logrus.Fatalf("Failed to initialize db: %s", dbErr.Error())
	}

	logrus.Print("Starting server...")
	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	logrus.Print("Restoring data from database to cache...")
	if dbToCacheErr := services.SetOrdersFromDbToCache(); dbToCacheErr != nil {
		logrus.Fatalf("Failed to restore data from database to cache: %s", dbErr.Error())
	}

	srv := new(L0task.Server)
	if err := srv.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil {
		logrus.Fatalf("Error occurred while running http server: %s", err.Error())
	}

	logrus.Print("Server started")
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")

	return viper.ReadInConfig()
}
