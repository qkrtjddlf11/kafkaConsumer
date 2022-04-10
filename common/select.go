package common

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/segmentio/kafka-go"
)

const (
	CRITICAL = "CRITICAL"
	WARNING  = "WARNING"
	OK       = "OK"
)

func SelectNameOfTelegraf(message kafka.Message, typeOf string, writeAPI api.WriteAPIBlocking) {
	var value string
	var measurementMessage string

	switch typeOf {
	case "mem":
		telegrafMemory := TelegrafMemory{}
		err := json.Unmarshal([]uint8(string(message.Value)), &telegrafMemory)
		if err != nil {
			log.Fatal(err)
		}

		value = fmt.Sprintf("%1f", telegrafMemory.Fields.UsedPercent)

		if telegrafMemory.Fields.UsedPercent > 80.0 {
			measurementMessage = CreateMessage("mem", CRITICAL, value)
		} else if telegrafMemory.Fields.UsedPercent > 70.0 {
			measurementMessage = CreateMessage("mem", WARNING, value)
		} else {
			measurementMessage = CreateMessage("mem", OK, value)
		}

		writeInfluxPoint(writeAPI, telegrafMemory.Tags.Host, telegrafMemory.Tags.HostnameIP, telegrafMemory.Tags.SvrID, telegrafMemory.Tags.Vrc, CRITICAL, "mem-used-percent", measurementMessage, value)
	}
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
