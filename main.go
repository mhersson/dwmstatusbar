package main

import (
	"os"

	"github.com/mhersson/dwmstatusbar/internal/dwmstatusbar"
)

var debug = false

func main() {
	if len(os.Args) > 1 {
		debug = true
	}

	dwmstatusbar.Run(debug)
}
