package main

import (
	"fmt"
	"time"
	"flag"
	"math/rand"
	"log"

	"code.google.com/p/go.net/websocket"
)


const BURST_SIZE = 10
const BURST_INTV = 100
const CONNECTIONS = 100
const SERVICE = 8080
const HOST = "127.0.0.1"
const ORIGIN = "127.0.0.1"
const MAX_SECOND = 100000


const debug = false


var p = Pool {
	connections: make(map[*Connection]bool),
	subscribe: make(chan *Connection),
	unsubscribe: make(chan *Connection),
}


func zap(id int, conn *Connection) {
	n := 0
	for {
		time.Sleep(time.Duration(rand.Intn(MAX_SECOND)) * time.Second)
		
		// sending a message...
		if debug { log.Printf("Client %d - Sending %d...", id, n) }
		message := fmt.Sprintf("Client: %d, Zap: %d", id, n)
		err := websocket.Message.Send(conn.ws, message)		
		if err != nil {
			fmt.Printf("Send Error: %s\n", err)
			break
		}
		if debug { log.Printf("Client %d - Sending %d [done]", id, n) }
		
		// waiting for the response...
		if debug { log.Printf("Client %d - Waiting response %d...", id, n) }
		for {
			var response string
			err = websocket.Message.Receive(conn.ws, &response)
			if err != nil {
				fmt.Printf("Receive Error: %s\n", err)
			}
			break
		}
		if debug { log.Printf("Client %d - Waiting response %d [done]", id, n) }
		n++
	}
}

func ready_to_work(id int, conn *Connection) {
	if debug { log.Printf("Client %d ready to work...", id) }
	
	p.subscribe <- conn
	zap(id, conn)
	p.unsubscribe <- conn
	
	conn.ws.Close()
}

func client(id int, host string, service int, origin string) {
	if debug { log.Printf("Creating a new client id=%d to connect at %s:%d...", id, host, service) }
	
	orig := fmt.Sprintf("http://%s/", origin)
	targ := fmt.Sprintf("ws://%s:%d/", host, service)
	
	ws, err := websocket.Dial(targ, "", orig)
	if err != nil {
		fmt.Printf("Dial Error for client: %d: %s\n", id, err)
		return
	}
	var conn = &Connection {
		ws: ws,
        }
	go ready_to_work(id, conn)
}

func simulate(
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
	var connections = flag.Int(
		"connections", 
		CONNECTIONS, 
		"Number of concurent connections")
	var host = flag.String(
		"host", 
		HOST, 
		"Server host")
	var port = 
		flag.Int(
		"port", 
		SERVICE, 
		"Server port")
	var origin = flag.String(
		"origin", 
		HOST, 
		"Client origin")
	var burst_size = flag.Int(
		"burst-size", 
		BURST_SIZE, 
		"Number of concurent connections per bucket send")
	var burst_intv = flag.Int(
		"burst-intv", 
		BURST_INTV, 
		"Interval between each bust of connections")
	
	flag.Parse()
	
	go p.Dispatch()
	
	var connid = 0
	go simulate(&connid, 
		*connections, 
		*host, 
		*port, 
		*origin, 
		*burst_size, 
		*burst_intv);
	
	
	fmt.Printf("Wait for initialization...\n")
	for {
		time.Sleep(200 * time.Millisecond)
		if len(p.connections) >= 1 {
			break
		}
	}

	for {
		fmt.Printf("Spawned: %d, Connected: %d\n", connid, len(p.connections))
		if len(p.connections) <= 0 {
			break
		}
		time.Sleep(1 * time.Second)
	}
}
