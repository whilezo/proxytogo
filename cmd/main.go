// Package main is main package of the proxy server application
package main

import (
	"log"
	"proxy/proxy"
	"sync"
)

func main() {
	var wg sync.WaitGroup

	config, err := proxy.LoadConfig("etc/config.yaml")
	if err != nil {
		log.Fatal(err)
	}

	for i := range config.Listeners {
		wg.Add(1)
		go proxy.StartProxy(&config.Listeners[i], config.Debug, &wg)
	}
	wg.Wait()
}
