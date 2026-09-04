package main

import (
	"os"
	"path/filepath"
	"strings"
)

const batterySysfsPath = "/sys/class/power_supply/BAT0"

// readBatterySysfs returns the trimmed contents of a file under the
// battery sysfs directory, or an empty string if it cannot be read.
func readBatterySysfs(name string) string {
	data, err := os.ReadFile(filepath.Join(batterySysfsPath, name))
	if err != nil {
		return ""
	}

	return strings.TrimSpace(string(data))
}

// getChargeThreshold returns a charge control threshold as "NN%", or
// "n/a" when the kernel does not expose it.
func getChargeThreshold(name string) string {
	value := readBatterySysfs(name)
	if value == "" {
		return "n/a"
	}

	return value + "%"
}
