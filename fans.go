package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

type fanSpeed struct {
	chip  string
	index int
	label string
	rpm   int
}

func getFanSpeeds() (string, bool) {
	fans := readFanSpeeds()
	if len(fans) == 0 {
		return "", false
	}

	needChip := distinctFanChips(fans) > 1
	parts := make([]string, 0, len(fans))

	for _, fan := range fans {
		parts = append(parts, formatFanSpeed(fan, needChip))
	}

	return strings.Join(parts, "\n              "), true
}

func formatFanSpeed(fan fanSpeed, needChip bool) string {
	name := fan.label
	if name == "" {
		name = fmt.Sprintf("fan%d", fan.index)
		if needChip && fan.chip != "" {
			name = fan.chip + " " + name
		}
	}

	return fmt.Sprintf("%s %d RPM", name, fan.rpm)
}

func readFanSpeeds() []fanSpeed {
	matches, err := filepath.Glob("/sys/class/hwmon/hwmon*/fan*_input")
	if err != nil {
		return nil
	}

	fans := make([]fanSpeed, 0, len(matches))

	for _, path := range matches {
		fan, ok := readFanSpeed(path)
		if !ok {
			continue
		}

		fans = append(fans, fan)
	}

	if hasNonACPIFan(fans) {
		filtered := fans[:0]
		for _, fan := range fans {
			if fan.chip == "acpi_fan" {
				continue
			}

			filtered = append(filtered, fan)
		}

		fans = filtered
	}

	sort.Slice(fans, func(i, j int) bool {
		if fans[i].chip != fans[j].chip {
			return fans[i].chip < fans[j].chip
		}

		return fans[i].index < fans[j].index
	})

	return fans
}

func readFanSpeed(path string) (fanSpeed, bool) {
	dir := filepath.Dir(path)
	base := filepath.Base(path)

	index, err := strconv.Atoi(
		strings.TrimSuffix(strings.TrimPrefix(base, "fan"), "_input"),
	)
	if err != nil {
		return fanSpeed{}, false
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return fanSpeed{}, false
	}

	rpm, err := strconv.Atoi(strings.TrimSpace(string(data)))
	if err != nil || rpm < 0 {
		return fanSpeed{}, false
	}

	chip := strings.TrimSpace(readSysfsFile(filepath.Join(dir, "name")))
	label := strings.TrimSpace(readSysfsFile(
		filepath.Join(dir, fmt.Sprintf("fan%d_label", index)),
	))

	return fanSpeed{
		chip:  chip,
		index: index,
		label: label,
		rpm:   rpm,
	}, true
}

func readSysfsFile(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		return ""
	}

	return string(data)
}

func hasNonACPIFan(fans []fanSpeed) bool {
	for _, fan := range fans {
		if fan.chip != "acpi_fan" {
			return true
		}
	}

	return false
}

func distinctFanChips(fans []fanSpeed) int {
	chips := map[string]struct{}{}
	for _, fan := range fans {
		chips[fan.chip] = struct{}{}
	}

	return len(chips)
}
