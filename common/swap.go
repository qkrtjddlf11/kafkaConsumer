package common

import "fmt"

type TelegrafSwap struct {
	Fields struct {
		Free        int     `json:"free"`
		Total       int     `json:"total"`
		Used        int     `json:"used"`
		UsedPercent float64 `json:"used_percent"`
	} `json:"fields"`
	Name string `json:"name"`
	Tags struct {
		Host       string `json:"host"`
		HostnameIP string `json:"hostname_ip"`
		Svctype    string `json:"svctype"`
		SvrID      string `json:"svr_id"`
		Vrc        string `json:"vrc"`
	} `json:"tags"`
	Timestamp int `json:"timestamp"`
}

func CheckTelegrafSwapUsedPercent(telegrafSwap TelegrafSwap, warning, critical int) (string, string, string) {
	var value string
	var level string
	var measurementMessage string

	value = fmt.Sprintf("%.1f", telegrafSwap.Fields.UsedPercent)

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

	return level, value, measurementMessage
}
