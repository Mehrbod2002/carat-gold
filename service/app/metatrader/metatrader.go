package metatrader

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"sync"
	"time"
)

var SharedConnection net.Conn

func startServerMetaTrader(
	errors chan<- error,
	wg *sync.WaitGroup,
	dataChannel chan<- interface{},
	stop chan struct{},
	adminChannel chan interface{},
	sharedReader chan map[string]interface{}) {
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
				go handleClientMetatrader(conn, dataChannel, adminChannel, sharedReader)
			}
		}()
	}
}

func handleClientMetatrader(conn net.Conn,
	dataChannel chan<- interface{},
	adminChannel chan interface{},
	sharedReader chan map[string]interface{}) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	buffer := make([]byte, 0, 1024)

	for {
		select {
		case <-adminChannel:
		default:
		}

		data, err := reader.ReadBytes('}')
		if err != nil {
			return
		}

		buffer = append(buffer[:0], data...)

		for {
			start := bytes.Index(buffer, []byte{'{'})
			end := bytes.Index(buffer, []byte{'}'})
			if start != -1 && end != -1 && start < end {
				completeJSON := buffer[start : end+1]
				go handleMetaTrader(completeJSON, dataChannel, adminChannel, sharedReader)

				buffer = buffer[end+1:]
			} else {
				break
			}
		}
	}
}

func handleMetaTrader(completeJSON []byte, dataChannel chan<- interface{}, adminChannel chan<- interface{}, sharedReader chan map[string]interface{}) {
	var dataMeta map[string]interface{}

	err := json.Unmarshal(completeJSON, &dataMeta)
	if err != nil {
		return
	}

	dataMeta["time"] = fmt.Sprintf("%d", time.Now().UTC().Unix())
	if _, admin := dataMeta["status"]; admin {
		sharedReader <- dataMeta
	} else {
		dataChannel <- dataMeta
	}
}
