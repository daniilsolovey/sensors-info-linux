package main

import (
	"fmt"
	"math"
	"os/exec"
	"strings"
	"time"

	"github.com/distatus/battery"
	"github.com/reconquest/pkg/log"
	"github.com/shirou/gopsutil/mem"
	"github.com/ssimunic/gosensors"
	wifiname "github.com/yelinaung/wifi-name"
)

const (
	showingTime = "5000"
	timeFormat  = "15:04"
)

func main() {
	// CPU temperature
	sensors, err := gosensors.NewFromSystem()

	cpuTemp := "error"

	if err != nil {
		log.Error(err)
	} else {
		cpuTemp = strings.Split(
			sensors.Chips["coretemp-isa-0000"]["Core 0"],
			" ",
		)[0]
	}

	// CPU frequency
	cpuFrequency := "error"

	frequency, err := getCPUFrequency()
	if err != nil {
		log.Error(err)
	} else {
		cpuFrequency = frequency
	}

	// Battery
	batteryPercent := "error"
	batteryState := "error"

	bat, err := battery.Get(0)
	if err != nil && !strings.Contains(
		fmt.Sprint(err),
		"State:Invalid state `Not charging",
	) {
		log.Error(err)
	} else {
		batteryPercent = fmt.Sprintf(
			"%.0f%%",
			math.Floor(bat.Current/bat.Full*100),
		)

		batteryState = bat.State.String()
	}

	// Ping
	var pingAVG string

	go func() {
		pingAVG = getPing()
	}()

	time.Sleep(1000 * time.Millisecond)

	// RAM
	totalRAM := "error"
	usedRAM := "error"
	ramUsage := "error"

	memory, err := mem.VirtualMemory()
	if err != nil {
		log.Error(err)
	} else {
		used := memory.Total - memory.Available
		usedPercent := float64(used) / float64(memory.Total) * 100

		totalRAM = fmt.Sprintf(
			"%.1f",
			float64(memory.Total)/1024/1024/1024,
		)

		usedRAM = fmt.Sprintf(
			"%.1f",
			float64(used)/1024/1024/1024,
		)

		ramUsage = fmt.Sprintf("%.0f%%", usedPercent)
	}

	ramValue := fmt.Sprintf(
		"(%s) %s(usedGB)\n%s%s(totalGB)",
		ramUsage,
		usedRAM,
		strings.Repeat(" ", 14+len(fmt.Sprintf("(%s) ", ramUsage))),
		totalRAM,
	)

	// Time
	localTime := time.Now()

	moscowLocation, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		log.Error(err)
	}

	moscowTime := localTime.In(moscowLocation)

	date := localTime.Format("02 January 2006")

	// Wi-Fi
	wifi := wifiname.WifiName()
	if wifi == "" {
		wifi = "disconnected"
	}

	// VPN
	vpn := "error"

	status, err := getCommonVPNStatus()
	if err != nil {
		log.Error(err)
	} else {
		vpn = status
	}

	info := []string{
		"<span foreground='#8e44ad' size='large'><b>◷ TIME</b></span>",
		fmt.Sprintf(
			"<span foreground='#222222'>  %-11s <b>%s</b></span>",
			"Local",
			localTime.Format(timeFormat),
		),
		fmt.Sprintf(
			"<span foreground='#222222'>  %-11s <b>%s</b></span>",
			"Date",
			date,
		),
		fmt.Sprintf(
			"<span foreground='#222222'>  %-11s <b>%s</b></span>",
			"Moscow",
			moscowTime.Format(timeFormat),
		),

		"",

		"<span foreground='#0083a8' size='large'><b>◉ NETWORK</b></span>",
		fmt.Sprintf(
			"<span foreground='#222222'>  %-11s <b>%s</b></span>",
			"Wi-Fi",
			wifi,
		),
		fmt.Sprintf(
			"<span foreground='#222222'>  %-11s <b>%s</b></span>",
			"Ping",
			pingAVG,
		),
		fmt.Sprintf(
			"<span foreground='#222222'>  %-11s <b>%s</b></span>",
			"VPN",
			vpn,
		),

		"",

		"<span foreground='#356aa0' size='large'><b>⚙ SYSTEM</b></span>",
		fmt.Sprintf(
			"<span foreground='#222222'>  %-11s <b>%s</b></span>",
			"CPU temp.",
			cpuTemp,
		),
		fmt.Sprintf(
			"<span foreground='#222222'>  %-11s <b>%s</b></span>",
			"CPU freq.",
			cpuFrequency,
		),
		fmt.Sprintf(
			"<span foreground='#222222'>  %-11s <b>%s</b></span>",
			"RAM",
			ramValue,
		),

		"",

		"<span foreground='#2e7d32' size='large'><b>⚡ POWER</b></span>",
		fmt.Sprintf(
			"<span foreground='#222222'>  %-11s <b>%s</b></span>",
			"Battery",
			batteryPercent,
		),
		fmt.Sprintf(
			"<span foreground='#222222'>  %-11s <b>%s</b></span>",
			"Status",
			batteryState,
		),
	}

	notify := exec.Command(
		"notify-send",
		"-t",
		showingTime,
		"System info",
		strings.Join(info, "\n"),
	)

	if err := notify.Run(); err != nil {
		log.Error(err)
	}
}
