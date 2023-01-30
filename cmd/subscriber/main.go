package main

import (
	"L0task"
	"L0task/pkg/handler"
	"L0task/pkg/repository"
	"L0task/pkg/repository/postgres"
	"L0task/pkg/service"
	"context"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/nats-io/stan.go"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	logrus.Print("Server initialization starting...")

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

	logrus.Print("Initializing services...")
	repos := repository.NewRepository(db)
	services := service.NewService(repos)

	logrus.Print("Restoring data from database to cache...")
	if dbToCacheErr := services.SetOrdersFromDbToCache(); dbToCacheErr != nil {
		logrus.Fatalf("Failed to restore data from database to cache: %s", dbErr.Error())
	}

	logrus.Print("Initializing nats-streaming server...")
	natsStreaming := handler.NewNats(services, validator.New())

	stanConn, stanConnErr := natsStreaming.Connect(
		viper.GetString("nats.cluster_id"),
		viper.GetString("nats.client_subscriber"),
		viper.GetString("nats.url_sub"))
	if stanConnErr != nil {
		return
	}
	defer func(sc stan.Conn) {
		if scErr := sc.Close(); scErr != nil {
			logrus.Errorf("Failed to close subscriber connection to nats streaming server: %s", scErr.Error())
		}
	}(stanConn)

	logrus.Print("Subscribing to nats subject \"orders\"...")
	var waitGroup sync.WaitGroup
	waitGroup.Add(1)
	go func() {
		if subErr := natsStreaming.Subscribe(&waitGroup, stanConn, viper.GetString("nats.subject")); subErr != nil {
			return
		}
	}()
	logrus.Print("Subscription to nats-streaming subject succeed")

	handlers := handler.NewHandler(services)

	logrus.Print("Starting server...")
	srv := new(L0task.Server)
	if err := srv.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil {
		logrus.Fatalf("Failed to run http server: %s", err.Error())
	}

	logrus.Print("Server started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	logrus.Print("Server shutting down...")

	if srvCloseErr := srv.Shutdown(context.Background()); srvCloseErr != nil {
		logrus.Errorf("Failed to shutdown server: %s", srvCloseErr.Error())
	}

	if dbCloseErr := db.Close(); dbCloseErr != nil {
		logrus.Errorf("Failed to close database connection: %s", dbCloseErr.Error())
	}

	if scErr := stanConn.Close(); scErr != nil {
		logrus.Errorf("Failed to close nats-streaming connection: %s", scErr.Error())
	}
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")

	return viper.ReadInConfig()
}
