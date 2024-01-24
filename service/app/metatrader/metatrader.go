package metatrader

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

var SharedConnection net.Conn
var SharedReader map[string]interface{} = make(map[string]interface{})

func startServerMetaTrader(
	errors chan<- error,
	wg *sync.WaitGroup,
	dataChannel chan<- interface{},
	stop chan struct{},
	adminChannel chan interface{}) {
	defer wg.Done()

	for {
		select {
		case <-stop:
			return
		default:
		}

		listener, err := net.Listen("tcp", ":"+os.Getenv("MQ5_PORT"))
		if err != nil {
			errors <- err
			continue
		}

		go func() {
			defer listener.Close()
			for {
				select {
				case <-stop:
					return
				default:
				}

				conn, err := listener.Accept()
				if err != nil {
					errors <- err
					return
				}

				SharedConnection = conn
				go handleClientMetatrader(conn, dataChannel, adminChannel)
			}
		}()
	}
}

func handleClientMetatrader(conn net.Conn,
	dataChannel chan<- interface{},
	adminChannel chan interface{}) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	buffer := make([]byte, 0, 1024)

	for {
		select {
		case <-adminChannel:
		default:
		}

		temp := make([]byte, 1024)
		_, err := reader.Read(temp)
		if err != nil {
			return
		}

		buffer = append(buffer, temp...)
		buffer = []byte(strings.ReplaceAll(string(buffer), "\x00", ""))

		for {
			var dataMeta map[string]interface{}
			err := json.Unmarshal(buffer, &dataMeta)
			if err == nil {
				go handleMetaTrader(dataMeta, dataChannel, adminChannel)
				buffer = buffer[:0]
			} else {
				break
			}
		}
	}
}

func handleMetaTrader(dataMeta map[string]interface{}, dataChannel chan<- interface{}, adminChannel chan<- interface{}) {
	dataMeta["time"] = fmt.Sprintf("%d", time.Now().UTC().Unix())
	id, ok := dataMeta["id"].(string)
	if ok {
		SharedReader[id] = dataMeta
	} else {
		dataChannel <- dataMeta
	}
}
