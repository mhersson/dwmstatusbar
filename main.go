package main

import (
	"os"
	"time"
)

var debug = false

func main() {
	if len(os.Args) > 1 {
		debug = true
	}

	dataUpdaters := []DataUpdater{
		{Command: DPMS, Channel: make(chan string), Interval: 10 * time.Second},
		{Command: KeyboardLayout, Channel: make(chan string), Interval: 1 * time.Second},
		{Command: ExternalIP, Channel: make(chan string), Interval: 600 * time.Second},
		{Command: Clock, Channel: make(chan string), Interval: 60 * time.Second},
	}

	for _, updater := range dataUpdaters {
		go updateData(updater)
	}

	go receive(dataUpdaters)

	select {}
}
