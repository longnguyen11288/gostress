package main

import (
        "code.google.com/p/go.net/websocket"
)

type Connection struct {
        ws *websocket.Conn
}

type Pool struct {
        connections map[*Connection]bool
	
        subscribe chan *Connection
        unsubscribe chan *Connection
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