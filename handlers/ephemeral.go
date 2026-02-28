package handlers

import (
	"os"
	"strconv"
	"time"
)

const defaultEphemeralDeleteSeconds = 3

var ephemeralDeleteDelay = loadEphemeralDeleteDelay()

func loadEphemeralDeleteDelay() time.Duration {
	value := os.Getenv("EPHEMERAL_DELETE_SECONDS")
	if value == "" {
		return time.Duration(defaultEphemeralDeleteSeconds) * time.Second
	}

	seconds, err := strconv.Atoi(value)
	if err != nil || seconds < 0 {
		return time.Duration(defaultEphemeralDeleteSeconds) * time.Second
	}

	return time.Duration(seconds) * time.Second
}
