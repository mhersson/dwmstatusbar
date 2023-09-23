package dwmstatusbar

import (
	"log"
	"os"
)

var (
	debug = false
	Log   = log.New(os.Stdout, "", 0)
)

func Run(setDebug bool) {
	debug = setDebug

	for _, updater := range dataUpdaters {
		go updateData(updater)
	}

	go receive(dataUpdaters)

	select {}
}
