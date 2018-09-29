package collectors

import (
	linuxproc "github.com/c9s/goprocinfo/linux"
	"log"
	"watchtopus/orm"
)

func CollectMem() (metrics []orm.MetricFloat) {
	metrics = make([]orm.MetricFloat, 0)

	//var stats runtime.MemStats
	//runtime.ReadMemStats(&stats)
	stat1, err := linuxproc.ReadMemInfo("/proc/meminfo")
	if err != nil {
		log.Fatal("stat read fail")
	}

	metrics = append(metrics, orm.MetricFloat{
		Key:         "mem.total",
		Val:         float64(stat1.MemTotal),
		Category:    "mem",
		SubCategory: "total",
		Component:   ""})

	metrics = append(metrics, orm.MetricFloat{
		Key:         "mem.free",
		Val:         float64(stat1.MemFree),
		Category:    "mem",
		SubCategory: "free",
		Component:   ""})

	// Calculation of used memory:
	// https://stackoverflow.com/questions/41224738/how-to-calculate-system-memory-usage-from-proc-meminfo-like-htop
	metrics = append(metrics, orm.MetricFloat{
		Key:         "mem.used",
		Val:         float64(stat1.MemTotal - stat1.MemFree),
		Category:    "mem",
		SubCategory: "used",
		Component:   ""})

	return
}
