package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var (
	ErrHMValue  error = errors.New("hour minute pair is incorrect")
	ErrHMFormat error = errors.New("time does not match the XX:XX format")
)

type HM struct {
	H int // часы
	M int // минуты
}

func ParseHM(value string) (HM, error) {
	if !IsValidHMString(value) {
		return HM{}, ErrHMFormat
	}

	fields := strings.Split(value, ":")
	h, _ := strconv.Atoi(fields[0])
	m, _ := strconv.Atoi(fields[1])

	return NewHM(h, m)
}

func NewHM(h int, m int) (HM, error) {
	var hm HM
	if !IsValidHour(h) || !IsValidMinute(m) {
		return hm, ErrHMValue
	}

	hm.H, hm.M = h, m

	return hm, nil
}

func (hm HM) String() string {
	return fmt.Sprintf("%02d:%02d", hm.H, hm.M)
}

func (hm HM) Round() int {
	res := hm.H
	if hm.M != 0 {
		res++
	}

	return res
}

func (hm HM) Compare(other HM) int {
	if hm.H > other.H {
		return 1
	} else if hm.H < other.H {
		return -1
	} else {
		if hm.M > other.M {
			return 1
		} else if hm.M < other.M {
			return -1
		} else {
			return 0
		}
	}
}

func Sub(a HM, b HM) HM {
	h := a.H - b.H
	m := a.M - b.M
	if m < 0 {
		h--
		m = 60 + m
	}

	return HM{
		H: h,
		M: m,
	}
}

func Add(a HM, b HM) HM {
	h := a.H + b.H
	m := a.M + b.M
	if m >= 60 {
		m = m % 60
		h++
	}

	return HM{
		H: h,
		M: m,
	}
}
