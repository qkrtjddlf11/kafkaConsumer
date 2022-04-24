package common

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go"
	"github.com/influxdata/influxdb-client-go/api"
	"github.com/segmentio/kafka-go"
)

type TelegrafCommon struct {
	SeqID      int    `json:"seq_id"`
	MonType    string `json:"mon_type"`
	Warning    int    `json:"warning"`
	Critical   int    `json:"critical"`
	Individual int    `json:"individual"`
	CreateDate string `json:"create_date"`
	UpdateDate string `json:"update_date"`
}

const (
	CRITICAL = "CRITICAL"
	WARNING  = "WARNING"
	OK       = "OK"
)

func getTelegrafCommonJson() []TelegrafCommon {
	file, _ := ioutil.ReadFile("/tmp/common.json")
	telegrafCommon := []TelegrafCommon{}
	_ = json.Unmarshal([]byte(file), &telegrafCommon)

	return telegrafCommon
}

func SelectNameOfTelegraf(message kafka.Message, typeOf string, writeAPI api.WriteAPIBlocking) {
	var value string
	var level string
	var critical int
	var warning int
	var measurementMessage string

	telegrafCommon := getTelegrafCommonJson()
	for _, v := range telegrafCommon {
		strtmp := strings.Split(v.MonType, "-")
		if strtmp[0] == typeOf {
			critical = v.Critical
			warning = v.Warning
			break
		}
	}

	switch typeOf {
	case "mem":
		telegrafMemory := TelegrafMemory{}
		if err := json.Unmarshal([]uint8(string(message.Value)), &telegrafMemory); err != nil {
			log.Fatal(err)
		}

		value = fmt.Sprintf("%.1f", telegrafMemory.Fields.UsedPercent)

		if telegrafMemory.Fields.UsedPercent > float64(critical) {
			level = CRITICAL
			measurementMessage = CreateMessage("mem", level, value, "")
		} else if telegrafMemory.Fields.UsedPercent > float64(warning) {
			level = WARNING
			measurementMessage = CreateMessage("mem", level, value, "")
		} else {
			level = OK
			measurementMessage = CreateMessage("mem", level, value, "")
		}
		writeInfluxPoint(writeAPI, telegrafMemory.Tags.Host, telegrafMemory.Tags.HostnameIP, telegrafMemory.Tags.SvrID, telegrafMemory.Tags.Vrc, level, "mem-used-percent", measurementMessage, value)

	case "cpu":
		telegrafCpu := TelegrafCPU{}
		if err := json.Unmarshal([]uint8(string(message.Value)), &telegrafCpu); err != nil {
			log.Fatal(err)
		}

		usedPercent := 100.0 - telegrafCpu.Fields.UsageIdle
		value = fmt.Sprintf("%.1f", usedPercent)

		if usedPercent > float64(critical) {
			level = CRITICAL
			measurementMessage = CreateMessage("cpu", level, value, "")
		} else if usedPercent > float64(warning) {
			level = WARNING
			measurementMessage = CreateMessage("cpu", level, value, "")
		} else {
			level = OK
			measurementMessage = CreateMessage("cpu", level, value, "")
		}
		writeInfluxPoint(writeAPI, telegrafCpu.Tags.Host, telegrafCpu.Tags.HostnameIP, telegrafCpu.Tags.SvrID, telegrafCpu.Tags.Vrc, level, "cpu-used-percent", measurementMessage, value)

	case "disk":
		telegrafDisk := TelegrafDisk{}
		if err := json.Unmarshal([]uint8(string(message.Value)), &telegrafDisk); err != nil {
			log.Fatal(err)
		}

		value = fmt.Sprintf("%.1f", telegrafDisk.Fields.UsedPercent)

		if telegrafDisk.Fields.UsedPercent > float64(critical) {
			level = CRITICAL
			measurementMessage = CreateMessage("disk", level, value, telegrafDisk.Tags.Path)
		} else if telegrafDisk.Fields.UsedPercent > float64(warning) {
			level = WARNING
			measurementMessage = CreateMessage("disk", level, value, telegrafDisk.Tags.Path)
		} else {
			level = OK
			measurementMessage = CreateMessage("disk", level, value, telegrafDisk.Tags.Path)
		}
		writeInfluxPoint(writeAPI, telegrafDisk.Tags.Host, telegrafDisk.Tags.HostnameIP, telegrafDisk.Tags.SvrID, telegrafDisk.Tags.Vrc, level, "disk-used-percent-"+telegrafDisk.Tags.Path, measurementMessage, value)

	case "swap":
		telegrafSwap := TelegrafSwap{}
		if err := json.Unmarshal([]uint8(string(message.Value)), &telegrafSwap); err != nil {
			log.Fatal("Error -> ", err)
		}

		if telegrafSwap.Fields.Total != 0 {
			value := fmt.Sprintf("%1f", telegrafSwap.Fields.UsedPercent)

			if telegrafSwap.Fields.UsedPercent > float64(critical) {
				level = CRITICAL
				measurementMessage = CreateMessage("swap", level, value, "")
			} else if telegrafSwap.Fields.UsedPercent > float64(warning) {
				level = WARNING
				measurementMessage = CreateMessage("swap", level, value, "")
			} else {
				level = OK
				measurementMessage = CreateMessage("swap", level, value, "")
			}
			writeInfluxPoint(writeAPI, telegrafSwap.Tags.Host, telegrafSwap.Tags.HostnameIP, telegrafSwap.Tags.SvrID, telegrafSwap.Tags.Vrc, level, "swap-used-percent", measurementMessage, value)
		}
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
