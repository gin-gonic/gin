package main

import (
	"runtime"
	"sync"
	"time"
)

var mutexStats sync.RWMutex
var savedStats map[string]uint64

func statsWorker() {
	c := time.Tick(1 * time.Second)
	for range c {
		var stats runtime.MemStats
		runtime.ReadMemStats(&stats)

		mutexStats.Lock()
		savedStats = map[string]uint64{
			"timestamp":    uint64(time.Now().Unix()),
			"HeapInuse":    stats.HeapInuse,
			"StackInuse":   stats.StackInuse,
			"NuGoroutines": uint64(runtime.NumGoroutine()),
			"Mallocs":      stats.Mallocs,
			"Frees":        stats.Mallocs,
			"Inbound":      uint64(messages.Get("inbound")),
			"Outbound":     uint64(messages.Get("outbound")),
		}
		mutexStats.Unlock()
	}
}

func Stats() map[string]uint64 {
	mutexStats.RLock()
	defer mutexStats.RUnlock()

	return savedStats
}
