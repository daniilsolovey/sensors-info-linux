package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func getCPUFrequency() (string, error) {
	frequencies, err := readSysfsCPUFrequencies()
	if err != nil || len(frequencies) == 0 {
		frequencies, err = readProcCPUFrequencies()
		if err != nil {
			return "", err
		}
	}

	efficientCPUs, err := readCPUListFile("/sys/devices/cpu_atom/cpus")
	if err != nil {
		efficientCPUs = nil
	}

	type group struct {
		cores map[int]struct{}
		total float64
		count int
	}

	efficient := group{cores: map[int]struct{}{}}
	other := group{cores: map[int]struct{}{}}

	for cpu, frequency := range frequencies {
		target := &other
		if _, ok := efficientCPUs[cpu]; ok {
			target = &efficient
		}

		target.total += frequency
		target.count++
		target.cores[cpuCoreID(cpu)] = struct{}{}
	}

	var parts []string

	if efficient.count > 0 {
		parts = append(parts, formatCPUGroup(
			"E-cores",
			len(efficient.cores),
			efficient.total/float64(efficient.count),
		))
	}

	if other.count > 0 {
		name := "P-cores"
		if efficient.count == 0 {
			name = "cores"
		}

		parts = append(parts, formatCPUGroup(
			name,
			len(other.cores),
			other.total/float64(other.count),
		))
	}

	if len(parts) == 0 {
		return "", fmt.Errorf("cpu frequency not found")
	}

	return strings.Join(parts, "\n              "), nil
}

func formatCPUGroup(name string, cores int, avgMHz float64) string {
	return fmt.Sprintf("%s %d %.0f(avgMHz)", name, cores, avgMHz)
}

func readSysfsCPUFrequencies() (map[int]float64, error) {
	matches, err := filepath.Glob(
		"/sys/devices/system/cpu/cpufreq/policy*/scaling_cur_freq",
	)
	if err != nil {
		return nil, err
	}

	frequencies := make(map[int]float64, len(matches))

	for _, path := range matches {
		cpu, ok := policyCPU(path)
		if !ok {
			continue
		}

		data, err := os.ReadFile(path)
		if err != nil {
			continue
		}

		kHz, err := strconv.ParseFloat(strings.TrimSpace(string(data)), 64)
		if err != nil || kHz <= 0 {
			continue
		}

		frequencies[cpu] = kHz / 1000
	}

	if len(frequencies) == 0 {
		return nil, fmt.Errorf("cpu frequency not found")
	}

	return frequencies, nil
}

func policyCPU(path string) (int, bool) {
	dir := filepath.Base(filepath.Dir(path))
	id, err := strconv.Atoi(strings.TrimPrefix(dir, "policy"))
	if err != nil {
		return 0, false
	}

	return id, true
}

func readProcCPUFrequencies() (map[int]float64, error) {
	file, err := os.Open("/proc/cpuinfo")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	frequencies := map[int]float64{}
	scanner := bufio.NewScanner(file)

	var (
		cpu    int
		hasCPU bool
	)

	for scanner.Scan() {
		line := scanner.Text()

		switch {
		case strings.HasPrefix(line, "processor"):
			_, value, ok := strings.Cut(line, ":")
			if !ok {
				continue
			}

			id, err := strconv.Atoi(strings.TrimSpace(value))
			if err != nil {
				return nil, err
			}

			cpu = id
			hasCPU = true

		case hasCPU && strings.HasPrefix(line, "cpu MHz"):
			_, value, ok := strings.Cut(line, ":")
			if !ok {
				continue
			}

			frequency, err := strconv.ParseFloat(
				strings.TrimSpace(value),
				64,
			)
			if err != nil {
				return nil, err
			}

			if frequency > 0 {
				frequencies[cpu] = frequency
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	if len(frequencies) == 0 {
		return nil, fmt.Errorf("cpu frequency not found")
	}

	return frequencies, nil
}

func cpuCoreID(cpu int) int {
	data, err := os.ReadFile(fmt.Sprintf(
		"/sys/devices/system/cpu/cpu%d/topology/core_id",
		cpu,
	))
	if err != nil {
		return cpu
	}

	id, err := strconv.Atoi(strings.TrimSpace(string(data)))
	if err != nil {
		return cpu
	}

	return id
}

func readCPUListFile(path string) (map[int]struct{}, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return parseCPUList(string(data))
}

func parseCPUList(s string) (map[int]struct{}, error) {
	s = strings.TrimSpace(s)
	cpus := map[int]struct{}{}

	if s == "" {
		return cpus, nil
	}

	for _, part := range strings.Split(s, ",") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		start, end, isRange := strings.Cut(part, "-")
		if !isRange {
			id, err := strconv.Atoi(part)
			if err != nil {
				return nil, err
			}

			cpus[id] = struct{}{}
			continue
		}

		from, err := strconv.Atoi(strings.TrimSpace(start))
		if err != nil {
			return nil, err
		}

		to, err := strconv.Atoi(strings.TrimSpace(end))
		if err != nil {
			return nil, err
		}

		if to < from {
			return nil, fmt.Errorf("invalid cpu list %q", s)
		}

		for id := from; id <= to; id++ {
			cpus[id] = struct{}{}
		}
	}

	return cpus, nil
}
