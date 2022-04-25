package common

import "fmt"

type TelegrafCPU struct {
	Fields struct {
		UsageGuest     int     `json:"usage_guest"`
		UsageGuestNice int     `json:"usage_guest_nice"`
		UsageIdle      float64 `json:"usage_idle"`
		UsageIowait    float64 `json:"usage_iowait"`
		UsageIrq       int     `json:"usage_irq"`
		UsageNice      float64 `json:"usage_nice"`
		UsageSoftirq   float64 `json:"usage_softirq"`
		UsageSteal     int     `json:"usage_steal"`
		UsageSystem    float64 `json:"usage_system"`
		UsageUser      float64 `json:"usage_user"`
	} `json:"fields"`
	Name string `json:"name"`
	Tags struct {
		CPU        string `json:"cpu"`
		Host       string `json:"host"`
		HostnameIP string `json:"hostname_ip"`
		Svctype    string `json:"svctype"`
		SvrID      string `json:"svr_id"`
		Vrc        string `json:"vrc"`
	} `json:"tags"`
	Timestamp int `json:"timestamp"`
}

func CheckTelegrafCPUUsedPercent(telegrafCpu TelegrafCPU, warning, critical int) (string, string, string) {
	var level string
	var measurementMessage string

	usedPercent := 100.0 - telegrafCpu.Fields.UsageIdle
	value := fmt.Sprintf("%.1f", usedPercent)

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

	return level, value, measurementMessage
}
