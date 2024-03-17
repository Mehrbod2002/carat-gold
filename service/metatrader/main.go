package main

import (
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var (
	serverAddress = ":5741"
)

func main() {
	errors := make(chan error)
	var wg sync.WaitGroup

	dataChannel := make(chan DataMeta, 200)

	stop := make(chan struct{})
	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
		<-sig
		close(stop)
	}()
	wg.Add(3)
	go startServerWSS(errors, &wg, 5050, dataChannel)
	go startServerMetaTrader(errors, &wg, dataChannel, stop)

	<-stop
	close(errors)
	wg.Wait()
}
