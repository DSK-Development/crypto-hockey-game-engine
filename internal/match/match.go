package match

import (
	"sync"
	"time"

	"github.com/dskdev/crypto-hockey-game-engine/internal/physics"
	"github.com/dskdev/crypto-hockey-game-engine/internal/protocol"
)

type Phase string

const (
	PhasePending      Phase = "PENDING"
	PhaseCountdown    Phase = "COUNTDOWN"
	PhaseLive         Phase = "LIVE"
	PhaseSettled      Phase = "SETTLED"
)

type Player struct {
	UserID     string
	TelegramID int64
	Username   string
}

type Spec struct {
	ID           string
	Stake        int64
	Players      [2]Player
	JoinDeadline time.Duration
	Duration     time.Duration
	GoalCap      int
	Now          time.Time
}

type Match struct {
	mu      sync.RWMutex
	spec    Spec
	phase   Phase
	score   protocol.Score
	puck    physics.Body
	mallets [2]physics.Body
	rink    physics.Rink
	deadlineJoin time.Time
}

func New(s Spec) *Match {
	m := &Match{
		spec:         s,
		phase:        PhasePending,
		rink:         physics.DefaultRink(),
		deadlineJoin: s.Now.Add(s.JoinDeadline),
	}
	m.resetPuck()
	m.mallets[0] = physics.Body{Pos: physics.Vec2{X: 100, Y: 200}, R: 24}
	m.mallets[1] = physics.Body{Pos: physics.Vec2{X: 700, Y: 200}, R: 24}
	return m
}

func (m *Match) ID() string             { return m.spec.ID }
func (m *Match) Spec() Spec             { return m.spec }
func (m *Match) Phase() Phase           { m.mu.RLock(); defer m.mu.RUnlock(); return m.phase }
func (m *Match) Score() protocol.Score  { m.mu.RLock(); defer m.mu.RUnlock(); return m.score }
func (m *Match) JoinDeadlineAt() time.Time { return m.deadlineJoin }

func (m *Match) SetPhase(p Phase) { m.mu.Lock(); m.phase = p; m.mu.Unlock() }

// ApplyGoal increments score for the scorer and returns true if the goal cap is hit.
func (m *Match) ApplyGoal(scorer string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	if scorer == "A" {
		m.score.A++
	} else if scorer == "B" {
		m.score.B++
	}
	m.resetPuck()
	return m.score.A >= m.spec.GoalCap || m.score.B >= m.spec.GoalCap
}

func (m *Match) resetPuck() {
	m.puck = physics.Body{Pos: physics.Vec2{X: 400, Y: 200}, R: 14}
}

// SetPuck is exported for tests / debug only.
func (m *Match) SetPuck(b physics.Body) { m.mu.Lock(); m.puck = b; m.mu.Unlock() }

// SetMalletTarget updates the kinematic target for a slot's mallet.
// Server clamps to that slot's half. Returns the applied target.
func (m *Match) SetMalletTarget(slot physics.Slot, target physics.Vec2) physics.Vec2 {
	m.mu.Lock()
	defer m.mu.Unlock()
	clamped := m.rink.ClampMallet(slot, target, m.mallets[slot].R)
	m.mallets[slot].Vel = clamped.Sub(m.mallets[slot].Pos)
	m.mallets[slot].Pos = clamped
	return clamped
}

// Tick advances physics by dt and returns "A"/"B"/"" if a goal occurred.
func (m *Match) Tick(dt float64) string {
	m.mu.Lock()
	defer m.mu.Unlock()
	for i := range m.mallets {
		physics.ResolvePuckMallet(&m.puck, m.mallets[i].Pos, m.mallets[i].Vel, m.mallets[i].R)
	}
	physics.StepPuck(&m.puck, m.rink, dt)
	scorer := physics.CheckGoal(m.puck, m.rink)
	if scorer != "" {
		if scorer == "A" {
			m.score.A++
		} else {
			m.score.B++
		}
		m.puck = physics.Body{Pos: physics.Vec2{X: 400, Y: 200}, R: 14}
	}
	return scorer
}

func (m *Match) Snapshot(now time.Time) protocol.SnapshotPayload {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return protocol.SnapshotPayload{
		TServer: now.UnixMilli(),
		MalletA: protocol.V2{X: m.mallets[0].Pos.X, Y: m.mallets[0].Pos.Y},
		MalletB: protocol.V2{X: m.mallets[1].Pos.X, Y: m.mallets[1].Pos.Y},
		Puck:    protocol.PuckSnap{X: m.puck.Pos.X, Y: m.puck.Pos.Y, VX: m.puck.Vel.X, VY: m.puck.Vel.Y},
		Score:   m.score,
	}
}
