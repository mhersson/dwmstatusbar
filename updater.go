package main

import (
	"fmt"
	"time"
)

type DataUpdater struct {
	Command  func() string
	Channel  chan string
	Interval time.Duration
}

func updateData(updater DataUpdater) {
	for {
		data := updater.Command()
		updater.Channel <- data
		time.Sleep(updater.Interval)
	}
}

func receive(dataUpdaters []DataUpdater) {
	var dpms, layout, ipaddress, clock string

	for {
		select {
		case data := <-dataUpdaters[0].Channel:
			dpms = data
		case data := <-dataUpdaters[1].Channel:
			layout = data
		case data := <-dataUpdaters[2].Channel:
			ipaddress = data
		case data := <-dataUpdaters[3].Channel:
			clock = data
		}

		status := fmt.Sprintf("󰌵 %s | 󰌌 %s | 󱇱 %s |  %s", dpms, layout, ipaddress, clock)

		if debug {
			fmt.Println(status)
		} else {
			ExecCommand("xsetroot", []string{"-name", status}, false)
		}
	}
}
