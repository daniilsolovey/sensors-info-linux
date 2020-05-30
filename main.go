package main

import (
	"fmt"
	"math"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/capnm/sysinfo"
	"github.com/distatus/battery"
	"github.com/reconquest/pkg/log"
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
	if err != nil {
		log.Fatal(err)
	}

	cpuTemp := "cpu temp: " + string(
		strings.Split(
			sensors.Chips["coretemp-isa-0000"]["Core 0"], " ",
		)[0],
	)

	// battery data:
	battery, err := battery.Get(0)
	if err != nil {
		log.Fatal(err)
	}

	batteryStatus := "battery: " + fmt.Sprintf(
		"%.0f", math.Floor(battery.Current/battery.Full*100),
	) +
		" %"

	// ping data:
	pinger, err := ping.NewPinger("www.google.com")
	if err != nil {
		log.Fatal(err)
	}

	pinger.Count = 1
	pinger.Run()
	stats := pinger.Statistics()
	pingAVG := "ping: " + stats.AvgRtt.String()

	// ram data:
	sysInfo := sysinfo.Get()
	totalRAM := "total ram: " + strconv.FormatUint(sysInfo.TotalRam, 10)
	freeRAM := "free ram: " + strconv.FormatUint(sysInfo.FreeRam, 10)

	// for date:
	hour, min, sec := time.Now().Clock()
	year, month, day := time.Now().Date()
	date := "time: " + strconv.Itoa(hour) + ":" + strconv.Itoa(min) +
		":" + strconv.Itoa(sec) + "\n" + "date: " + strconv.Itoa(day) +
		"-" + month.String() + "-" + strconv.Itoa(year)

	// wi-fi name:
	wifiName := "wi-fi: " + wifiname.WifiName()

	var info []string
	info = append(info, date, "", totalRAM, freeRAM, "", pingAVG, wifiName, "", cpuTemp, batteryStatus)
	notify := exec.Command("notify-send", "-t", showingTime, "info", strings.Join(info, "\n"))
	err = notify.Run()
	if err != nil {
		log.Fatal(err)
	}

}
