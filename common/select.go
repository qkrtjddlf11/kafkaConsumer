package common

import (
	"context"
	"encoding/json"
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
	//file, _ := ioutil.ReadFile("/Users/parksungil/go/src/github.com/qkrtjddlf11/kafkaConsumer/common.json")
	file, _ := ioutil.ReadFile("/tmp/common.json")
	telegrafCommon := []TelegrafCommon{}
	_ = json.Unmarshal([]byte(file), &telegrafCommon)

	return telegrafCommon
}

func SelectNameOfTelegraf(message kafka.Message, typeOf string, writeAPI api.WriteAPIBlocking) {
	var critical int
	var warning int

	telegrafCommon := getTelegrafCommonJson()
	for _, v := range telegrafCommon {
		strtmp := strings.Split(v.MonType, "-")
		if strtmp[0] == typeOf {
			critical = v.Critical
			warning = v.Warning
			//log.Println(warning, critical, typeOf)
			break
		}
	}

	switch typeOf {
	case "mem":
		telegrafMemory := TelegrafMemory{}
		if err := json.Unmarshal([]uint8(string(message.Value)), &telegrafMemory); err != nil {
			log.Fatal(err)
		}

		level, value, measurementMessage := CheckTelegrafMemoryUsedPercent(telegrafMemory, warning, critical)
		writeInfluxPoint(writeAPI, telegrafMemory.Tags.Host, telegrafMemory.Tags.HostnameIP, telegrafMemory.Tags.SvrID, telegrafMemory.Tags.Vrc, level, "mem-used-percent", measurementMessage, value)

	case "cpu":
		telegrafCpu := TelegrafCPU{}
		if err := json.Unmarshal([]uint8(string(message.Value)), &telegrafCpu); err != nil {
			log.Fatal(err)
		}

		level, value, measurementMessage := CheckTelegrafCPUUsedPercent(telegrafCpu, warning, critical)
		writeInfluxPoint(writeAPI, telegrafCpu.Tags.Host, telegrafCpu.Tags.HostnameIP, telegrafCpu.Tags.SvrID, telegrafCpu.Tags.Vrc, level, "cpu-used-percent", measurementMessage, value)

	case "disk":
		telegrafDisk := TelegrafDisk{}
		if err := json.Unmarshal([]uint8(string(message.Value)), &telegrafDisk); err != nil {
			log.Fatal(err)
		}

		level, value, measurementMessage := CheckTelegrafDiskUsedPercent(telegrafDisk, warning, critical)
		writeInfluxPoint(writeAPI, telegrafDisk.Tags.Host, telegrafDisk.Tags.HostnameIP, telegrafDisk.Tags.SvrID, telegrafDisk.Tags.Vrc, level, "disk-used-percent-"+telegrafDisk.Tags.Path, measurementMessage, value)

	case "swap":
		telegrafSwap := TelegrafSwap{}
		if err := json.Unmarshal([]uint8(string(message.Value)), &telegrafSwap); err != nil {
			log.Fatal("Error -> ", err)
		}

		if telegrafSwap.Fields.Total != 0 {
			level, value, measurementMessage := CheckTelegrafSwapUsedPercent(telegrafSwap, warning, critical)
			writeInfluxPoint(writeAPI, telegrafSwap.Tags.Host, telegrafSwap.Tags.HostnameIP, telegrafSwap.Tags.SvrID, telegrafSwap.Tags.Vrc, level, "swap-used-percent", measurementMessage, value)
		}

	//{"fields":{"load1":0,"load15":0.05,"load5":0.01,"n_cpus":4,"n_users":2},"name":"system","tags":{"host":"KAFKA-VM01","hostname_ip":"KAFKA-VM01_172.30.1.210","svctype":"Control","svr_id":"kafka-vm01","vrc":"Control"},"timestamp":1652878860}
	case "system":
		telegrafLoad5 := Load5{}
		if err := json.Unmarshal([]uint8(string(message.Value)), &telegrafLoad5); err != nil {
			log.Fatal("Error -> ", err)
		}

		if telegrafLoad5.Fields.NCpus != 0 {
			level, value, measurementMessage := CheckTelegrafLoad5Percent(telegrafLoad5, warning, critical)
			writeInfluxPoint(writeAPI, telegrafLoad5.Tags.Host, telegrafLoad5.Tags.HostnameIP, telegrafLoad5.Tags.SvrID, telegrafLoad5.Tags.Vrc, level, "load5-used-percent", measurementMessage, value)
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
