package main

import (
	"os/exec"
	"strings"

	"github.com/pkg/errors"
)

func getVPNStatus() (string, error) {
	out, err := exec.Command("nordvpn", "status").Output()
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
