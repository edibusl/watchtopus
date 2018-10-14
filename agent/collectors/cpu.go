package collectors

import (
	"fmt"
	linuxproc "github.com/c9s/goprocinfo/linux"
	"github.com/op/go-logging"
	"github.com/spf13/viper"
	"math"
	"time"
	"watchtopus/orm"
)

var logger = logging.MustGetLogger("watchtopus")

func CollectCpu(ch chan []orm.MetricFloat) {
	// Read first sample
	stat1, err := linuxproc.ReadStat("/proc/stat")
	if err != nil {
		logger.Fatal("Stat read failed")
	}

	// Wait
	time.Sleep(1 * time.Second)

	// Read second sample
	stat2, err := linuxproc.ReadStat("/proc/stat")
	if err != nil {
		logger.Fatal("Stat read failed")
	}

	metrics := make([]orm.MetricFloat, 0)

	// Go through all CPU cores
	for i := 0; i < len(stat1.CPUStats); i++ {
		// Calc the diffs between 2 samples
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
		calcPercent := func(val, total uint64) float64 { return math.Trunc(float64(val)/float64(total)*100.0*100.0) / 100.0 }

		metrics = append(metrics, orm.MetricFloat{
			HostId:      viper.GetString("hostId"),
			Key:         fmt.Sprintf("cpu.user.%s", diff.Id),
			Val:         calcPercent(diff.User, total),
			Category:    "cpu",
			SubCategory: "user",
			Component:   diff.Id})

		metrics = append(metrics, orm.MetricFloat{
			HostId:      viper.GetString("hostId"),
			Key:         fmt.Sprintf("cpu.system.%s", diff.Id),
			Val:         calcPercent(diff.System, total),
			Category:    "cpu",
			SubCategory: "system",
			Component:   diff.Id})
	}

	ch <- metrics
}
