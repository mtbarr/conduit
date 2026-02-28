package handlers

import (
	"testing"
	"time"
)

func TestLoadCooldownDuration_Defaults(t *testing.T) {
	t.Setenv("REPORTBUG_COOLDOWN_SECONDS", "")
	if got := loadReportBugCooldownDuration(); got != 60*time.Second {
		t.Fatalf("expected default 60s, got %s", got)
	}
}

func TestLoadCooldownDuration_Env(t *testing.T) {
	t.Setenv("REPORTBUG_COOLDOWN_SECONDS", "15")
	if got := loadReportBugCooldownDuration(); got != 15*time.Second {
		t.Fatalf("expected 15s, got %s", got)
	}
}

func TestCheckAndMarkCooldown(t *testing.T) {
	cooldownMu.Lock()
	lastActionByKey = map[string]time.Time{}
	cooldownMu.Unlock()

	duration := 2 * time.Second
	allowed, remaining := checkAndMarkCooldown("user-1", duration)
	if !allowed || remaining != 0 {
		t.Fatalf("expected first call allowed, got %v %s", allowed, remaining)
	}

	allowed, remaining = checkAndMarkCooldown("user-1", duration)
	if allowed || remaining <= 0 {
		t.Fatalf("expected second call blocked with remaining, got %v %s", allowed, remaining)
	}
}
