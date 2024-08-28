package main

import (
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/gorilla/websocket"
)

var (
	serverAddress  = ":5741"
	wssClients     = make(map[*websocket.Conn]struct{})
	wssClientsLock sync.Mutex
	upgrader       = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	lastData     DataMeta
	lastDataLock sync.Mutex
)

func main() {
	log.SetOutput(os.Stdout)
	errors := make(chan error)
	dataChannel := make(chan DataMeta, 100)
	stop := make(chan struct{})
	var wg sync.WaitGroup

	wg.Add(2)
	go startServerWSS(errors, &wg, 5050, dataChannel)
	go startServerMetaTrader(errors, &wg, dataChannel, stop)

	wg.Wait()
	close(errors)
	close(dataChannel)
}
