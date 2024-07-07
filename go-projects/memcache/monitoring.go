package memcache

import (
	"fmt"
	"runtime"
	"time"
)

func MonitorMemory(memoryThreshold uint64) {
	var memStats runtime.MemStats

	for {
		runtime.ReadMemStats(&memStats)
		if memStats.Alloc > memoryThreshold {
			fmt.Println("Memory usage exceeded the threshold. Cleaning up...")
			//TODO: Implement the cleanup logic
			panic("Memory usage exceeded the threshold")
		}
		time.Sleep(10 * time.Minute)
	}
}
