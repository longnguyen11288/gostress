package main

import (
	"fmt"
	"log"
	"time"

	"code.google.com/p/go.net/websocket"
)

const BURST_SIZE = 50
const BURST_INTV = 200 * time.Millisecond
const CONNECTIONS = 10000
const SERVICE = 8080
const HOST = "127.0.0.1"
const ORIGIN = "127.0.0.1"

var p = Pool {
	connections: make(map[*Connection]bool),
	subscribe: make(chan *Connection),
	unsubscribe: make(chan *Connection),
}

func client(id int) {
	origin := fmt.Sprintf("http://%s/", ORIGIN)
	target := fmt.Sprintf("ws://%s:%d/", HOST, SERVICE)
	
	ws, err := websocket.Dial(target, "", origin)
	if err != nil {
		log.Print(err)
	}
	
	var conn = &Connection {
		ws: ws,
        }
	p.subscribe <- conn
	for {
		var message string
		err = websocket.Message.Receive(ws, &message)
		if err != nil {
			break
		}
	}
	p.unsubscribe <- conn
}

func flood(connid *int) {
	for {
		for y := 0; y < BURST_SIZE; y++ {
			if *connid >= CONNECTIONS {
				break
			}
			go client(*connid)
			*connid++
		}
		time.Sleep(BURST_INTV)
	}
}

func main() {
	var connid = 0
	
	go p.Dispatch()
	go flood(&connid);
	
	for {
		time.Sleep(1 * time.Second)
		log.Printf("Connections spawned: %d", connid)
		if len(p.connections) <= 0 {
			break
		}
	}
}
