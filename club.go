package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"sort"
)

var (
	ErrEventOrder = errors.New("event placed in incorrect order timewise")
	ErrEventDesk  = errors.New("incorrect desk number")
)

type ClubConfig struct {
	DeskCount int
	Start     HM
	Stop      HM
	Rate      int
}

type Club struct {
	EventQueue *Queue[Event]

	ClientQueue *Queue[string]
	Clients     map[string]int

	Desks       []*Desk
	Sessions    []*Session
	VacantTotal int

	Config ClubConfig
}

func NewClub(eq *Queue[Event], cfg ClubConfig) *Club {
	desks := make([]*Desk, 0, cfg.DeskCount)
	for range cfg.DeskCount {
		desks = append(desks, &Desk{})
	}

	return &Club{
		EventQueue:  eq,
		ClientQueue: &Queue[string]{},
		Clients:     make(map[string]int),
		VacantTotal: cfg.DeskCount,
		Desks:       desks,
		Sessions:    make([]*Session, cfg.DeskCount),
		Config:      cfg,
	}
}

func ParseClub(in io.Reader) (*Club, error) {
	lineCount := 1

	s := bufio.NewScanner(in)

	s.Scan()
	line := s.Text()
	deskCount, err := ParsePositive(line)
	if err != nil {
		return nil, fmt.Errorf("line %d: %w", lineCount, err)
	}
	lineCount++

	s.Scan()
	line = s.Text()
	start, stop, err := ParseStartStop(line)
	if err != nil {
		return nil, fmt.Errorf("line %d: %w", lineCount, err)
	}
	lineCount++

	s.Scan()
	line = s.Text()
	rate, err := ParsePositive(line)
	if err != nil {
		return nil, fmt.Errorf("line %d: %w", lineCount, err)
	}
	lineCount++

	cfg := ClubConfig{
		DeskCount: deskCount,
		Start:     start,
		Stop:      stop,
		Rate:      rate,
	}

	eventQueue := &Queue[Event]{}
	var prevTime HM
	for s.Scan() {
		line = s.Text()

		e, err := ParseEvent(line)
		if err != nil {
			return nil, fmt.Errorf("line %d: %w", lineCount, err)
		}
		if prevTime.Compare(e.Time()) == 1 {
			return nil, fmt.Errorf("line %d: %w", lineCount, ErrEventOrder)
		}
		if cde, ok := e.(ClientDeskEvent); ok {
			if cde.DeskNumber > deskCount {
				return nil, fmt.Errorf("line %d: %w", lineCount, ErrEventDesk)
			}
		}

		prevTime = e.Time()
		eventQueue.Push(e)
		lineCount++
	}

	return NewClub(eventQueue, cfg), nil
}

func (c *Club) HandleIn(e Event) Event {
	fmt.Println(e.String())

	switch e.ID() {
	case InClientEnter:
		ce := e.(ClientEvent)

		client := ce.ClientName
		if c.Clients[client] != 0 {
			ee := ErrorEvent{
				BaseEvent: BaseEvent{
					OccuredAt: ce.OccuredAt,
					EID:       OutError,
				},
				Error: ErrYouShallNotPass,
			}

			return ee
		}

		if ce.OccuredAt.Compare(c.Config.Start) == -1 || ce.OccuredAt.Compare(c.Config.Stop) == 1 {
			ee := ErrorEvent{
				BaseEvent: BaseEvent{
					OccuredAt: ce.OccuredAt,
					EID:       OutError,
				},
				Error: ErrNotOpenYet,
			}

			return ee
		}

		c.Clients[client] = -1
	case InClientWait:
		ce := e.(ClientEvent)

		if c.VacantTotal != 0 {
			ee := ErrorEvent{
				BaseEvent: BaseEvent{
					OccuredAt: ce.OccuredAt,
					EID:       OutError,
				},
				Error: ErrICanWaitNoLonger,
			}

			return ee
		}

		if c.ClientQueue.Length() > c.Config.DeskCount {
			ce2 := ClientEvent{
				BaseEvent: BaseEvent{
					OccuredAt: ce.OccuredAt,
					EID:       OutClientLeave,
				},
				ClientName: ce.ClientName,
			}

			return ce2
		}

		c.ClientQueue.Push(ce.ClientName)
	case InClientLeave:
		ce := e.(ClientEvent)

		if c.Clients[ce.ClientName] == 0 {
			ee := ErrorEvent{
				BaseEvent: BaseEvent{
					OccuredAt: ce.OccuredAt,
					EID:       OutError,
				},
				Error: ErrClientUnknown,
			}

			return ee
		}

		if c.Clients[ce.ClientName] > 0 {
			deskNumber := c.Clients[ce.ClientName] - 1
			start := c.Sessions[deskNumber].Start

			delta := Sub(ce.OccuredAt, start)
			c.Desks[deskNumber].Time = Add(c.Desks[deskNumber].Time, delta)
			c.Desks[deskNumber].Balance += delta.Round() * c.Config.Rate

			c.Sessions[deskNumber] = nil

			c.VacantTotal++
		}

		defer delete(c.Clients, ce.ClientName)

		if clientName, ok := c.ClientQueue.Pop(); ok {
			cde := ClientDeskEvent{
				BaseEvent: BaseEvent{
					OccuredAt: ce.OccuredAt,
					EID:       OutClientSit,
				},
				ClientName: clientName,
				DeskNumber: c.Clients[ce.ClientName],
			}

			return cde
		}
	case InClientSit:
		cde := e.(ClientDeskEvent)

		if c.Sessions[cde.DeskNumber-1] != nil {
			ee := ErrorEvent{
				BaseEvent: BaseEvent{
					OccuredAt: cde.OccuredAt,
					EID:       OutError,
				},
				Error: ErrPlaceIsBusy,
			}

			return ee
		}

		if c.Clients[cde.ClientName] == 0 {
			ee := ErrorEvent{
				BaseEvent: BaseEvent{
					OccuredAt: cde.OccuredAt,
					EID:       OutError,
				},
				Error: ErrClientUnknown,
			}

			return ee
		}

		if c.Clients[cde.ClientName] > 0 {
			deskNumber := c.Clients[cde.ClientName] - 1
			start := c.Sessions[deskNumber].Start

			delta := Sub(cde.OccuredAt, start)
			c.Desks[deskNumber].Time = Add(c.Desks[deskNumber].Time, delta)
			c.Desks[deskNumber].Balance += delta.Round() * c.Config.Rate

			c.Sessions[deskNumber] = nil

			c.VacantTotal++
		}

		c.Clients[cde.ClientName] = cde.DeskNumber
		c.Sessions[cde.DeskNumber-1] = &Session{
			Start: cde.OccuredAt,
		}
		c.VacantTotal--
	}

	return nil
}

func (c *Club) HandleOut(e Event) {
	fmt.Println(e.String())

	switch e.ID() {
	case OutClientLeave:
		ce := e.(ClientEvent)

		if c.Clients[ce.ClientName] > 0 {
			deskNumber := c.Clients[ce.ClientName] - 1
			start := c.Sessions[deskNumber].Start

			delta := Sub(ce.OccuredAt, start)
			c.Desks[deskNumber].Time = Add(c.Desks[deskNumber].Time, delta)
			c.Desks[deskNumber].Balance += delta.Round() * c.Config.Rate

			c.Sessions[deskNumber] = nil

			c.VacantTotal++
		}

		delete(c.Clients, ce.ClientName)
	case OutClientSit:
		cde := e.(ClientDeskEvent)

		c.Clients[cde.ClientName] = cde.DeskNumber
		c.Sessions[cde.DeskNumber-1] = &Session{
			Start: cde.OccuredAt,
		}
		c.VacantTotal--
	case OutError:
		ee := e.(ErrorEvent)
		_ = ee
	}
}

func (c *Club) Close() {
	clientsRem := make([]string, 0, len(c.Clients))
	for client := range c.Clients {
		clientsRem = append(clientsRem, client)
	}
	sort.Strings(clientsRem)
	for _, client := range clientsRem {
		ce := ClientEvent{
			BaseEvent: BaseEvent{
				OccuredAt: c.Config.Stop,
				EID:       OutClientLeave,
			},
			ClientName: client,
		}

		c.HandleOut(ce)
	}
}

func (c *Club) Run() {
	fmt.Println(c.Config.Start.String())

	for e, ok := c.EventQueue.Pop(); ok; e, ok = c.EventQueue.Pop() {
		e2 := c.HandleIn(e)
		if e2 != nil {
			c.HandleOut(e2)
		}
	}

	c.Close()

	fmt.Println(c.Config.Stop.String())

	for i, desk := range c.Desks {
		fmt.Printf("%d %s\n", i+1, desk.String())
	}
}
