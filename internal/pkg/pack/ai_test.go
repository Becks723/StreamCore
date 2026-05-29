package pack

import "testing"

func TestParseBotConfigProactiveDefaults(t *testing.T) {
	raw := `{"proactive":{"enabled":true}}`

	cfg := ParseBotConfig(&raw)

	if !cfg.Proactive.Enabled {
		t.Fatal("expected proactive enabled")
	}
	if cfg.Proactive.QuietMinutes != 5 {
		t.Fatalf("expected default quiet_minutes 5, got %d", cfg.Proactive.QuietMinutes)
	}
	if cfg.Proactive.CooldownMinutes != 30 {
		t.Fatalf("expected default cooldown_minutes 30, got %d", cfg.Proactive.CooldownMinutes)
	}
	if cfg.Proactive.MaxPerDay != 5 {
		t.Fatalf("expected default max_per_day 5, got %d", cfg.Proactive.MaxPerDay)
	}
}

func TestParseBotConfigProactiveDisabledKeepsZeroValues(t *testing.T) {
	raw := `{"proactive":{"enabled":false}}`

	cfg := ParseBotConfig(&raw)

	if cfg.Proactive.Enabled {
		t.Fatal("expected proactive disabled")
	}
	if cfg.Proactive.QuietMinutes != 0 || cfg.Proactive.CooldownMinutes != 0 || cfg.Proactive.MaxPerDay != 0 {
		t.Fatalf("expected disabled proactive config to keep zero values, got %+v", cfg.Proactive)
	}
}
