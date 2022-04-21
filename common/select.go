package common

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go"
	"github.com/influxdata/influxdb-client-go/api"
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
		if err := json.Unmarshal([]uint8(string(message.Value)), &telegrafMemory); err != nil {
			log.Fatal(err)
		}

		value = fmt.Sprintf("%1f", telegrafMemory.Fields.UsedPercent)

		if telegrafMemory.Fields.UsedPercent > 80.0 {
			measurementMessage = CreateMessage("mem", CRITICAL, value, "")
		} else if telegrafMemory.Fields.UsedPercent > 70.0 {
			measurementMessage = CreateMessage("mem", WARNING, value, "")
		} else {
			measurementMessage = CreateMessage("mem", OK, value, "")
		}
		writeInfluxPoint(writeAPI, telegrafMemory.Tags.Host, telegrafMemory.Tags.HostnameIP, telegrafMemory.Tags.SvrID, telegrafMemory.Tags.Vrc, CRITICAL, "mem-used-percent", measurementMessage, value)

	case "cpu":
		telegrafCpu := TelegrafCPU{}
		if err := json.Unmarshal([]uint8(string(message.Value)), &telegrafCpu); err != nil {
			log.Fatal(err)
		}

		usedPercent := 100.0 - telegrafCpu.Fields.UsageIdle
		value = fmt.Sprintf("%1f", usedPercent)

		if usedPercent > 50.0 {
			measurementMessage = CreateMessage("cpu", CRITICAL, value, "")
		} else if usedPercent > 30.0 {
			measurementMessage = CreateMessage("cpu", WARNING, value, "")
		} else {
			measurementMessage = CreateMessage("cpu", OK, value, "")
		}
		writeInfluxPoint(writeAPI, telegrafCpu.Tags.Host, telegrafCpu.Tags.HostnameIP, telegrafCpu.Tags.SvrID, telegrafCpu.Tags.Vrc, CRITICAL, "cpu-used-percent", measurementMessage, value)

	case "disk":
		telegrafDisk := TelegrafDisk{}
		if err := json.Unmarshal([]uint8(string(message.Value)), &telegrafDisk); err != nil {
			log.Fatal(err)
		}

		value = fmt.Sprintf("%1f", telegrafDisk.Fields.UsedPercent)

		if telegrafDisk.Fields.UsedPercent > 90 {
			measurementMessage = CreateMessage("disk", CRITICAL, value, telegrafDisk.Tags.Path)
		} else if telegrafDisk.Fields.UsedPercent > 85 {
			measurementMessage = CreateMessage("disk", WARNING, value, telegrafDisk.Tags.Path)
		} else {
			measurementMessage = CreateMessage("disk", OK, value, telegrafDisk.Tags.Path)
		}
		writeInfluxPoint(writeAPI, telegrafDisk.Tags.Host, telegrafDisk.Tags.HostnameIP, telegrafDisk.Tags.SvrID, telegrafDisk.Tags.Vrc, CRITICAL, "disk-"+telegrafDisk.Tags.Path+"-used-percent", measurementMessage, value)
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
