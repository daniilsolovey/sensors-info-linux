package main

import (
	"os/exec"
	"strings"
)

func getCommonVPNStatus() (string, error) {
	out, err := exec.Command(
		"nmcli",
		"-t",
		"-f", "TYPE",
		"con", "show",
		"--active",
	).Output()
	if err != nil {
		return "", err
	}

	for _, connType := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		switch connType {
		case "vpn", "wireguard", "tun":
			return "Connected", nil
		}
	}

	return "Disconnected", nil
}
