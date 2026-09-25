package main

import (
	"fmt"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
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

	// Fixed notification ID: dunst replaces the existing window instead of
	// stacking a new one to the right.
	notifyID = "991337"
	stackTag = "sensors-info"

	// Minimal interval between two notifications. Protects dunst from being
	// flooded when the hotkey is held down.
	minInterval = 500 * time.Millisecond
)

func main() {
	unlock, ok := acquireRunGuard()
	if !ok {
		return
	}
	defer unlock()

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

	// The kernel reports "Not charging" when the charge limit is
	// reached, which the battery library does not understand.
	if state := readBatterySysfs("status"); state != "" {
		batteryState = state
	}

	chargeStart := getChargeThreshold("charge_control_start_threshold")
	chargeEnd := getChargeThreshold("charge_control_end_threshold")

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

	// Power profile
	powerProfile := "error"

	profile, err := getPowerProfile()
	if err != nil {
		log.Error(err)
	} else {
		powerProfile = profile
	}

	info := []string{
		"<span foreground='#8e44ad' size='large'><b>◷ TIME</b></span>",
		fmt.Sprintf(
			"<span foreground='#cdd6f4'>  %-11s <b>%s</b></span>",
			"Local",
			localTime.Format(timeFormat),
		),
		fmt.Sprintf(
			"<span foreground='#cdd6f4'>  %-11s <b>%s</b></span>",
			"Date",
			date,
		),
		fmt.Sprintf(
			"<span foreground='#cdd6f4'>  %-11s <b>%s</b></span>",
			"Moscow",
			moscowTime.Format(timeFormat),
		),

		"",

		"<span foreground='#0083a8' size='large'><b>◉ NETWORK</b></span>",
		fmt.Sprintf(
			"<span foreground='#cdd6f4'>  %-11s <b>%s</b></span>",
			"Wi-Fi",
			wifi,
		),
		fmt.Sprintf(
			"<span foreground='#cdd6f4'>  %-11s <b>%s</b></span>",
			"Ping",
			pingAVG,
		),
		fmt.Sprintf(
			"<span foreground='#cdd6f4'>  %-11s <b>%s</b></span>",
			"VPN",
			vpn,
		),

		"",

		"<span foreground='#356aa0' size='large'><b>⚙ SYSTEM</b></span>",
		fmt.Sprintf(
			"<span foreground='#cdd6f4'>  %-11s <b>%s</b></span>",
			"CPU temp.",
			cpuTemp,
		),
		fmt.Sprintf(
			"<span foreground='#cdd6f4'>  %-11s <b>%s</b></span>",
			"CPU freq.",
			cpuFrequency,
		),
	}

	info = append(info,
		fmt.Sprintf(
			"<span foreground='#cdd6f4'>  %-11s <b>%s</b></span>",
			"RAM",
			ramValue,
		),
	)

	if fanSpeeds, ok := getFanSpeeds(); ok {
		info = append(info, fmt.Sprintf(
			"<span foreground='#cdd6f4'>  %-11s <b>%s</b></span>",
			"Fans",
			fanSpeeds,
		))
	}

	info = append(info,
		"",

		"<span foreground='#2e7d32' size='large'><b>⚡ POWER</b></span>",
		fmt.Sprintf(
			"<span foreground='#cdd6f4'>  %-11s <b>%s</b></span>",
			"Profile",
			powerProfile,
		),
		fmt.Sprintf(
			"<span foreground='#cdd6f4'>  %-11s <b>%s - %s</b></span>",
			"Charge lim",
			chargeStart,
			chargeEnd,
		),
		fmt.Sprintf(
			"<span foreground='#cdd6f4'>  %-11s <b>%s</b></span>",
			"Status",
			batteryState,
		),
		fmt.Sprintf(
			"<span foreground='#cdd6f4'>  %-11s <b>%s</b></span>",
			"Charge lvl",
			batteryPercent,
		),
	)

	notify := exec.Command(
		"notify-send",
		"-t", showingTime,
		"-r", notifyID,
		"-h", "string:x-dunst-stack-tag:"+stackTag,
		"System info",
		strings.Join(info, "\n"),
	)

	if err := notify.Run(); err != nil {
		log.Error(err)
	}
}

// acquireRunGuard makes sure only one instance runs at a time and that
// notifications are not sent more often than minInterval. Returns false when
// this invocation should be skipped.
func acquireRunGuard() (func(), bool) {
	dir := os.Getenv("XDG_RUNTIME_DIR")
	if dir == "" {
		dir = os.TempDir()
	}

	path := filepath.Join(dir, "sensors-info.lock")

	file, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0o600)
	if err != nil {
		log.Error(err)
		return func() {}, true
	}

	// Another instance is running right now: skip.
	if err := syscall.Flock(int(file.Fd()), syscall.LOCK_EX|syscall.LOCK_NB); err != nil {
		file.Close()
		return nil, false
	}

	unlock := func() {
		syscall.Flock(int(file.Fd()), syscall.LOCK_UN)
		file.Close()
	}

	// Previous instance finished too recently: skip.
	var last int64
	if _, err := fmt.Fscan(file, &last); err == nil {
		if time.Since(time.Unix(0, last)) < minInterval {
			unlock()
			return nil, false
		}
	}

	if err := file.Truncate(0); err != nil {
		log.Error(err)
	}
	if _, err := file.WriteAt(
		[]byte(fmt.Sprint(time.Now().UnixNano())), 0,
	); err != nil {
		log.Error(err)
	}

	return unlock, true
}
