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
		cpuTemp = "cpu temp: error"
	} else {
		cpuTemp = "cpu temp: " + string(
			strings.Split(
				sensors.Chips["coretemp-isa-0000"]["Core 0"], " ",
			)[0],
		)
	}

	var cpuFrequency string
	frequency, err := getCPUFrequency()
	if err != nil {
		log.Error(err)
		cpuFrequency = "cpu frequency: error"
	} else {
		cpuFrequency = "cpu frequency: " + frequency + " Mhz"
	}

	// battery data:
	var batteryStatus string
	battery, err := battery.Get(0)
	if err != nil {
		log.Error(err)
		batteryStatus = "battery: error"
	} else {
		batteryStatus = "battery: " + fmt.Sprintf(
			"%.0f", math.Floor(battery.Current/battery.Full*100),
		) +
			" %"
	}

	// ping data:
	var pingAVG string
	pinger, err := ping.NewPinger("www.google.com")
	if err != nil {
		log.Error(err)
		pingAVG = "ping: error"
	} else {
		pinger.Count = 1
		pinger.Run()
		stats := pinger.Statistics()
		pingAVG = "ping: " + fmt.Sprintf("%.0f", float64(stats.AvgRtt)/1000000) + " ms"
	}

	// ram data:
	var totalRAM string
	var freeRAM string
	memory, err := mem.VirtualMemory()
	if err != nil {
		log.Error(err)
		totalRAM = "total ram: error"
		freeRAM = "free ram: error"
	} else {
		totalRAM = "total ram: " + fmt.Sprintf("%.1f", float64(memory.Total)/1000000000) + " GB"
		freeRAM = "free ram: " + fmt.Sprintf("%.1f", float64(memory.Free)/1000000000) + " GB"
	}

	// for date:
	hour, min, sec := time.Now().Clock()
	year, month, day := time.Now().Date()
	date := "time: " + strconv.Itoa(hour) + ":" + strconv.Itoa(min) +
		":" + strconv.Itoa(sec) + "\n" + "date: " + strconv.Itoa(day) +
		"-" + month.String() + "-" + strconv.Itoa(year)

	// wi-fi name:
	wifiName := "wi-fi: " + wifiname.WifiName()

	var info []string
	info = append(info, date, "", totalRAM, freeRAM, "", pingAVG, wifiName, "", cpuTemp, cpuFrequency, "", batteryStatus)
	notify := exec.Command("notify-send", "-t", showingTime, "info", strings.Join(info, "\n"))
	err = notify.Run()
	if err != nil {
		log.Error(err)
	}
}
