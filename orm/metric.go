package orm

import "time"

type MetricFloat struct {
	Key         string    `json:"key"`
	Val         float64   `json:"val"`
	Category    string    `json:"category"`    //cpu, mem, network
	SubCategory string    `json:"subcategory"` //user, system, idle
	Component   string    `json:"component"`   //Optional - cpu0, cpu1, cpu2, etc..
	Timestamp   time.Time `json:"timestamp"`
}
