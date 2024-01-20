package metatrader

import (
	"bufio"
	"bytes"
	"carat-gold/models"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"sync"
	"time"
)

func startServerMetaTrader(errors chan<- error, wg *sync.WaitGroup, dataChannel chan<- models.DataMeta) {
	defer wg.Done()

	for {
		listener, err := net.Listen("tcp", os.Getenv("MQ5_PORT"))
		if err != nil {
			errors <- err
			continue
		}

		go func() {
			defer listener.Close()
			for {
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

func handleClientMetatrader(conn net.Conn, dataChannel chan<- models.DataMeta) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	buffer := make([]byte, 0, 1024)
	for {
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
				go handleMetaTrader(completeJSON, dataChannel)

				buffer = buffer[end+1:]
			} else {
				break
			}
		}
	}
}

func handleMetaTrader(completeJSON []byte, dataChannel chan<- models.DataMeta) {
	var dataMeta models.DataMeta

	err := json.Unmarshal(completeJSON, &dataMeta)
	if err != nil {
		return
	}

	dataMeta.Time = fmt.Sprintf("%d", time.Now().UTC().Unix())
	dataChannel <- dataMeta
}
