package tests

import (
	"testing"
	"watchtopus/agent/collectors"
	"watchtopus/orm"
)

func TestCollectorCpu(t *testing.T) {
	ch := make(chan []orm.MetricFloat)
	go collectors.CollectCpu(ch)
	metricsCpu := <-ch

	// Assert user cpu metric
	metric := findMetricByKey(metricsCpu, "cpu.user.cpu0", t)
	if metric.Category != "cpu" {
		t.Errorf("Bad category for cpu.user.cpu0")
	}
	if metric.SubCategory != "user" {
		t.Errorf("Bad subcategory for cpu.user.cpu0")
	}
	if metric.Component != "cpu0" {
		t.Errorf("Bad component for cpu.user.cpu0")
	}
	if metric.Val < 0 || metric.Val > 100 {
		t.Errorf("Bad value for cpu.user.cpu0")
	}

	// Assert user cpu metric
	metric = findMetricByKey(metricsCpu, "cpu.system.cpu0", t)
	if metric.Category != "cpu" {
		t.Errorf("Bad category for cpu.system.cpu0")
	}
	if metric.SubCategory != "system" {
		t.Errorf("Bad subcategory for cpu.system.cpu0")
	}
	if metric.Component != "cpu0" {
		t.Errorf("Bad component for cpu.system.cpu0")
	}
	if metric.Val < 0 || metric.Val > 100 {
		t.Errorf("Bad value for cpu.system.cpu0")
	}
}

func TestCollectorMem(t *testing.T) {
	ch := make(chan []orm.MetricFloat)
	go collectors.CollectMem(ch)
	metricsMem := <-ch

	// Assert user cpu metric
	metricFree := findMetricByKey(metricsMem, "mem.free", t)
	metricUsed := findMetricByKey(metricsMem, "mem.used", t)
	metricTotal := findMetricByKey(metricsMem, "mem.total", t)

	// Assert free mem metric fields
	if metricFree.Category != "mem" {
		t.Errorf("Bad category for mem.free")
	}
	if metricFree.SubCategory != "free" {
		t.Errorf("Bad subcategory for mem.free")
	}
	if metricFree.Component != "" {
		t.Errorf("Bad component for mem.free")
	}
	if metricFree.Val < 0 || metricFree.Val > (100*1024*1024*1024) {
		t.Errorf("Bad value for mem.free")
	}

	// Assert free + used = total
	if metricFree.Val+metricUsed.Val != metricTotal.Val {
		t.Errorf("Bad value for mem.total")
	}
}

func findMetricByKey(metrics []orm.MetricFloat, key string, t *testing.T) orm.MetricFloat {
	for _, metric := range metrics {
		if metric.Key == key {
			return metric
		}
	}

	t.Errorf("Couldn't find key %s", key)

	return orm.MetricFloat{}
}
