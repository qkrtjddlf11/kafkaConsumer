package common

type TelegrafDisk struct {
	Fields struct {
		Free        int64   `json:"free"`
		InodesFree  int     `json:"inodes_free"`
		InodesTotal int     `json:"inodes_total"`
		InodesUsed  int     `json:"inodes_used"`
		Total       int64   `json:"total"`
		Used        int64   `json:"used"`
		UsedPercent float64 `json:"used_percent"`
	} `json:"fields"`
	Name string `json:"name"`
	Tags struct {
		Device     string `json:"device"`
		Fstype     string `json:"fstype"`
		Host       string `json:"host"`
		HostnameIP string `json:"hostname_ip"`
		Mode       string `json:"mode"`
		Path       string `json:"path"`
		Svctype    string `json:"svctype"`
		SvrID      string `json:"svr_id"`
		Vrc        string `json:"vrc"`
	} `json:"tags"`
	Timestamp int `json:"timestamp"`
}
