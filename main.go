package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go"
	"github.com/influxdata/influxdb-client-go/api"
	"github.com/segmentio/kafka-go"
)

type telegrafJson struct {
	Fields struct {
	} `json:"fields"`
	Name string `json:"name"`
	Tags struct {
	} `json:"tags"`
	Timestamp int `json:"timestamp"`
}

// Kafka options
var (
	topic = flag.String(
		"topic",
		"",
		"Kafka Topic\nUsage : -topic=telegraf")

	partition = flag.Int(
		"partition",
		-1,
		"Topic Partition\nUsage : -partition=3")

	kafkaBrokers = flag.String(
		"brokers",
		"",
		"Kafka Broker Servers\nUsage : -brokers=172.30.1.210:9092")

	influx = flag.String(
		"influx",
		"",
		"InfluxDB Server\nUsage : -influx=172.30.1.220:8086")
)

// initializeKafka initialize Kafka configuration
func initializeKafka() *kafka.Reader {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{*kafkaBrokers},
		GroupID: "telegraf",
		Topic:   *topic,
		//Partition: 1,
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

func initializeInfluxDB() api.WriteAPIBlocking {
	client := influxdb2.NewClient(*influx, "")
	writeAPI := client.WriteAPIBlocking("", "telegraf")

	return writeAPI
}

func WriteInfluxPoint(w api.WriteAPIBlocking) {
	p := influxdb2.NewPoint("alertServer",
		map[string]string{
			"host":        "t.Tags.Host",
			"hostname_ip": "t.Tags.Hostname_IP",
			"svr_id":      "t.Tags.Svr_Id",
			"vrc":         "t.Tags.Vrc",
			"level":       "level",
			"alertName":   "alertName",
		},
		map[string]interface{}{
			"message": "message",
			"value":   "value",
		},
		time.Now())

	if err := w.WritePoint(context.Background(), p); err != nil {
		log.Fatal(err)
	}
}

func main() {
	flag.Parse()

	if *topic == "" {
		printUsageAndErrorAndExit("-topic is required")
	}
	if *partition == -1 {
		printUsageAndErrorAndExit("-partition is required")
	}
	if *kafkaBrokers == "" {
		printUsageAndErrorAndExit("-brokers is required")
	}
	if *influx == "" {
		printUsageAndErrorAndExit("-influx is required")
	}

	r := initializeKafka()
	writeAPI := initializeInfluxDB()
	ctx := context.Background()
	for {
		message, err := r.FetchMessage(ctx)
		if err != nil {
			fmt.Println(err)
			break
		}

		fmt.Printf("Message at Topic -> %v, Partition -> %v, Offset -> %v, Key -> %s, Value -> %s\n", message.Topic, message.Partition, message.Offset, string(message.Key), string(message.Value))
		WriteInfluxPoint(writeAPI)
		if err := r.CommitMessages(ctx, message); err != nil {
			log.Fatal("Failed to commit messages :", err)
		}
	}

	if err := r.Close(); err != nil {
		log.Fatal("Failed to close reader :", err)
	}
}

func printUsageAndErrorAndExit(message string) {
	fmt.Fprintln(os.Stderr, "Error :", message)
	fmt.Println("Available command line options :")
	flag.PrintDefaults()
	os.Exit(1)
}
