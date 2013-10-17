package main

import (
	"fmt"
	"net/http"
	"log"
	"time"
	"flag"
	
	"code.google.com/p/go.net/websocket"
)

const SERVICE = 8080
const HOST = ""
const LISTEN_INTV = 400 * time.Millisecond

var p = Pool {
	connections: make(map[*Connection]bool),
	subscribe: make(chan *Connection),
	unsubscribe: make(chan *Connection),
}

const DEBUG = true
func debug(msg string, args ...interface{}) {
	if DEBUG == true {
		log.Print(msg, args)
	}
}

func serv(conn *Connection) {
	for {
		// waiting for a client message
		var receive string
		err := websocket.Message.Receive(conn.ws, &receive)
		if err != nil {
			fmt.Printf("Receive Error: %s\n", err)
			break
		}
		log.Printf("Receives Message: %s", receive)
		
		// message receive, responding to the client.
		message := fmt.Sprintf("Response: [%s]", receive)
		err = websocket.Message.Send(conn.ws, message)
		if err != nil {
			fmt.Printf("Send Error: %s\n", err)
			break
		}
		log.Printf("Response sent.")

		// A client sends a message in minimun every 1s.
		// so we can wait before check for a new message.
		time.Sleep(LISTEN_INTV)
	}
}

func socketHandler(ws *websocket.Conn) {
	var conn = &Connection {
		ws: ws,
	}
	
	p.subscribe <- conn
	serv(conn)
	p.unsubscribe <- conn
}

func displayStats() {
	for {
                fmt.Printf("connections: %d\n", len(p.connections))
		time.Sleep(1 * time.Second)
        }
}

func main() {
	var host = flag.String("host", HOST, "Server listen host")
	var port = flag.Int("port", SERVICE, "Server listen port")
	
	flag.Parse()

	go p.Dispatch()
	go displayStats()
	
	http.Handle("/", websocket.Handler(socketHandler))
	err := http.ListenAndServe(fmt.Sprintf("%s:%d", *host, *port), nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}
