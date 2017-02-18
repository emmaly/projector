package projector

import (
	"bufio"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

// When using Wireshark, `(data.data contains 05:00:06:00:00:03:00)` works well for grabbing commands

// Projector is a projector
type Projector struct {
	Conn          net.Conn
	IP            string
	CommandPrefix string
	Properties    *Properties
	eventChan     chan Event
	DebugOutput   bool
	messageChan   chan []string
	done          chan struct{}
	wg            sync.WaitGroup
	connected     bool
}

// Errors
var (
	ErrNotConnected = errors.New("not connected")
)

// New creates a new Projector and connects to it
func New(ip string, eventChan chan Event, connectTimeout time.Duration) (*Projector, error) {
	p := &Projector{
		IP:            ip,
		CommandPrefix: "05000600000300",
		Properties:    &Properties{},
		eventChan:     eventChan,
		messageChan:   make(chan []string),
		done:          make(chan struct{}),
	}
	go p.listen()
	return p, p.Connect(connectTimeout)
}

// Close disconnects from a Projector.
// This blocks until it has finished closing and cleaning up.
func (p *Projector) Close() {
	close(p.done)
	p.Conn.Close()
	p.wg.Wait()
}

// Connect (re)connects the Projector
// This blocks until it either succeeds or hits the provided expiration
func (p *Projector) Connect(timeout time.Duration) error {
	p.wg.Add(1)
	defer p.wg.Done()
	if p.connected {
		return nil
	}
	var err error
	p.Conn, err = net.DialTimeout("tcp", p.IP+":41794", timeout)
	if err != nil {
		p.connected = false
		return err
	}
	p.connected = true
	go p.receiveMessages()
	p.RefreshProperties()
	return nil
}

func (p *Projector) listen() {
	p.wg.Add(1)
	defer p.wg.Done()
	for {
		select {
		case <-p.done:
			return
		case m := <-p.messageChan:
			for i, v := range m {
				fmt.Printf("IN: [%d] [%s]\n", i, v)
			}
		}
	}
}

func (p *Projector) receiveMessages() {
	p.wg.Add(1)
	defer p.wg.Done()
	b := bufio.NewReader(p.Conn)
	for {
		select {
		case <-p.done:
			return
		default:
			r, err := b.ReadBytes(byte(5))
			if err != nil {
				log.Fatal(err)
				p.connected = false
				return
			}
			p.messageChan <- p.ParseProperty(r)
		}
	}
}

// RawCommand sends a raw command to the projector.
// You probably should not use this unless you've got no other option.
// Commands are provided in hex.
func (p *Projector) RawCommand(h string) error {
	buf, err := hex.DecodeString(h)
	if err != nil {
		return err
	}
	if p.Conn == nil {
		return ErrNotConnected
	}
	_, err = p.Conn.Write(buf)
	return err
}

// Command sends a command to the projector.
// You probably should not use this unless you've got no other option.
// Commands are provided in hex.
func (p *Projector) Command(h string) error {
	return p.RawCommand(p.CommandPrefix + h)
}

// ReconnectIfNeeded will reconnect+retry if the error is a (dis)connection-related error,
// and will otherwise simply passthorugh the error.
func (p *Projector) ReconnectIfNeeded(fn func() error) error {
	err := fn()

	if err == ErrNotConnected { // this is not the only error that should cause a reconnection...  FIXME
		log.Println("[debug] Reconnecting to projector...")
		err = p.Connect(time.Second * 5)
		log.Printf("[debug] Reconnection:[%s]\n", err)
		if err != nil {
			return err
		}

		// This is the retry.  If the connection is still down, it just stays down and returns the error.
		// If they want it to retry *again*, they can just wrap it however they want...
		err = fn()
	}

	log.Printf("[debug] RIN:[%s]\n", err)
	return err
}
