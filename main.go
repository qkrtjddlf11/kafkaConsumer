package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/qkrtjddlf11/kafkaConsumer/common"

	influxdb2 "github.com/influxdata/influxdb-client-go"
	"github.com/influxdata/influxdb-client-go/api"
	"github.com/segmentio/kafka-go"
)

type Telegraf struct {
	Fields struct {
	} `json:"fields"`
	Name string `json:"name"`
	Tags struct {
	} `json:"tags"`
	Timestamp int `json:"timestamp"`
}

type Env struct {
	Topic     string
	Partition string
	Brokers   string
	Influx    string
}

// initializeKafka initialize Kafka configuration
func initializeKafka(env Env) *kafka.Reader {
	// Autthentication For SASL
	/*
		mechanism := plain.Mechanism{
			Username: "admin",
			Password: "admin-secret",
		}

		dailer := &kafka.Dialer{
			Timeout:       time.Second * 5,
			DualStack:     true,
			SASLMechanism: mechanism,
		}
	*/

	parseBrokers := strings.Split(env.Brokers, ",")
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: parseBrokers,
		GroupID: "telegraf",
		Topic:   env.Topic,
		//Partition:         *partition,
		MinBytes:          10e3, // 10KB
		MaxBytes:          10e6, // 10MB
		MaxAttempts:       5,
		MaxWait:           time.Second * 10,
		HeartbeatInterval: time.Second * 3,
		CommitInterval:    time.Second * 1,
	})
	r.SetOffset(-1)

	return r
}

func initializeInfluxDB(env Env) api.WriteAPIBlocking {
	client := influxdb2.NewClient(env.Influx, "")
	writeAPI := client.WriteAPIBlocking("", "telegraf")

	return writeAPI
}

func getEnv() Env {
	if err := godotenv.Load("/usr/service/etc/.env"); err != nil {
		log.Fatal("Error loading .env file.")
	}

	env := Env{
		Topic:     os.Getenv("TOPIC"),
		Partition: os.Getenv("PARTITION"),
		Brokers:   os.Getenv("BROKERS"),
		Influx:    os.Getenv("INFLUX"),
	}

	return env
}

func main() {
	env := getEnv()
	r := initializeKafka(env)
	writeAPI := initializeInfluxDB(env)

	ctx := context.Background()
	for {
		message, err := r.FetchMessage(ctx)
		if err != nil {
			log.Panic(err)
			break
		}

		//log.Printf("Message at Topic -> %v, Partition -> %v, Offset -> %v, Key -> %s, Value -> %s\n", message.Topic, message.Partition, message.Offset, string(message.Key), string(message.Value))

		telegraf := Telegraf{}
		err = json.Unmarshal([]uint8(string(message.Value)), &telegraf)
		if err != nil {
			log.Fatal(err)
		}

		common.SelectNameOfTelegraf(message, telegraf.Name, writeAPI)

		if err := r.CommitMessages(ctx, message); err != nil {
			log.Fatal("Failed to commit messages :", err)
		}
	}

	if err := r.Close(); err != nil {
		log.Fatal("Failed to close reader :", err)
	}
}

func gotEnv() {
	panic("unimplemented")
}

func printUsageAndErrorAndExit(message string) {
	fmt.Fprintln(os.Stderr, "Error :", message)
	fmt.Println("Available command line options :")
	flag.PrintDefaults()
	os.Exit(1)
}
