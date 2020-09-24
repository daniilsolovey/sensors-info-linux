package main

import (
	"fmt"
	"math"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/distatus/battery"
	"github.com/reconquest/pkg/log"
	"github.com/shirou/gopsutil/mem"
	"github.com/sparrc/go-ping"
	"github.com/ssimunic/gosensors"
	wifiname "github.com/yelinaung/wifi-name"
)

const (
	showingTime = "3000"
)

func main() {
	// cpu data:
	sensors, err := gosensors.NewFromSystem()
	var cpuTemp string
	if err != nil {
		log.Error(err)
		cpuTemp = " CPU TEMP: error"
	} else {
		cpuTemp = " CPU TEMP: " + string(
			strings.Split(
				sensors.Chips["coretemp-isa-0000"]["Core 0"], " ",
			)[0],
		)
	}

	var cpuFrequency string
	frequency, err := getCPUFrequency()
	if err != nil {
		log.Error(err)
		cpuFrequency = " CPU FREQUENCY: error"
	} else {
		cpuFrequency = " CPU FREQUENCY: " + frequency + " Mhz"
	}

	// battery data:
	var batteryStatus string
	battery, err := battery.Get(0)
	if err != nil {
		log.Error(err)
		batteryStatus = " BATTERY: error"
	} else {
		batteryStatus = " BATTERY: " + fmt.Sprintf(
			"%.0f", math.Floor(battery.Current/battery.Full*100),
		) + " %" +
			"\n CHARGE STATUS: " + battery.State.String()
	}

	// ping data:
	var pingAVG string
	pinger, err := ping.NewPinger("www.google.com")
	if err != nil {
		log.Error(err)
		pingAVG = " PING: error"
	} else {
		pinger.Count = 1
		pinger.Run()
		stats := pinger.Statistics()
		pingAVG = " PING: " + fmt.Sprintf("%.0f", float64(stats.AvgRtt)/1000000) + " ms"
	}

	// ram data:
	var totalRAM string
	var freeRAM string
	memory, err := mem.VirtualMemory()
	if err != nil {
		log.Error(err)
		totalRAM = " TOTAL RAM: error"
		freeRAM = " FREE RAM: error"
	} else {
		totalRAM = " TOTAL RAM: " + fmt.Sprintf("%.1f", float64(memory.Total)/1000000000) + " GB"
		freeRAM = " FREE RAM: " + fmt.Sprintf("%.1f", float64(memory.Available)/1000000000) + " GB"
	}

	// for date:
	hour, min, sec := time.Now().Clock()
	year, month, day := time.Now().Date()
	date := "TIME: " + strconv.Itoa(hour) + ":" + strconv.Itoa(min) +
		":" + strconv.Itoa(sec) + "\n" + "DATE: " + strconv.Itoa(day) +
		"-" + month.String() + "-" + strconv.Itoa(year)

	// wi-fi name:
	wifiName := " WI-FI: " + wifiname.WifiName()

	// vpn status:
	var vpnStatus string
	status, err := getVPNStatus()
	if err != nil {
		log.Error(err)
		vpnStatus = " NORDVPN: error"
	} else {
		vpnStatus = " NORDVPN: " + status
	}

	var info []string
	info = append(
		info,
		"<span color='#000000' font='21px'>"+date+"</span>",
		"\nNETWORK:",
		"<span color='#0083c9' font='18px'><b>"+pingAVG+"</b></span>",
		"<span color='#0083c9' font='18px'><b>"+wifiName+"</b></span>",
		"<span color='#0083c9' font='18px'><b>"+vpnStatus+"</b></span>",
		"\nSYSTEM:",
		"<span color='#ff0037' font='18px'><b>"+cpuTemp+"</b></span>",
		"<span color='#0026ff' font='18px'><b>"+cpuFrequency+"</b></span>",
		"<span color='#b300ff' font='18px'><b>"+totalRAM+"</b></span>",
		"<span color='#327501' font='18px'><b>"+freeRAM+"</b></span>",
		"\nBATTERY:",
		"<span color='#008783' font='18px'><b>"+batteryStatus+"</b></span>",
	)
	notify := exec.Command(
		"notify-send",
		"-t",
		showingTime,
		"info",
		strings.Join(info, "\n"),
	)
	err = notify.Run()
	if err != nil {
		log.Error(err)
	}
}
