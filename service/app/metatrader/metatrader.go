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

func startServerMetaTrader(
	errors chan<- error,
	wg *sync.WaitGroup,
	dataChannel chan<- interface{},
	stop chan struct{}) {
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

				go handleClientMetatrader(conn, dataChannel)
			}
		}()
	}
}

func handleClientMetatrader(conn net.Conn,
	dataChannel chan<- interface{}) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	buffer := make([]byte, 0, 1024)
	for {
		data, err := reader.ReadBytes('}')
		if err != nil {
			fmt.Println(865, err)
			return
		}

		buffer = append(buffer[:0], data...)

		for {
			start := bytes.Index(buffer, []byte{'{'})
			end := bytes.Index(buffer, []byte{'}'})
			if start != -1 && end != -1 && start < end {
				completeJSON := buffer[start : end+1]
				var dataMeta map[string]interface{}
				err := json.Unmarshal(completeJSON, &dataMeta)
				if err == nil {
					go handleMetaTrader(dataMeta, dataChannel)
					buffer = buffer[end+1:]
				}
			} else {
				break
			}
		}
	}
}

func handleMetaTrader(dataMeta map[string]interface{}, dataChannel chan<- interface{}) {
	dataMeta["time"] = fmt.Sprintf("%d", time.Now().UTC().Unix())
	dataChannel <- dataMeta
}
