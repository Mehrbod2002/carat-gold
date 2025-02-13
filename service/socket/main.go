package main

import (
	"log"
	"os"
	"sync"
)

var (
	serverAddress = ":5742"
)

func main() {
	log.SetOutput(os.Stdout)
	errors := make(chan error)
	dataChannel := make(chan DataMeta)
	stop := make(chan struct{})
	var wg sync.WaitGroup

	wg.Add(2)
	go startServerWSS(errors, &wg, 5051, dataChannel)
	go startServerMetaTrader(errors, &wg, dataChannel, stop)

	wg.Wait()
	close(errors)
	close(dataChannel)
}
