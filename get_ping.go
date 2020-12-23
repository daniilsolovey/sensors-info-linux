package main

import (
	"fmt"

	"github.com/reconquest/pkg/log"
	"github.com/sparrc/go-ping"
)

func getPing() string {
	pinger, err := ping.NewPinger("www.google.com")
	var pingAVG string
	if err != nil {
		log.Error(err)
		pingAVG = " PING: error"
	} else {
		pinger.Count = 1
		pinger.Run()
		stats := pinger.Statistics()
		pingAVG = " PING: " + fmt.Sprintf("%.0f", float64(stats.AvgRtt)/1000000) + " ms"
	}

	return pingAVG
}
