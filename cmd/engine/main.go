package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/dskdev/crypto-hockey-game-engine/internal/account"
	"github.com/dskdev/crypto-hockey-game-engine/internal/bot"
	"github.com/dskdev/crypto-hockey-game-engine/internal/config"
	"github.com/dskdev/crypto-hockey-game-engine/internal/httpapi"
	"github.com/dskdev/crypto-hockey-game-engine/internal/match"
	"github.com/dskdev/crypto-hockey-game-engine/internal/ws"
)

type starter struct {
	settler      match.Settler
	logger       *slog.Logger
	forfeitGrace time.Duration
	mu           sync.RWMutex
	runners      map[string]*match.Runner
}

func (s *starter) Start(m *match.Match) {
	r := match.NewRunner(m, s.settler, s.logger, s.forfeitGrace)
	s.mu.Lock()
	s.runners[m.ID()] = r
	s.mu.Unlock()
	go r.Run(context.Background())
}

type authBridge struct{ ac *account.Client }

func (a authBridge) AuthTelegram(ctx context.Context, initData string) (string, int64, string, error) {
	res, err := a.ac.AuthTelegram(ctx, initData)
	if err != nil {
		return "", 0, "", err
	}
	return res.Player.ID, res.Player.TelegramID, res.Player.Username, nil
}

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	cfg, err := config.Load()
	if err != nil {
		logger.Error("config", "err", err)
		os.Exit(1)
	}

	reg := match.NewRegistry()
	acc := account.New(cfg.AccountBaseURL, cfg.AccountServiceToken)
	bc := bot.New(cfg.BotBaseURL, cfg.BotServiceToken)
	st := match.NewSettler(acc, bc)
	start := &starter{settler: st, logger: logger, forfeitGrace: cfg.ForfeitGrace, runners: map[string]*match.Runner{}}

	router := httpapi.New(httpapi.Deps{
		Registry: reg, ServiceToken: cfg.ServiceToken,
		Starter: start,
		JoinDeadline: cfg.MatchJoinDeadline, MatchDuration: cfg.MatchDuration, GoalCap: cfg.GoalCap,
	})

	wsSrv := ws.NewServer(reg, authBridge{ac: acc}, func(m *match.Match, slot int, c *ws.Conn) {
		start.mu.RLock()
		r, ok := start.runners[m.ID()]
		start.mu.RUnlock()
		if ok {
			r.Join(c)
		}
	})

	mux := http.NewServeMux()
	mux.Handle("/", router)
	mux.HandleFunc("/ws", wsSrv.Handle)

	srv := &http.Server{Addr: cfg.HTTPAddr, Handler: mux, ReadHeaderTimeout: 5 * time.Second}
	go func() {
		logger.Info("listening", "addr", cfg.HTTPAddr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("listen", "err", err)
			os.Exit(1)
		}
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = srv.Shutdown(ctx)
}
