package main

import (
	"fmt"
	"strings"
	"time"
)

const (
	dpms   = "dpms"
	layout = "layout"
	vpn    = "vpn"
	extip  = "extip"
	clock  = "clock"
)

var printOrder = []string{dpms, layout, vpn, extip, clock}

type DataUpdater struct {
	Command  func() string
	Channel  chan string
	Name     string
	Prefix   string
	Data     string
	Interval time.Duration
	Enabled  bool
}

var dataUpdaters = map[string]*DataUpdater{
	dpms: {
		Command:  DPMS,
		Channel:  make(chan string),
		Name:     dpms,
		Data:     "",
		Prefix:   "󰌵",
		Interval: 1 * time.Second,
		Enabled:  true,
	},
	layout: {
		Command:  KeyboardLayout,
		Channel:  make(chan string),
		Name:     layout,
		Data:     "",
		Prefix:   "| 󰌌",
		Interval: 1 * time.Second,
		Enabled:  true,
	},
	vpn: {
		Command:  PIA,
		Channel:  make(chan string),
		Name:     vpn,
		Data:     "",
		Prefix:   "| 󱇱",
		Interval: 10 * time.Second,
		Enabled:  true,
	},
	extip: {
		Command:  ExternalIP,
		Channel:  make(chan string),
		Name:     extip,
		Data:     "",
		Prefix:   "| 󰅟",
		Interval: 600 * time.Second,
		Enabled:  true,
	},
	clock: {
		Command:  Clock,
		Channel:  make(chan string),
		Name:     clock,
		Data:     "",
		Prefix:   "| ",
		Interval: 60 * time.Second,
		Enabled:  true,
	},
}

func updateData(updater *DataUpdater) {
	if !updater.Enabled {
		return
	}

	for {
		data := updater.Command()
		updater.Channel <- data
		time.Sleep(updater.Interval)
	}
}

// nolint: cyclop
func receive(dataUpdaters map[string]*DataUpdater) {
	for {
		select {
		case dataUpdaters[dpms].Data = <-dataUpdaters[dpms].Channel:
		case dataUpdaters[layout].Data = <-dataUpdaters[layout].Channel:
		case dataUpdaters[vpn].Data = <-dataUpdaters[vpn].Channel:
		case dataUpdaters[extip].Data = <-dataUpdaters[extip].Channel:
		case dataUpdaters[clock].Data = <-dataUpdaters[clock].Channel:
		}

		status := ""

		for _, name := range printOrder {
			updater := dataUpdaters[name]

			if name == extip && updater.Data == dataUpdaters[vpn].Data {
				if updater.Data != "" {
					updater.Interval = 3600 * time.Second
				}

				continue
			}

			// status := fmt.Sprintf("󰌵 %s | 󰌌 %s | 󱇱 %s |  %s", dpms, layout, ipaddress, clock)
			status += fmt.Sprintf("%s %s ", updater.Prefix, updater.Data)
		}

		if debug {
			fmt.Println(status)
		} else {
			ExecCommand("xsetroot", []string{"-name", strings.TrimSpace(status)}, false)
		}
	}
}
