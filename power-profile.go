package main

import (
	"fmt"
	"os/exec"
	"strings"
)

func getPowerProfile() (string, error) {
	out, err := exec.Command("powerprofilesctl", "get").Output()
	if err != nil {
		return "", err
	}

	profile := strings.TrimSpace(string(out))
	if profile == "" {
		return "", fmt.Errorf("power profile is empty")
	}

	return profile, nil
}
