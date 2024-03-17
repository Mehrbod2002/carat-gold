package main

import (
	"log"
	"sync"
)

var (
	serverAddress = ":5741"
)

func main() {
	errors := make(chan error)
	dataChannel := make(chan DataMeta)
	stop := make(chan struct{})
	var wg sync.WaitGroup

	wg.Add(2)
	go startServerWSS(errors, &wg, 5050, dataChannel)
	go startServerMetaTrader(errors, &wg, dataChannel, stop)

	go func() {
		wg.Wait()
		close(errors)
		close(dataChannel)
	}()

	for err := range errors {
		log.Println("Error:", err)
	}
}
