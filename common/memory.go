package common

import (
	"fmt"
)

type TelegrafMemory struct {
	Fields struct {
		Active           int     `json:"active"`
		Available        int64   `json:"available"`
		AvailablePercent float64 `json:"available_percent"`
		Buffered         int     `json:"buffered"`
		Cached           int     `json:"cached"`
		CommitLimit      int64   `json:"commit_limit"`
		CommittedAs      int     `json:"committed_as"`
		Dirty            int     `json:"dirty"`
		Free             int     `json:"free"`
		HighFree         int     `json:"high_free"`
		HighTotal        int     `json:"high_total"`
		HugePageSize     int     `json:"huge_page_size"`
		HugePagesFree    int     `json:"huge_pages_free"`
		HugePagesTotal   int     `json:"huge_pages_total"`
		Inactive         int     `json:"inactive"`
		LowFree          int     `json:"low_free"`
		LowTotal         int     `json:"low_total"`
		Mapped           int     `json:"mapped"`
		PageTables       int     `json:"page_tables"`
		Shared           int     `json:"shared"`
		Slab             int     `json:"slab"`
		Sreclaimable     int     `json:"sreclaimable"`
		Sunreclaim       int     `json:"sunreclaim"`
		SwapCached       int     `json:"swap_cached"`
		SwapFree         int     `json:"swap_free"`
		SwapTotal        int     `json:"swap_total"`
		Total            int64   `json:"total"`
		Used             int     `json:"used"`
		UsedPercent      float64 `json:"used_percent"`
		VmallocChunk     int64   `json:"vmalloc_chunk"`
		VmallocTotal     int64   `json:"vmalloc_total"`
		VmallocUsed      int     `json:"vmalloc_used"`
		WriteBack        int     `json:"write_back"`
		WriteBackTmp     int     `json:"write_back_tmp"`
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

func CheckTelegrafMemoryUsedPercent(telegrafMemory TelegrafMemory, warning, critical int) (string, string, string) {
	var level string
	var measurementMessage string

	value := fmt.Sprintf("%.1f", telegrafMemory.Fields.UsedPercent)

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

	return level, value, measurementMessage
}
