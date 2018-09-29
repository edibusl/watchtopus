package collectors

import (
	"fmt"
	linuxproc "github.com/c9s/goprocinfo/linux"
	"log"
	"math"
	"time"
)


//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// Option 2 - Parse by myself:
//idle0, total0 := getCPUSample()
//time.Sleep(3 * time.Second)
//idle1, total1 := getCPUSample()
//idleTicks := float64(idle1 - idle0)
//totalTicks := float64(total1 - total0)
//cpuUsage := 100 * (totalTicks - idleTicks) / totalTicks
//fmt.Printf("CPU usage is %f%% [busy: %f, total: %f]\n", cpuUsage, totalTicks-idleTicks, totalTicks)

// Option 3 - Use mpstat -P ALL 1 1 (it makes the calculations by itslef)
// Stackoverflow: https://stackoverflow.com/questions/11356330/getting-cpu-usage-with-golang
//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~


type MetricFloat struct {
	Key string
	Val float64
	Category string //cpu, memory, network
	SubCategory string //user, system, idle
	Component string //Optional - cpu0, cpu1, cpu2, etc..
}

func CollectCpu() (metrics []MetricFloat) {
	stat1, err := linuxproc.ReadStat("/proc/stat")
	if err != nil {
		log.Fatal("stat read fail")
	}

	time.Sleep(1 * time.Second)

	stat2, err := linuxproc.ReadStat("/proc/stat")
	if err != nil {
		log.Fatal("stat read fail")
	}


	metrics = make([]MetricFloat, 0)

	for i := 0; i < len(stat1.CPUStats); i++ {
		diff := linuxproc.CPUStat{
			stat1.CPUStats[i].Id,
			stat2.CPUStats[i].User - stat1.CPUStats[i].User,
			stat2.CPUStats[i].Nice - stat1.CPUStats[i].Nice,
			stat2.CPUStats[i].System - stat1.CPUStats[i].System,
			stat2.CPUStats[i].Idle - stat1.CPUStats[i].Idle,
			stat2.CPUStats[i].IOWait - stat1.CPUStats[i].IOWait,
			stat2.CPUStats[i].IRQ - stat1.CPUStats[i].IRQ,
			stat2.CPUStats[i].SoftIRQ - stat1.CPUStats[i].SoftIRQ,
			stat2.CPUStats[i].Steal - stat1.CPUStats[i].Steal,
			stat2.CPUStats[i].Guest - stat1.CPUStats[i].Guest,
			stat2.CPUStats[i].GuestNice - stat1.CPUStats[i].GuestNice}

		total := diff.User + diff.Nice + diff.System + diff.Idle + diff.IOWait + diff.IRQ + diff.SoftIRQ + diff.Steal + diff.Guest + diff.GuestNice


		// Lambda func to calc percentage out of total and round up to 2 decimals places after dot
		calcPercent := func(val, total uint64) float64 {return math.Trunc(float64(val) / float64(total) * 100.0 * 100.0) / 100.0}


		metrics = append(metrics, MetricFloat{
			Key: fmt.Sprintf("cpu.user.%s", diff.Id),
			Val: calcPercent(diff.User, total),
			Category: "cpu",
			SubCategory: "user",
			Component: diff.Id})

		metrics = append(metrics, MetricFloat{
			Key: fmt.Sprintf("cpu.system.%s", diff.Id),
			Val: calcPercent(diff.System, total),
			Category: "cpu",
			SubCategory: "system",
			Component: diff.Id})
	}

	return
}