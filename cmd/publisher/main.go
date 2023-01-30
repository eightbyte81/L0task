package main

import (
	"github.com/nats-io/stan.go"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"io"
	"os"
)

func main() {
	logrus.Print("[Publisher] Server initialization starting...")
	logrus.Print("[Publisher] Initializing configs...")
	if configErr := initConfig(); configErr != nil {
		logrus.Fatalf("[Publisher] Error initializing configs: %s", configErr.Error())
	}

	logrus.Print("[Publisher] Initializing nats-streaming server...")
	stanConn, stanConnErr := stan.Connect(
		viper.GetString("nats.cluster_id"),
		viper.GetString("nats.client_producer"),
		stan.NatsURL(viper.GetString("nats.url_pub")))
	if stanConnErr != nil {
		logrus.Fatalf("[Publisher] Failed to connect to the nats-streaming server: %s", stanConnErr.Error())
	}
	defer func(sc stan.Conn) {
		if scErr := sc.Close(); scErr != nil {
			logrus.Errorf("Failed to close subscriber connection to nats streaming server: %s", scErr.Error())
		}
	}(stanConn)

	logrus.Print("[Publisher] Parsing JSON model...")
	jsonData, jsonErr := os.Open(viper.GetString("json.static_model_path"))
	if jsonErr != nil {
		logrus.Fatalf("[Publisher] Failed to open JSON file: %s", jsonErr)
	}
	defer func(jd *os.File) {
		if jdErr := jd.Close(); jdErr != nil {
			logrus.Errorf("[Publisher] Failed to close JSON file: %s", jdErr.Error())
		}
	}(jsonData)

	byteData, readErr := io.ReadAll(jsonData)
	if readErr != nil {
		logrus.Fatalf("[Publisher] Failed to serialize JSON file: %s", readErr)
	}

	logrus.Print("[Publisher] Sending data to nats-streaming server...")
	if pubErr := stanConn.Publish(viper.GetString("nats.subject"), byteData); pubErr != nil {
		logrus.Fatalf("[Publisher] Failed to publish data to nats-streaming server: %s", pubErr)
	}

	logrus.Print("[Publisher] Publishing static JSON file to nats-streaming server succeed")
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")

	return viper.ReadInConfig()
}
