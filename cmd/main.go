// Package main is main package of the proxy server application
package main

import (
	"log"
	"proxy/proxy"
	"sync"

	"github.com/sirupsen/logrus"
)

func main() {
	var wg sync.WaitGroup

	config, err := proxy.LoadConfig("etc/config.yaml")
	if err != nil {
		log.Fatal(err)
	}

	if config.Debug {
		logrus.SetLevel(logrus.DebugLevel)
	}

	for i := range config.Listeners {
		wg.Add(1)
		go proxy.StartProxy(&config.Listeners[i], &wg)
	}
	wg.Wait()
}
