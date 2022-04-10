package common

type TelegrafCPU struct {
	Fields struct {
		UsageGuest     int     `json:"usage_guest"`
		UsageGuestNice int     `json:"usage_guest_nice"`
		UsageIdle      float64 `json:"usage_idle"`
		UsageIowait    float64 `json:"usage_iowait"`
		UsageIrq       int     `json:"usage_irq"`
		UsageNice      int     `json:"usage_nice"`
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
