package match

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/dskdev/crypto-hockey-game-engine/internal/physics"
	"github.com/dskdev/crypto-hockey-game-engine/internal/protocol"
)

type ConnHandle interface {
	Slot() int
	UserID() string
	WriteEnvelope(ctx context.Context, e protocol.ServerEnvelope) error
	ReadEnvelope(ctx context.Context) (protocol.ClientEnvelope, error)
	CloseSignal() <-chan struct{}
}

type Settler interface {
	Settle(ctx context.Context, m *Match, winner string, reason string)
}

type Runner struct {
	m       *Match
	joins   chan ConnHandle
	conns   [2]ConnHandle
	mu      sync.Mutex
	settler Settler
	logger  *slog.Logger
	forfeitGrace time.Duration
}

func NewRunner(m *Match, settler Settler, logger *slog.Logger, forfeitGrace time.Duration) *Runner {
	return &Runner{
		m:            m,
		joins:        make(chan ConnHandle, 2),
		settler:      settler,
		logger:       logger,
		forfeitGrace: forfeitGrace,
	}
}

func (r *Runner) Join(c ConnHandle) { r.joins <- c }

// Run blocks until the match is settled. Caller spawns it in a goroutine.
func (r *Runner) Run(ctx context.Context) {
	const tickHz = 60
	const snapEvery = 3 // → 20 Hz
	tick := time.NewTicker(time.Second / tickHz)
	defer tick.Stop()

	joinDeadline := time.NewTimer(time.Until(r.m.JoinDeadlineAt()))
	defer joinDeadline.Stop()
	matchTimer := time.NewTimer(r.m.Spec().Duration)
	matchTimer.Stop() // start when LIVE
	defer matchTimer.Stop()

	var disconnectAt [2]*time.Time
	frame := 0
	for {
		select {
		case <-ctx.Done():
			return
		case c := <-r.joins:
			r.attach(c)
			if r.bothJoined() {
				if !joinDeadline.Stop() {
					select { case <-joinDeadline.C: default: }
				}
				r.m.SetPhase(PhaseCountdown)
				r.broadcastState(ctx, PhaseCountdown, intPtr(3000), nil)
				time.AfterFunc(3*time.Second, func() {
					if !r.m.trySetLive() {
						return // match already settled before countdown ended
					}
					matchTimer.Reset(r.m.Spec().Duration)
					r.broadcastState(ctx, PhaseLive, nil, intPtr(int(r.m.Spec().Duration/time.Millisecond)))
				})
			}
		case <-joinDeadline.C:
			if r.m.Phase() == PhasePending {
				r.settler.Settle(ctx, r.m, "", "no_join")
				return
			}
		case <-matchTimer.C:
			r.endByScore(ctx, "timeout")
			return
		case <-tick.C:
			if r.m.Phase() != PhaseLive {
				continue
			}
			r.drainInputs(ctx)
			if scorer := r.m.Tick(1.0 / float64(tickHz)); scorer != "" {
				r.broadcastGoal(ctx, scorer, r.m.Score())
				if r.m.Score().A >= r.m.Spec().GoalCap || r.m.Score().B >= r.m.Spec().GoalCap {
					r.endByScore(ctx, "score")
					return
				}
			}
			r.checkDisconnects(ctx, &disconnectAt)
			frame++
			if frame%snapEvery == 0 {
				r.broadcast(ctx, protocol.ServerEnvelope{
					Type:     protocol.TypeSnapshot,
					Snapshot: ptrSnap(r.m.Snapshot(time.Now())),
				})
			}
		}
	}
}

func (r *Runner) attach(c ConnHandle) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.conns[c.Slot()] = c
}

func (r *Runner) bothJoined() bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.conns[0] != nil && r.conns[1] != nil
}

func (r *Runner) drainInputs(ctx context.Context) {
	r.mu.Lock()
	conns := r.conns
	r.mu.Unlock()
	for _, c := range conns {
		if c == nil {
			continue
		}
		// non-blocking: rely on a short read deadline implemented by ConnHandle.
		env, err := c.ReadEnvelope(ctx)
		if err != nil || env.Type != protocol.TypeInput || env.Input == nil {
			continue
		}
		r.m.SetMalletTarget(physics.Slot(c.Slot()),
			physics.Vec2{X: env.Input.MalletTarget.X, Y: env.Input.MalletTarget.Y})
	}
}

func (r *Runner) checkDisconnects(ctx context.Context, at *[2]*time.Time) {
	r.mu.Lock()
	conns := r.conns
	r.mu.Unlock()
	now := time.Now()
	for i, c := range conns {
		if c == nil {
			continue
		}
		select {
		case <-c.CloseSignal():
			if at[i] == nil {
				t := now
				at[i] = &t
				r.logger.Info("conn dropped", "slot", i)
			}
			if now.Sub(*at[i]) >= r.forfeitGrace {
				winnerSlot := 1 - i
				winner := r.m.Spec().Players[winnerSlot].UserID
				slotName := [2]string{"A", "B"}[winnerSlot]
				score := r.m.Score()
				r.broadcast(ctx, protocol.ServerEnvelope{
					Type: protocol.TypeMatchEnd,
					MatchEnd: &protocol.MatchEndPayload{
						WinnerUserID: &winner,
						WinnerSlot:   &slotName,
						Reason:       "forfeit",
						FinalScore:   score,
					},
				})
				r.settler.Settle(ctx, r.m, winner, "forfeit")
				return
			}
		default:
			if at[i] != nil {
				at[i] = nil
			}
		}
	}
}

func (r *Runner) endByScore(ctx context.Context, reason string) {
	s := r.m.Score()
	var winner, winnerSlot string
	if s.A > s.B {
		winner = r.m.Spec().Players[0].UserID
		winnerSlot = "A"
	} else if s.B > s.A {
		winner = r.m.Spec().Players[1].UserID
		winnerSlot = "B"
	}
	r.broadcast(ctx, protocol.ServerEnvelope{
		Type: protocol.TypeMatchEnd,
		MatchEnd: &protocol.MatchEndPayload{
			WinnerUserID: ptrStr(winner),
			WinnerSlot:   ptrStr(winnerSlot),
			Reason:       reason,
			FinalScore:   s,
		},
	})
	r.settler.Settle(ctx, r.m, winner, reason)
}

func (r *Runner) broadcast(ctx context.Context, e protocol.ServerEnvelope) {
	r.mu.Lock()
	conns := r.conns
	r.mu.Unlock()
	for _, c := range conns {
		if c == nil {
			continue
		}
		_ = c.WriteEnvelope(ctx, e)
	}
}

func (r *Runner) broadcastState(ctx context.Context, p Phase, countdown, dur *int) {
	r.broadcast(ctx, protocol.ServerEnvelope{
		Type: protocol.TypeMatchState,
		MatchState: &protocol.MatchStatePayload{Phase: string(p), CountdownMs: countdown, DurationLeftMs: dur},
	})
}

func (r *Runner) broadcastGoal(ctx context.Context, scorer string, s protocol.Score) {
	r.broadcast(ctx, protocol.ServerEnvelope{
		Type: protocol.TypeGoal,
		Goal: &protocol.GoalPayload{Scorer: scorer, Score: s},
	})
}

func intPtr(i int) *int                                             { return &i }
func ptrSnap(s protocol.SnapshotPayload) *protocol.SnapshotPayload { return &s }
func ptrStr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
