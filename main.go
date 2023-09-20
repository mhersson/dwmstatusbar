package main

import (
	"os"
)

var debug = false

func main() {
	if len(os.Args) > 1 {
		debug = true
	}

	for _, updater := range dataUpdaters {
		go updateData(updater)
	}

	go receive(dataUpdaters)

	select {}
}
