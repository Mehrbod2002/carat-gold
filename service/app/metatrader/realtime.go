package metatrader

import (
	"fmt"
	"net/http"
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

func HandleWebSocket(w http.ResponseWriter, r *http.Request, dataChannel <-chan interface{}) {
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
				fmt.Println(data)
				err := conn.WriteJSON(data)
				if err != nil {
					return
				}
			case <-closeSignal:
				return
			}
		}
	}()

	go func() {
		for {
			var receivedData interface{}
			err := conn.ReadJSON(&receivedData)
			if err != nil {
				return
			}
		}
	}()

	<-closeSignal
}
