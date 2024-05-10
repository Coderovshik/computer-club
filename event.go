package main

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrYouShallNotPass  = errors.New("YouShallNotPass")
	ErrNotOpenYet       = errors.New("NotOpenYet")
	ErrPlaceIsBusy      = errors.New("PlaceIsBusy")
	ErrClientUnknown    = errors.New("ClientUnknown")
	ErrICanWaitNoLonger = errors.New("ICanWaitNoLonger!")
)

var (
	ErrEventValueCount = errors.New("incorrect number of values representing an event")
	ErrEventID         = errors.New("incorrect event id")
)

type EventID int

const (
	InClientEnter EventID = iota + 1
	InClientSit
	InClientWait
	InClientLeave
)

const (
	OutClientLeave EventID = iota + 11
	OutClientSit
	OutError
)

type Event interface {
	Time() HM
	ID() EventID
	String() string
}

type BaseEvent struct {
	OccuredAt HM
	EID       EventID
}

func (e BaseEvent) Time() HM {
	return e.OccuredAt
}

func (e BaseEvent) ID() EventID {
	return e.EID
}

func (e BaseEvent) String() string {
	return fmt.Sprintf("%s %d", e.OccuredAt, e.EID)
}

type ClientEvent struct {
	BaseEvent
	ClientName string
}

func (e ClientEvent) String() string {
	return fmt.Sprintf("%s %s", e.BaseEvent.String(), e.ClientName)
}

type ClientDeskEvent struct {
	BaseEvent
	ClientName string
	DeskNumber int
}

func (e ClientDeskEvent) String() string {
	return fmt.Sprintf("%s %s %d", e.BaseEvent.String(), e.ClientName, e.DeskNumber)
}

type ErrorEvent struct {
	BaseEvent
	Error error
}

func (e ErrorEvent) String() string {
	return fmt.Sprintf("%s %s", e.BaseEvent.String(), e.Error.Error())
}

func ParseEvent(s string) (Event, error) {
	fields := strings.Split(s, " ")
	if len(fields) < 3 {
		return nil, ErrEventValueCount
	}

	t, err := ParseHM(fields[0])
	if err != nil {
		return nil, err
	}

	id, err := ParsePositive(fields[1])
	if err != nil {
		return nil, err
	}

	clientName, err := ParseClientName(fields[2])
	if err != nil {
		return nil, err
	}

	switch eid := EventID(id); eid {
	case InClientEnter, InClientWait, InClientLeave:
		if len(fields) != 3 {
			return nil, ErrEventValueCount
		}

		return ClientEvent{
			BaseEvent: BaseEvent{
				OccuredAt: t,
				EID:       eid,
			},
			ClientName: clientName,
		}, nil
	case InClientSit:
		if len(fields) != 4 {
			return nil, ErrEventValueCount
		}

		deskNumber, err := ParsePositive(fields[3])
		if err != nil {
			return nil, err
		}

		return ClientDeskEvent{
			BaseEvent: BaseEvent{
				OccuredAt: t,
				EID:       eid,
			},
			ClientName: clientName,
			DeskNumber: deskNumber,
		}, nil
	default:
		return nil, ErrEventID
	}
}
