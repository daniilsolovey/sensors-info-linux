package main

import (
	"io/ioutil"
	"sort"
	"strconv"
	"strings"
)

func getCPUFrequency() (string, error) {
	contents, err := ioutil.ReadFile("/proc/cpuinfo")
	if err != nil {
		return "", err
	}

	var total []float64
	lines := strings.Split(string(contents), "\n")
	for _, line := range lines {
		if strings.Contains(line, "cpu MHz") {
			frequency := strings.Split(line, ": ")
			intFrequency, err := strconv.ParseFloat(frequency[1], 64)
			if err != nil {
				return "", err
			}

			total = append(total, intFrequency)
		}
	}

	var totalInt []int

	for _, item := range total {
		convertedItem := int(item)
		totalInt = append(totalInt, convertedItem)
	}

	sort.Slice(totalInt, func(i, j int) bool {
		return totalInt[i] > totalInt[j]
	})

	return strconv.Itoa(totalInt[0]), nil
}
