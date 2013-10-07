package main

import (
	"fmt"
	"log"
	"time"
	"flag"

	"code.google.com/p/go.net/websocket"
)

const BURST_SIZE = 50
const BURST_INTV = 500
const CONNECTIONS = 20000
const SERVICE = 8080
const HOST = "127.0.0.1"
const ORIGIN = "127.0.0.1"

var p = Pool {
	connections: make(map[*Connection]bool),
	subscribe: make(chan *Connection),
	unsubscribe: make(chan *Connection),
}


func listen(conn *Connection) {
	p.subscribe <- conn
	for {
		var message string
		err := websocket.Message.Receive(conn.ws, &message)
		if err != nil {
			log.Print("Receive: ", err)
			break
		}
	}
	p.unsubscribe <- conn
	conn.ws.Close()
}

func client(id int, host string, service int, origin string) {
	orig := fmt.Sprintf("http://%s/", origin)
	targ := fmt.Sprintf("ws://%s:%d/", host, service)
	
	ws, err := websocket.Dial(targ, "", orig)
	if err != nil {
		log.Print("Dial: ", err)
		return
	}
	var conn = &Connection {
		ws: ws,
        }
	go listen(conn)
}

func flood(
	connid *int, 
	connections int, 
	host string, 
	service int, 
	origin string, 
	burst_size int, 
	burst_intv int) {
	for {
		for y := 0; y < burst_size; y++ {
			if *connid >= connections {
				break
			}
			client(*connid, host, service, origin)
			*connid++
		}
		time.Sleep(time.Duration(burst_intv) * time.Millisecond)
	}
}

func main() {
	var connections = flag.Int("connections", CONNECTIONS, "Number of concurent connections")
	var host = flag.String("host", HOST, "Server host")
	var port = flag.Int("port", SERVICE, "Server port")
	var origin = flag.String("origin", HOST, "Client origin")
	var burst_size = flag.Int("burst-size", BURST_SIZE, "Number of concurent connections per bucket send")
	var burst_intv = flag.Int("burst-intv", BURST_INTV, "Interval between each bust of connections")
	
	flag.Parse()

	var connid = 0
	go p.Dispatch()
	go flood(&connid, *connections, *host, *port, *origin, *burst_size, *burst_intv);


	log.Printf("Wait for initialization...")
	for {
		if len(p.connections) >= 1 {
			break
		}
		time.Sleep(200 * time.Millisecond)
	}
	
	for {
		time.Sleep(1 * time.Second)
		log.Printf("Connections spawned: %d", connid)
		if len(p.connections) <= 0 {
			break
		}
	}
}
