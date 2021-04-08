package main

import (
	"time"
)

func getLocalAndMoscowTime() (*time.Time, *time.Time, error) {
	timeMoscow, err := getTimeInLocation("Europe/Moscow")
	if err != nil {
		return nil, nil, err
	}

	timeLocal, err := getTimeInLocation("Asia/Omsk")
	if err != nil {
		return nil, nil, err
	}

	return timeMoscow, timeLocal, nil
}

func getTimeInLocation(name string) (*time.Time, error) {
	location, err := time.LoadLocation(name)
	if err != nil {
		return nil, err
	}

	result := time.Now().In(location)
	return &result, nil
}
