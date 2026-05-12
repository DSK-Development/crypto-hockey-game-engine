package match

import (
	"context"
	"log/slog"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/dskdev/crypto-hockey-game-engine/internal/protocol"
)

type fakeConn struct {
	slot   int
	userID string
	out    chan protocol.ServerEnvelope
	in     chan protocol.ClientEnvelope
	done   chan struct{}
	once   sync.Once
}

func newFakeConn(slot int, userID string) *fakeConn {
	return &fakeConn{slot: slot, userID: userID,
		out:  make(chan protocol.ServerEnvelope, 64),
		in:   make(chan protocol.ClientEnvelope, 64),
		done: make(chan struct{})}
}

func (f *fakeConn) Slot() int    { return f.slot }
func (f *fakeConn) UserID() string { return f.userID }
func (f *fakeConn) WriteEnvelope(_ context.Context, e protocol.ServerEnvelope) error {
	f.out <- e
	return nil
}
func (f *fakeConn) ReadEnvelope(_ context.Context) (protocol.ClientEnvelope, error) {
	select {
	case env := <-f.in:
		return env, nil
	default:
		return protocol.ClientEnvelope{}, context.Canceled
	}
}
func (f *fakeConn) CloseSignal() <-chan struct{} { return f.done }
func (f *fakeConn) close()                       { f.once.Do(func() { close(f.done) }) }

type fakeSettler struct {
	called int32
	winner string
	reason string
}

func (s *fakeSettler) Settle(_ context.Context, _ *Match, w, r string) {
	atomic.StoreInt32(&s.called, 1)
	s.winner = w
	s.reason = r
}

func TestRunner_NoJoinDeadline(t *testing.T) {
	m := New(Spec{
		ID: "m1",
		Players: [2]Player{{UserID: "u1"}, {UserID: "u2"}},
		JoinDeadline: 50 * time.Millisecond,
		Duration: time.Second, GoalCap: 5,
		Now: time.Now(),
	})
	s := &fakeSettler{}
	r := NewRunner(m, s, slog.Default(), time.Second)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r.Run(ctx)
	if atomic.LoadInt32(&s.called) != 1 || s.reason != "no_join" {
		t.Fatalf("expected no_join settle, got reason=%q", s.reason)
	}
}
