package projector

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"log"
	"net"
	"sync"
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
	done          chan struct{}
	wg            sync.WaitGroup
}

// New creates a new Projector and connects to it
func New(ip string, eventChan chan Event) (*Projector, error) {
	p := &Projector{
		IP:            ip,
		CommandPrefix: "05000600000300",
		Properties:    &Properties{},
		eventChan:     eventChan,
		done:          make(chan struct{}),
	}
	var err error
	p.Conn, err = net.Dial("tcp", ip+":41794")
	if err != nil {
		return nil, err
	}
	go p.listen()
	p.RefreshProperties()
	return p, nil
}

// Close disconnects from a Projector.
// This blocks until it has finished closing and cleaning up.
func (p *Projector) Close() {
	close(p.done)
	p.Conn.Close()
	p.wg.Wait()
}

func (p *Projector) listen() {
	p.wg.Add(1)
	defer p.wg.Done()
	in := make(chan []string)
	go p.receiveMessages(in)
	for {
		select {
		case <-p.done:
			return
		case m := <-in:
			for i, v := range m {
				fmt.Printf("IN: [%d] [%s]\n", i, v)
			}
		}
	}
}

func (p *Projector) receiveMessages(in chan []string) {
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
			}
			in <- p.ParseProperty(r)
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
	_, err = p.Conn.Write(buf)
	return err
}

// Command sends a command to the projector.
// You probably should not use this unless you've got no other option.
// Commands are provided in hex.
func (p *Projector) Command(h string) error {
	return p.RawCommand(p.CommandPrefix + h)
}
