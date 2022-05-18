package common

import "fmt"

type Load5 struct {
	Fields struct {
		Load1  float64 `json:"load1"`
		Load15 float64 `json:"load15"`
		Load5  float64 `json:"load5"`
		NCpus  int     `json:"n_cpus"`
		NUsers int     `json:"n_users"`
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

func CheckTelegrafLoad5Percent(telegrafLoad5 Load5, warning, critical int) (string, string, string) {
	var level string
	var measurementMessage string

	nCpus := telegrafLoad5.Fields.NCpus * 100
	load5 := telegrafLoad5.Fields.Load5 * float64(nCpus)

	value := load5 / float64(nCpus)
	strValue := fmt.Sprintf("%.1f", load5/float64(nCpus))

	if value > float64(critical) {
		level = CRITICAL
		measurementMessage = CreateMessage("load5", level, strValue, "")
	} else if value > float64(warning) {
		level = WARNING
		measurementMessage = CreateMessage("load5", level, strValue, "")
	} else {
		level = OK
		measurementMessage = CreateMessage("load5", level, strValue, "")
	}

	return level, strValue, measurementMessage
}
