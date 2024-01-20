package metatrader

import (
	"carat-gold/models"
	"fmt"
	"net/http"
	"os"
	"sync"

	"github.com/gorilla/websocket"
)

var wssClients = make(map[*websocket.Conn]struct{})
var wssClientsLock sync.Mutex
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func startServerWSS(errors chan<- error, wg *sync.WaitGroup, dataChannel <-chan models.DataMeta) {
	defer wg.Done()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		handleWebSocket(w, r, dataChannel)
	})

	err := http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("FEED_REALTIME")), nil)
	if err != nil {
		errors <- fmt.Errorf("WebSocket server error: %v", err)
	}
}

func handleWebSocket(w http.ResponseWriter, r *http.Request, dataChannel <-chan models.DataMeta) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	wssClientsLock.Lock()
	wssClients[conn] = struct{}{}
	wssClientsLock.Unlock()

	closeSignal := make(chan struct{})
	defer func() {
		wssClientsLock.Lock()
		delete(wssClients, conn)
		wssClientsLock.Unlock()
		close(closeSignal)
		conn.Close()
	}()

	go func() {
		for {
			select {
			case data := <-dataChannel:
				err := conn.WriteJSON(data)
				if err != nil {
					return
				}
			case <-closeSignal:
				return
			}
		}
	}()

	<-closeSignal
}
