package main

import (
	"errors"
	"strconv"
	"strings"
)

var (
	ErrNotPositiveInt = errors.New("not a positive integer")
	ErrNotClientName  = errors.New("not a client name")
	ErrValueCount     = errors.New("incorrect number of values")
	ErrWorkTime       = errors.New("incorrect club working period")
)

func ParsePositive(s string) (int, error) {
	num, err := strconv.Atoi(s)
	if err != nil || num <= 0 {
		return 0, ErrNotPositiveInt
	}

	return num, nil
}

func ParseClientName(s string) (string, error) {
	if !IsValidClientName(s) {
		return "", ErrNotClientName
	}

	return s, nil
}

func ParseStartStop(s string) (HM, HM, error) {
	var start, stop HM

	fields := strings.Split(s, " ")
	if len(fields) != 2 {
		return start, stop, ErrValueCount
	}

	start, err := ParseHM(fields[0])
	if err != nil {
		return start, stop, err
	}

	stop, err = ParseHM(fields[1])
	if err != nil {
		return start, stop, err
	}

	if start.Compare(stop) != -1 {
		return HM{}, HM{}, ErrWorkTime
	}

	return start, stop, nil
}
