package util

import (
	"github.com/c9s/goprocinfo/linux"
	"linux-web-top/data"
)

func PopulateNCpu(stat *data.Stat) error {
	statInfo, err := linux.ReadStat("/proc/stat")
	if err != nil {
		return err
	}
	stat.NCpu = len(statInfo.CPUStats)
	stat.CpuStats = make([]data.CpuStat, stat.NCpu)
	return nil
}

func PopulateMemInfo(stat *data.Stat) error {
	meminfo, err := linux.ReadMemInfo("/proc/meminfo")
	if err != nil {
		return err
	}
	stat.MemFree = meminfo.MemFree
	stat.MemBuff = meminfo.Buffers + meminfo.Cached
	stat.MemShared = meminfo.Shmem
	stat.MemUsed = meminfo.MemTotal - (stat.MemFree + stat.MemBuff + stat.MemShared)
	return nil
}

func PopulateCpuInfo(stat *data.Stat) error {
	statInfo, err := linux.ReadStat("/proc/stat")
	if err != nil {
		return err
	}
	for i, cpuStat := range statInfo.CPUStats {
		stat.CpuStats[i].Idle = cpuStat.Idle + cpuStat.IOWait
		stat.CpuStats[i].Total = cpuStat.User +
			cpuStat.Nice +
			cpuStat.System +
			stat.CpuStats[i].Idle +
			cpuStat.IRQ +
			cpuStat.SoftIRQ +
			cpuStat.Steal +
			cpuStat.Guest +
			cpuStat.GuestNice
	}
	return nil
}

func PopulateNetworkBandwith(stat *data.Stat) error {
	netInfos, err := linux.ReadNetworkStat("/proc/net/dev")
	if err != nil {
		return err
	}
	var sumRx uint64
	var sumTx uint64
	for _, info := range netInfos {
		sumRx += info.RxBytes
		sumTx += info.TxBytes
	}
	stat.NetRx = sumRx / 1024
	stat.NetTx = sumTx / 1024
	return nil
}
