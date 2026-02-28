package handlers

import (
	"testing"
	"time"
)

func TestLoadCooldownDuration_Defaults(t *testing.T) {
	reset := setCooldownStateForTest(1 * time.Second)
	defer reset()

	t.Setenv("REPORTBUG_COOLDOWN_SECONDS", "")
	if got := loadCooldownDuration(); got != 60*time.Second {
		t.Fatalf("expected default 60s, got %s", got)
	}
}

func TestLoadCooldownDuration_Env(t *testing.T) {
	t.Setenv("REPORTBUG_COOLDOWN_SECONDS", "15")
	if got := loadCooldownDuration(); got != 15*time.Second {
		t.Fatalf("expected 15s, got %s", got)
	}
}

func TestCheckAndMarkCooldown(t *testing.T) {
	reset := setCooldownStateForTest(2 * time.Second)
	defer reset()

	allowed, remaining := checkAndMarkCooldown("user-1")
	if !allowed || remaining != 0 {
		t.Fatalf("expected first call allowed, got %v %s", allowed, remaining)
	}

	allowed, remaining = checkAndMarkCooldown("user-1")
	if allowed || remaining <= 0 {
		t.Fatalf("expected second call blocked with remaining, got %v %s", allowed, remaining)
	}
}
