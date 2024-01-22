package metatrader

import (
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

func startServerWSS(errors chan<- error,
	wg *sync.WaitGroup,
	dataChannel <-chan interface{},
	adminChannel chan interface{}) {
	defer wg.Done()

	http.HandleFunc("/feed", func(w http.ResponseWriter, r *http.Request) {
		handleWebSocket(w, r, dataChannel, adminChannel)
	})

	server := &http.Server{Addr: fmt.Sprintf(":%s", os.Getenv("FEED_REALTIME"))}
	err := server.ListenAndServe()
	if err != nil {
		errors <- fmt.Errorf("WebSocket server error: %v", err)
	}
}

func handleWebSocket(w http.ResponseWriter, r *http.Request, dataChannel <-chan interface{}, adminChannel chan interface{}) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	// isAdmin := utils.ValidateAdmin(r.URL.Query().Get("Authorization"))

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

	go func() {
		for {
			var receivedData interface{}
			err := conn.ReadJSON(&receivedData)
			if err != nil {
				return
			}

			// if isAdmin {
			// 	adminChannel <- receivedData
			// }
		}
	}()

	<-closeSignal
}
