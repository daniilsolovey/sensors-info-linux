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

	var isConnected bool
	for _, item := range result {
		if item == "Status: Connected" {
			isConnected = true
		}
	}

	if isConnected && len(result) >= 5 {
		return result[0] + "\n" + " " + result[4], nil
	}

	return result[0], nil
}
