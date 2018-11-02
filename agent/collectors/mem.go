package collectors

import (
	linuxproc "github.com/c9s/goprocinfo/linux"
	"github.com/spf13/viper"
	"watchtopus/orm"
)

func CollectMem(ch chan []orm.MetricFloat) {
	metrics := make([]orm.MetricFloat, 0)

	//var stats runtime.MemStats
	//runtime.ReadMemStats(&stats)
	stat1, err := linuxproc.ReadMemInfo("/proc/meminfo")
	if err != nil {
		logger.Fatal("Stat read failed")

		ch <- metrics

		return
	}

	metrics = append(metrics, orm.MetricFloat{
		HostId:      viper.GetString("hostId"),
		Key:         "mem.total",
		Val:         float64(stat1.MemTotal),
		Category:    "mem",
		SubCategory: "total",
		Component:   ""})

	metrics = append(metrics, orm.MetricFloat{
		HostId:      viper.GetString("hostId"),
		Key:         "mem.free",
		Val:         float64(stat1.MemFree),
		Category:    "mem",
		SubCategory: "free",
		Component:   ""})

	// Calculation of used memory:
	// https://stackoverflow.com/questions/41224738/how-to-calculate-system-memory-usage-from-proc-meminfo-like-htop
	metrics = append(metrics, orm.MetricFloat{
		HostId:      viper.GetString("hostId"),
		Key:         "mem.used",
		Val:         float64(stat1.MemTotal - stat1.MemFree),
		Category:    "mem",
		SubCategory: "used",
		Component:   ""})

	ch <- metrics
}
