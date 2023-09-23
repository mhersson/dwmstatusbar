package dwmstatusbar

import (
	"log"
	"os"

	"github.com/spf13/afero"
)

var (
	debug         = false
	Log           = log.New(os.Stdout, "", 0)
	ExternalIPURL = "https://icanhazip.com"
	Fsys          afero.Fs
)

func Run(setDebug bool) {
	debug = setDebug

	Fsys = afero.NewOsFs()

	for _, updater := range dataUpdaters {
		go updateData(updater)
	}

	go receive(dataUpdaters)

	select {}
}
