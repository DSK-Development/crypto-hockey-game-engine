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
