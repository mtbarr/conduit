package handlers

import (
	"math"
	"os"
	"strconv"
	"sync"
	"time"
)

const defaultCooldownSeconds = 60

var (
	cooldownMu       sync.Mutex
	lastReportByUser = map[string]time.Time{}
	cooldownDuration = loadCooldownDuration()
)

func loadCooldownDuration() time.Duration {
	value := os.Getenv("REPORTBUG_COOLDOWN_SECONDS")
	if value == "" {
		return time.Duration(defaultCooldownSeconds) * time.Second
	}

	seconds, err := strconv.Atoi(value)
	if err != nil || seconds < 1 {
		return time.Duration(defaultCooldownSeconds) * time.Second
	}

	return time.Duration(seconds) * time.Second
}

func checkAndMarkCooldown(userID string) (bool, time.Duration) {
	now := time.Now()

	cooldownMu.Lock()
	defer cooldownMu.Unlock()

	if last, ok := lastReportByUser[userID]; ok {
		elapsed := now.Sub(last)
		if elapsed < cooldownDuration {
			return false, cooldownDuration - elapsed
		}
	}

	lastReportByUser[userID] = now
	return true, 0
}

func formatCooldownRemaining(remaining time.Duration) string {
	seconds := int(math.Ceil(remaining.Seconds()))
	if seconds < 1 {
		seconds = 1
	}
	return strconv.Itoa(seconds) + "s"
}
