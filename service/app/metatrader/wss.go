package metatrader

import (
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func InitiateMetatrader() {
	errors := make(chan error)
	var wg sync.WaitGroup

	dataChannel := make(chan interface{}, 200)
	adminChannel := make(chan interface{}, 200)
	stop := make(chan struct{})
	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
		<-sig
		close(stop)
	}()

	wg.Add(2)
	go startServerWSS(errors, &wg, dataChannel, adminChannel)
	go startServerMetaTrader(errors, &wg, dataChannel, adminChannel)

	<-stop
	close(errors)
	wg.Wait()
}
