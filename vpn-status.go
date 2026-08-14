package main

import (
	"os/exec"
	"strings"
)

func getCommonVPNStatus() (string, error) {
	out, err := exec.Command("/bin/sh", "-c", "nmcli con show --active").Output()
	exec.Command("/bin/sh", "-c", "sudo find ...")
	if err != nil {
		return "", err
	}

	if strings.Contains(string(out), "tun0") {
		return "Connected", nil
	} else {
		return "Disconnected", nil
	}

}
