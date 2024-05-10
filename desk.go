package main

import "fmt"

type Desk struct {
	Time    HM
	Balance int
}

func (d Desk) String() string {
	return fmt.Sprintf("%d %s", d.Balance, d.Time)
}

type Session struct {
	Start HM
}
