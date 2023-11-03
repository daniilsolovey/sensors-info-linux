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
	"github.com/ssimunic/gosensors"
	wifiname "github.com/yelinaung/wifi-name"
)

const (
	showingTime = "3000"
	timeFormat  = "15:04:01"
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
	if err != nil && !strings.Contains(fmt.Sprint(err), "State:Invalid state `Not charging") {
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
	go func() {
		pingAVG = getPing()
	}()
	time.Sleep(1000 * time.Millisecond)

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

	moscowTime, localTime, err := getLocalAndMoscowTime()
	if err != nil {
		log.Error(err)
	}

	// for date:
	year, month, day := localTime.Date()
	date := "DATE LOCAL: " + strconv.Itoa(day) +
		"-" + month.String() + "-" + strconv.Itoa(year)

	// for local date:
	localTimeResult := "TIME LOCAL: " + localTime.Format(timeFormat)
	// for moscow date:
	moscowTimeResult := "Time Moscow: " + moscowTime.Format(timeFormat)
	// wi-fi name:
	wifiName := " WI-FI: " + wifiname.WifiName()

	// vpn status:
	var vpnStatus string
	status, err := getCommonVPNStatus()
	if err != nil {
		log.Error(err)
		vpnStatus = " VPN: error"
	} else {
		vpnStatus = " VPN: " + status
	}

	var info []string
	info = append(
		info,
		"<span color='#B22222' font='17px'><b>"+localTimeResult+"</b></span>",
		"<span color='#B22222' font='17px'><b>"+date+"</b></span>",
		"<span color='#B22222' font='17px'><b>"+moscowTimeResult+"</b></span>",
		"<span color='#0083c9' font='17px'><b>"+pingAVG+"</b></span>",
		"<span color='#0083c9' font='17px'><b>"+wifiName+"</b></span>",
		"<span color='#0083c9' font='17px'><b>"+vpnStatus+"</b></span>",
		"<span color='#0026ff' font='17px'><b>"+cpuTemp+"</b></span>",
		"<span color='#0026ff' font='17px'><b>"+cpuFrequency+"</b></span>",
		"<span color='#32CD32' font='17px'><b>"+batteryStatus+"</b></span>",
		"<span color='#0026ff' font='17px'><b>"+totalRAM+"</b></span>",
		"<span color='#0026ff' font='17px'><b>"+freeRAM+"</b></span>",
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
