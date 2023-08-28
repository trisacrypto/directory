package bff

import "time"

// currentTime shadows the time.Now function so that it can be mocked out for testing.
var currentTime = time.Now

// MockTime is used by tests to mock the current time.
func MockTime(newTime time.Time) {
	currentTime = func() time.Time { return newTime }
}

// ResetTime is used by tests to reset the time method to use the original time.Now.
func ResetTime() {
	currentTime = time.Now
}
