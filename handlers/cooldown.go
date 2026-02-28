package handlers

import (
	"math"
	"os"
	"strconv"
	"sync"
	"time"
)

const (
	defaultReportBugCooldownSeconds      = 60
	defaultRequestFeatureCooldownSeconds = 60
)

var (
	cooldownMu                     sync.Mutex
	lastActionByKey                = map[string]time.Time{}
	reportbugCooldownDuration      = loadReportBugCooldownDuration()
	requestfeatureCooldownDuration = loadRequestFeatureCooldownDuration()
)

func loadReportBugCooldownDuration() time.Duration {
	value := os.Getenv("REPORTBUG_COOLDOWN_SECONDS")
	if value == "" {
		return time.Duration(defaultReportBugCooldownSeconds) * time.Second
	}

	seconds, err := strconv.Atoi(value)
	if err != nil || seconds < 1 {
		return time.Duration(defaultReportBugCooldownSeconds) * time.Second
	}

	return time.Duration(seconds) * time.Second
}

func loadRequestFeatureCooldownDuration() time.Duration {
	value := os.Getenv("REQUESTFEATURE_COOLDOWN_SECONDS")
	if value == "" {
		return time.Duration(defaultRequestFeatureCooldownSeconds) * time.Second
	}

	seconds, err := strconv.Atoi(value)
	if err != nil || seconds < 1 {
		return time.Duration(defaultRequestFeatureCooldownSeconds) * time.Second
	}

	return time.Duration(seconds) * time.Second
}

func checkAndMarkCooldown(key string, duration time.Duration) (bool, time.Duration) {
	now := time.Now()

	cooldownMu.Lock()
	defer cooldownMu.Unlock()

	if last, ok := lastActionByKey[key]; ok {
		elapsed := now.Sub(last)
		if elapsed < duration {
			return false, duration - elapsed
		}
	}

	lastActionByKey[key] = now
	return true, 0
}

func formatCooldownRemaining(remaining time.Duration) string {
	seconds := int(math.Ceil(remaining.Seconds()))
	if seconds < 1 {
		seconds = 1
	}
	return strconv.Itoa(seconds) + "s"
}
