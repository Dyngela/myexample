package main

import (
	"memcache"
)

func main() {
	config := memcache.NewConfig()
	go memcache.MonitorMemory(config.MemoryThreshold)
	store := memcache.NewStore(config.DefaultExpiration)
	store.Set("key", "value")
	server := memcache.NewServer()
	server.Start(config)

}
