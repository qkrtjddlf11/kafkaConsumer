package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/qkrtjddlf11/kafkaConsumer/common"

	influxdb2 "github.com/influxdata/influxdb-client-go"
	"github.com/influxdata/influxdb-client-go/api"
	"github.com/segmentio/kafka-go"
)

type telegraf struct {
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

func writeInfluxPoint(w api.WriteAPIBlocking, host, hostname_ip, svr_id, vrc, level, alertName, message, value string) {
	p := influxdb2.NewPoint("alertServer",
		map[string]string{
			"host":        host,
			"hostname_ip": hostname_ip,
			"svr_id":      svr_id,
			"vrc":         vrc,
			"level":       level,
			"alertName":   alertName,
		},
		map[string]interface{}{
			"message": message,
			"value":   value,
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
			log.Println(err)
			break
		}

		log.Printf("Message at Topic -> %v, Partition -> %v, Offset -> %v, Key -> %s, Value -> %s\n", message.Topic, message.Partition, message.Offset, string(message.Key), string(message.Value))

		telegraf := telegraf{}
		err = json.Unmarshal([]uint8(string(message.Value)), &telegraf)
		if err != nil {
			log.Fatal(err)
		}

		switch telegraf.Name {
		case "mem":
			telegrafMemory := common.TelegrafMemory{}
			err := json.Unmarshal([]uint8(string(message.Value)), &telegrafMemory)
			if err != nil {
				log.Fatal(err)
			}

			if telegrafMemory.Fields.UsedPercent > 80.0 {
				value := fmt.Sprintf("%1f", telegrafMemory.Fields.UsedPercent)
				writeInfluxPoint(writeAPI, telegrafMemory.Tags.Host, telegrafMemory.Tags.HostnameIP, telegrafMemory.Tags.SvrID, telegrafMemory.Tags.Vrc, "mem-used-percent", "", "<TEST>", value)
			}
		}

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
