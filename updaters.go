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
	Data     string
	Prefix   string
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
		Interval: 10 * time.Second,
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
		Prefix:   "| 󱇱",
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

			if name == extip && dataUpdaters[vpn].Enabled {
				continue
			}

			if name == vpn && updater.Enabled {
				if dataUpdaters[extip].Enabled {
					extipUpdater := dataUpdaters[extip]

					// If vpn is enabled only check the external ip against icanhazip.com every hour
					extipUpdater.Interval = 3600 * time.Second

					// Check to verify that the external ip address reported by the vpn
					// is actually the external ip address shown to the world. The use
					// case appeared with a bug in PIA that caused the routing table to
					// be renamed with a .pacsave extension on an Archlinux update,
					// leaving PIA to report connected, but the connection was never
					// used. Since the extip interval is set to 1 hour when vpn is
					// enabled, this will stay on IP mismatch until a new external check
					// is made in one hour
					// TBD if this should be ignored and the extip should be disabled if
					// vpn is enabled
					if extipUpdater.Data != updater.Data && (extipUpdater.Data != "" && updater.Data != "") {
						status += fmt.Sprintf("%s IP mismatch ", updater.Prefix)

						continue
					}
				}
			}

			status += fmt.Sprintf("%s %s ", updater.Prefix, updater.Data)
		}
		// status := fmt.Sprintf("󰌵 %s | 󰌌 %s | 󱇱 %s |  %s", dpms, layout, ipaddress, clock)

		if debug {
			fmt.Println(status)
		} else {
			ExecCommand("xsetroot", []string{"-name", strings.TrimSpace(status)}, false)
		}
	}
}
