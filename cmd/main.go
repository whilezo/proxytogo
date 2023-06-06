// Package main is main package of the proxy server application
package main

import (
	"fmt"
	"proxy/proxy"
	"sync"
)

func main() {
	var wg sync.WaitGroup

	config, err := proxy.LoadConfig("etc/config.yaml")
	if err != nil {
		panic(err)
	}

	fmt.Println(config.Listeners)
	for i := range config.Listeners {
		wg.Add(1)
		go proxy.StartProxy(&config.Listeners[i], config.Debug)
	}
	wg.Wait()
}
