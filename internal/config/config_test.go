package config

import (
	"testing"
)

func TestLoad_RequiresServiceToken(t *testing.T) {
	t.Setenv("ACCOUNT_BASE_URL", "http://acc")
	t.Setenv("ACCOUNT_SERVICE_TOKEN", "s")
	t.Setenv("BOT_BASE_URL", "http://bot")
	t.Setenv("BOT_SERVICE_TOKEN", "s")
	if _, err := Load(); err == nil {
		t.Fatal("expected error for missing SERVICE_TOKEN")
	}
}

func TestLoad_DefaultsApply(t *testing.T) {
	t.Setenv("SERVICE_TOKEN", "x")
	t.Setenv("ACCOUNT_BASE_URL", "http://acc")
	t.Setenv("ACCOUNT_SERVICE_TOKEN", "s")
	t.Setenv("BOT_BASE_URL", "http://bot")
	t.Setenv("BOT_SERVICE_TOKEN", "s")
	c, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if c.GoalCap != 5 {
		t.Fatalf("goal cap default: %d", c.GoalCap)
	}
	if c.HTTPAddr != ":8081" {
		t.Fatalf("addr default: %s", c.HTTPAddr)
	}
}
