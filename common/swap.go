package common

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
