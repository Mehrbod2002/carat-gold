package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"sync"
	"time"
)

type DataMeta struct {
	Time      string  `json:"time" bson:"time"`
	Symbol    string  `json:"symbol" bson:"symbol"`
	Ask       float64 `json:"ask" bson:"ask"`
	Bid       float64 `json:"bid" bson:"bid"`
	High      float64 `json:"high" bson:"high"`
	Low       float64 `json:"low" bson:"low"`
	Open      float64 `json:"open" bson:"open"`
	Close     float64 `json:"close" bson:"close"`
	Type      string  `json:"type" bson:"type"`
	Timeframe string  `json:"timeframe" bson:"timeframe"`
}

// var lastBarsDict = make(map[string]DataMeta)
// var lastBarsDictLock sync.Mutex

func handleMetaTrader(completeJSON []byte, dataChannel chan<- DataMeta) {
	var dataMeta DataMeta

	err := json.Unmarshal(completeJSON, &dataMeta)
	if err != nil {
		return
	}
	dataMeta.Time = fmt.Sprintf("%d", time.Now().UTC().Unix())
	dataChannel <- dataMeta

}

func startServerMetaTrader(errors chan<- error, wg *sync.WaitGroup, dataChannel chan<- DataMeta, stop chan struct{}) {
	defer wg.Done()

	for {
		select {
		case <-stop:
			return
		default:
		}

		listener, err := net.Listen("tcp", serverAddress)
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

func handleClientMetatrader(conn net.Conn, dataChannel chan<- DataMeta) {
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
