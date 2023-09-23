package dwmstatusbar

import (
	"fmt"
	"strings"
	"time"
)

const (
	battery = "battery"
	clock   = "clock"
	dpms    = "dpms"
	extip   = "extip"
	layout  = "layout"
	vpn     = "vpn"
	xset    = "xset"
)

var printOrder = []string{dpms, layout, vpn, extip, battery, clock}

type DataUpdater struct {
	Command  func(string) string
	Channel  chan string
	Name     string
	Prefix   string
	Data     string
	OldData  string
	Parent   string
	Interval time.Duration
	Enabled  bool
}

var dataUpdaters = map[string]*DataUpdater{
	xset: {
		Command:  Xset,
		Channel:  make(chan string),
		Name:     xset,
		Data:     "",
		Interval: 1 * time.Second,
		Enabled:  true,
	},
	dpms: {
		Command:  DPMS,
		Channel:  make(chan string),
		Name:     dpms,
		Data:     "",
		Prefix:   "󰌵",
		Interval: 1 * time.Second,
		Enabled:  true,
		Parent:   xset,
	},
	layout: {
		Command:  KeyboardLayout,
		Channel:  make(chan string),
		Name:     layout,
		Data:     "",
		Prefix:   "| 󰌌",
		Interval: 1 * time.Second,
		Enabled:  true,
		Parent:   xset,
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
	battery: {
		Command:  Battery,
		Channel:  make(chan string),
		Name:     battery,
		Data:     "",
		Prefix:   "|  ",
		Interval: 60 * time.Second,
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
	for {
		if !updater.Enabled {
			return
		}

		var data string

		if updater.Parent != "" {
			data = updater.Command(dataUpdaters[updater.Parent].Data)
		} else {
			data = updater.Command("")
		}

		if updater.OldData != data {
			updater.Channel <- data
			updater.OldData = data
		}

		time.Sleep(updater.Interval)
	}
}

// nolint: cyclop
func receive(dataUpdaters map[string]*DataUpdater) {
	for {
		select {
		case dataUpdaters[battery].Data = <-dataUpdaters[battery].Channel:
		case dataUpdaters[clock].Data = <-dataUpdaters[clock].Channel:
		case dataUpdaters[dpms].Data = <-dataUpdaters[dpms].Channel:
		case dataUpdaters[extip].Data = <-dataUpdaters[extip].Channel:
		case dataUpdaters[layout].Data = <-dataUpdaters[layout].Channel:
		case dataUpdaters[vpn].Data = <-dataUpdaters[vpn].Channel:
		case dataUpdaters[xset].Data = <-dataUpdaters[xset].Channel:
		}

		status := ""

		for _, name := range printOrder {
			updater := dataUpdaters[name]

			if name == battery && updater.Data == "No Battery" {
				updater.Enabled = false

				continue
			}

			if name == extip && updater.Data == dataUpdaters[vpn].Data {
				if updater.Data != "" {
					updater.Interval = 3600 * time.Second
				}

				continue
			}

			// status := fmt.Sprintf("󰌵 %s | 󰌌 %s | 󱇱 %s |  %s", dpms, layout, ipaddress, clock)
			if updater.Data != "" {
				status += fmt.Sprintf("%s %s ", updater.Prefix, updater.Data)
			}
		}

		if debug {
			fmt.Println(status)
		} else {
			ExecCommand("xsetroot", []string{"-name", strings.TrimSpace(status)}, false)
		}
	}
}
