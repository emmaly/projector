package main

import (
	"fmt"
	"log"

	"github.com/dustywilson/projector"
)

func main() {
	eventChan := make(chan projector.Event)
	p, err := projector.New("10.10.10.4", eventChan)
	if err != nil {
		log.Fatal(err)
	}
	defer p.Close()
	// p.DebugOutput = true

	for {
		e := <-eventChan
		fmt.Printf("EVENT: %+v\n", e)
	}

	// var i int
	// for {
	// 	fmt.Printf("\n\n%s\n\n%+v\n", time.Now(), p.Properties)
	// 	time.Sleep(time.Second)
	// 	i++
	// 	if i%5 == 0 {
	// 		// p.LanguageEn()
	// 	}
	// }
}
