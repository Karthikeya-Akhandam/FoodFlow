package lib

import "time"

// Clock interface for time abstraction (useful for testing)
type Clock interface {
	Now() time.Time
}

// RealClock implements Clock using real time
type RealClock struct{}

func (RealClock) Now() time.Time {
	return time.Now()
}

// MockClock implements Clock for testing
type MockClock struct {
	Time time.Time
}

func (m MockClock) Now() time.Time {
	return m.Time
}

// Global clock instance
var GlobalClock Clock = RealClock{}

// Now returns the current time using the global clock
func Now() time.Time {
	return GlobalClock.Now()
}

// SetClock sets the global clock (mainly for testing)
func SetClock(clock Clock) {
	GlobalClock = clock
}
