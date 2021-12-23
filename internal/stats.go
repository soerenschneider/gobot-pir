package internal

import (
	"log"
	"sync"
	"time"
)

const maxEvents = 2048

type SensorStats struct {
	events []time.Time
	m      sync.Mutex
}

func NewSensorStats() *SensorStats {
	return &SensorStats{
		events: make([]time.Time, 0),
	}
}

func (s *SensorStats) NewEvent() {
	s.m.Lock()
	defer s.m.Unlock()

	if len(s.events) < maxEvents {
		s.events = append(s.events, time.Now())
	} else {
		log.Println("Not adding further events to stats, slice is full")
	}
}

func (s *SensorStats) GetStatsSliceSize() int {
	s.m.Lock()
	defer s.m.Unlock()
	return len(s.events)
}

func (s *SensorStats) GetEventCountNewerThan(window time.Duration) int {
	s.m.Lock()
	defer s.m.Unlock()
	idx := s.getIndexOfEventsNewerThan(time.Now().Add(-window))
	return len(s.events) - idx
}

func (s *SensorStats) getIndexOfEventsNewerThan(timestamp time.Time) int {
	for index, event := range s.events {
		if event.After(timestamp) {
			return index
		}
	}

	return len(s.events)
}

func (s *SensorStats) PurgeEventsBefore(timestamp time.Time) {
	s.m.Lock()
	defer s.m.Unlock()
	marker := s.getIndexOfEventsNewerThan(timestamp)
	s.events = s.events[marker:]
}
