package main

import (
	"errors"
	"os/exec"
	"strings"
)

func getVPNStatus() (string, error) {
	out, err := exec.Command("/bin/sh", "-c", "sudo nordvpn status").Output()
	exec.Command("/bin/sh", "-c", "sudo find ...")
	if err != nil {
		return "", err
	}

	result := strings.Split(string(out), "\n")
	if len(result) == 0 {
		err := errors.New("unable to get nordvpn status")
		return "", err
	}

	if len(result) >= 10 {
		return "Connected" + "\n" + " " + result[4], nil
	}

	return "Disconnected", nil
}
