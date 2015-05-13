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
		"Mallocs":      stats.Mallocs,
		"Frees":        stats.Mallocs,
		"Inbound":      uint64(messages.Get("inbound")),
		"Outbound":     uint64(messages.Get("outbound")),
	}
}
