package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func getCPUFrequency() (string, error) {
	file, err := os.Open("/proc/cpuinfo")
	if err != nil {
		return "", err
	}
	defer file.Close()

	var (
		totalFrequency float64
		maxFrequency   float64
		count          int
	)

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()

		if !strings.HasPrefix(line, "cpu MHz") {
			continue
		}

		_, value, ok := strings.Cut(line, ":")
		if !ok {
			continue
		}

		frequency, err := strconv.ParseFloat(
			strings.TrimSpace(value),
			64,
		)
		if err != nil {
			return "", err
		}

		totalFrequency += frequency
		count++

		if frequency > maxFrequency {
			maxFrequency = frequency
		}
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	if count == 0 {
		return "", fmt.Errorf("cpu frequency not found")
	}

	avgFrequency := totalFrequency / float64(count)

	return fmt.Sprintf(
		"AVG %.0f / MAX %.0f",
		avgFrequency,
		maxFrequency,
	), nil
}
