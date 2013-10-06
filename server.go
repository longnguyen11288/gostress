package main

import (
	"fmt"
	"net/http"
	"log"
	"time"
	"code.google.com/p/go.net/websocket"
)

var p = Pool {
	connections: make(map[*Connection]bool),
	subscribe: make(chan *Connection),
	unsubscribe: make(chan *Connection),
}

func feed(conn *Connection, n int) {
	message := fmt.Sprintf("Chunk: %d", n)
	err := websocket.Message.Send(conn.ws, message)
	if err == nil {
		time.Sleep(10 * time.Second)
		feed(conn, n+1)
	}
}

func socketHandler(ws *websocket.Conn) {
	var conn = &Connection {
		ws: ws,
	}
	
	p.subscribe <- conn
	feed(conn, 1)
	p.unsubscribe <- conn
}

func displayStats() {
	for {
                log.Printf("Clients: %d", len(p.connections))
		time.Sleep(1 * time.Second)
        }
}

func main() {
	go p.Dispatch()
	go displayStats()
	
	http.Handle("/", websocket.Handler(socketHandler))
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}
