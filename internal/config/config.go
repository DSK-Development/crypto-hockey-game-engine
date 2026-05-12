package config

import (
	"fmt"
	"time"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	HTTPAddr            string        `env:"HTTP_ADDR" envDefault:":8081"`
	ServiceToken        string        `env:"SERVICE_TOKEN,required"`
	AccountBaseURL      string        `env:"ACCOUNT_BASE_URL,required"`
	AccountServiceToken string        `env:"ACCOUNT_SERVICE_TOKEN,required"`
	BotBaseURL          string        `env:"BOT_BASE_URL,required"`
	BotServiceToken     string        `env:"BOT_SERVICE_TOKEN,required"`
	MatchJoinDeadline   time.Duration `env:"MATCH_JOIN_DEADLINE" envDefault:"60s"`
	MatchDuration       time.Duration `env:"MATCH_DURATION"       envDefault:"120s"`
	ForfeitGrace        time.Duration `env:"FORFEIT_GRACE"        envDefault:"15s"`
	GoalCap             int           `env:"GOAL_CAP"             envDefault:"5"`
	LogLevel            string        `env:"LOG_LEVEL"            envDefault:"info"`
}

func Load() (Config, error) {
	var c Config
	if err := env.Parse(&c); err != nil {
		return c, fmt.Errorf("config: %w", err)
	}
	return c, nil
}
