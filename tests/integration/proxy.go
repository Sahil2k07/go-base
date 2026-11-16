package integration

import (
	"io"
	"net"
	"sync"
)

// proxy forwards a fixed local address to whichever container the current test started,
// so the app's load-once config never needs to change.
type proxy struct {
	mu     sync.Mutex
	target string
	conns  map[net.Conn]struct{}
}

func startProxy(addr string) (*proxy, error) {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}

	p := &proxy{conns: map[net.Conn]struct{}{}}

	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				return
			}

			go p.forward(conn)
		}
	}()

	return p, nil
}

func (p *proxy) SetTarget(addr string) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.target = addr
}

// DropConnections closes everything still open, so a finished test cannot leak into the next one.
func (p *proxy) DropConnections() {
	p.mu.Lock()
	defer p.mu.Unlock()

	for conn := range p.conns {
		_ = conn.Close()
	}
}

func (p *proxy) forward(client net.Conn) {
	p.mu.Lock()
	target := p.target
	p.mu.Unlock()

	upstream, err := net.Dial("tcp", target)
	if err != nil {
		_ = client.Close()
		return
	}

	p.mu.Lock()
	p.conns[client] = struct{}{}
	p.conns[upstream] = struct{}{}
	p.mu.Unlock()

	done := make(chan struct{}, 2)

	go func() {
		_, _ = io.Copy(upstream, client)
		done <- struct{}{}
	}()

	go func() {
		_, _ = io.Copy(client, upstream)
		done <- struct{}{}
	}()

	<-done

	_ = client.Close()
	_ = upstream.Close()

	p.mu.Lock()
	delete(p.conns, client)
	delete(p.conns, upstream)
	p.mu.Unlock()
}
