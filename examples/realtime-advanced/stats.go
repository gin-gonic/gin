package main

import (
	"runtime"
	"time"
)

func Stats() map[string]uint64 {
	var stats runtime.MemStats
	runtime.ReadMemStats(&stats)

	return map[string]uint64{
		"timestamp":    uint64(time.Now().Unix()),
		"HeapInuse":    stats.HeapInuse,
		"StackInuse":   stats.StackInuse,
		"NuGoroutines": uint64(runtime.NumGoroutine()),
		//"Latency":      latency,
		"Mallocs": stats.Mallocs,
		"Frees":   stats.Mallocs,
		// "HeapIdle":     stats.HeapIdle,
		// "HeapInuse":    stats.HeapInuse,
		// "HeapReleased": stats.HeapReleased,
		// "HeapObjects":  stats.HeapObjects,
	}
}
