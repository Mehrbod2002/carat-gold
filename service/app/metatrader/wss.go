package metatrader

import (
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func InitiateMetatrader(dataChannel chan interface{}) {
	errors := make(chan error)
	var wg sync.WaitGroup

	stop := make(chan struct{})
	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
		<-sig
		close(stop)
	}()

	wg.Add(1)
	go startServerMetaTrader(errors, &wg, dataChannel, stop)

	<-stop
	close(errors)
	wg.Wait()
}
