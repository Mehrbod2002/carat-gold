package main

import (
	"encoding/json"
	"fmt"
	"net"
	"sync"
	"time"
)

type DataMeta struct {
	Time       string  `json:"time" bson:"time"`
	Symbol     string  `json:"symbol" bson:"symbol"`
	Ask        float64 `json:"ask" bson:"ask"`
	Bid        float64 `json:"bid" bson:"bid"`
	High       float64 `json:"high" bson:"high"`
	Low        float64 `json:"low" bson:"low"`
	Open       float64 `json:"open" bson:"open"`
	Close      float64 `json:"close" bson:"close"`
	Type       string  `json:"type" bson:"type"`
	Timeframe  string  `json:"timeframe" bson:"timeframe"`
	ProfitDay  float64 `json:"profit_day" bson:"profit_day"`
	Profithour float64 `json:"profit_hour" bson:"profit_hour"`
	ProfitWeek float64 `json:"profit_week" bson:"profit_week"`
}

func startServerMetaTrader(errors chan<- error, wg *sync.WaitGroup, dataChannel chan<- DataMeta, stop chan struct{}) {
	defer wg.Done()

	listener, err := net.Listen("tcp", serverAddress)
	if err != nil {
		errors <- err
		return
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			errors <- err
			continue
		}

		go handleConnection(conn, dataChannel, errors, stop)
	}
}

func handleConnection(conn net.Conn, dataChannel chan<- DataMeta, errors chan<- error, stop chan struct{}) {
	defer conn.Close()

	decoder := json.NewDecoder(conn)

	sentTimestamps := make(map[string]struct{})

	for {
		select {
		case <-stop:
			return
		default:
			var data DataMeta

			err := decoder.Decode(&data)
			if err != nil {
				if _, ok := err.(*json.SyntaxError); ok {
					continue
				} else {
					errors <- err
					break
				}
			}

			if _, ok := sentTimestamps[data.Time]; !ok {
				data.Time = fmt.Sprintf("%d", time.Now().UTC().Unix())
				dataChannel <- data

				sentTimestamps[data.Time] = struct{}{}

				time.Sleep(100 * time.Millisecond)
			}
		}
	}
}
