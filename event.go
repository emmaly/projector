package projector

import (
	"fmt"
	"time"
)

// Event is emitted when an event or property change occurs
type Event struct {
	EventType  EventType
	Field      string
	ChangeFrom interface{}
	ChangeTo   interface{}
	Timestamp  time.Time
	Properties Properties
}

// EventType is an event type
type EventType string

// Event types
var (
	EventPropertyChange EventType = "propertychange"
)

func (p *Projector) emitEvent(e Event) {
	if e.Timestamp.IsZero() {
		e.Timestamp = time.Now()
	}
	if p.DebugOutput {
		if e.EventType == EventPropertyChange {
			fmt.Printf("CHANGED %s from [%v] to [%v]\n", e.Field, e.ChangeFrom, e.ChangeTo)
		}
	}
	if p.eventChan != nil {
		go func(e Event) {
			p.eventChan <- e
		}(e)
	}
}
