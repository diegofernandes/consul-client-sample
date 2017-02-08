package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"time"

	consulapi "github.com/hashicorp/consul/api"
)

func main() {
	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	config := consulapi.DefaultConfig()
	config.Address = "localhost:8500"
	consul, err := consulapi.NewClient(config)
	if err != nil {
		log.Print(err)
	}

	kv := consul.KV()

	go func() {
		var index uint64 = 0
		for {
			index = watcher(kv, index)
		}
	}()

	<-sigs
}

func watcher(kv *consulapi.KV, index uint64) uint64 {
	keys, meta, err := kv.List("test", &consulapi.QueryOptions{
		WaitTime:  10 * time.Second,
		WaitIndex: index,
	})

	log.Printf("watcher...")

	if err != nil {
		log.Fatal(err)
		return 0
	}

	for _, key := range keys {
		if key.ModifyIndex > index {
			log.Printf("Key: %s Value: %s", key.Key, key.Value)
		}
	}

	return meta.LastIndex

}
