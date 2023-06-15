package data

type Stat struct {
	NCpu      int       `json:"nCpu"`
	CpuStats  []CpuStat `json:"cpuStats"`
	MemUsed   uint64    `json:"mem_used"`
	MemBuff   uint64    `json:"mem_buff"`
	MemShared uint64    `json:"mem_shared"`
	MemFree   uint64    `json:"mem_free"`
	NetRx     uint64    `json:"net_rx"`
	NetTx     uint64    `json:"net_tx"`
}

type CpuStat struct {
	Idle  uint64 `json:"idle"`
	Total uint64 `json:"total"`
}
