package main

import (
	"sync/atomic"
        "code.google.com/p/go.net/websocket"
)

type Connection struct {
        ws *websocket.Conn
}

type Pool struct {
        connections map[*Connection]bool
	
        subscribe chan *Connection
        unsubscribe chan *Connection
	
	// stats
	nrecv uint64
	nsend uint64
}

func (p *Pool) Dispatch() {
        for {
                select {
                case c := <- p.subscribe:
                        p.connections[c] = true;
                case c := <- p.unsubscribe:
                        delete(p.connections, c);
                }
        }
}

func (p *Pool) IncrRecv() {
	atomic.AddUint64(&p.nrecv, 1)
}

func (p *Pool) IncrSend() {
	atomic.AddUint64(&p.nsend, 1)
}
