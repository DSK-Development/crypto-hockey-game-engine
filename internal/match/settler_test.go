package match

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/dskdev/crypto-hockey-game-engine/internal/account"
	"github.com/dskdev/crypto-hockey-game-engine/internal/bot"
)

type fakeAcc struct{ calls int32; last account.MatchResult }
func (f *fakeAcc) PostMatchResult(_ context.Context, r account.MatchResult) error {
	atomic.AddInt32(&f.calls, 1); f.last = r; return nil
}

type fakeBot struct{ calls int32; last bot.MatchResult }
func (f *fakeBot) NotifyResult(_ context.Context, _ string, r bot.MatchResult) error {
	atomic.AddInt32(&f.calls, 1); f.last = r; return nil
}

func TestSettler_ConcurrentCallsSettleExactlyOnce(t *testing.T) {
	m := New(Spec{
		ID: "race", Stake: 100,
		Players:      [2]Player{{UserID: "u1"}, {UserID: "u2"}},
		JoinDeadline: time.Second, Duration: time.Second, GoalCap: 5,
		Now: time.Now(),
	})
	a := &fakeAcc{}
	b := &fakeBot{}
	s := NewSettler(a, b)

	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			s.Settle(context.Background(), m, "u1", "score")
		}()
	}
	wg.Wait()

	if atomic.LoadInt32(&a.calls) != 1 {
		t.Fatalf("PostMatchResult called %d times, want 1", atomic.LoadInt32(&a.calls))
	}
	if atomic.LoadInt32(&b.calls) != 1 {
		t.Fatalf("NotifyResult called %d times, want 1", atomic.LoadInt32(&b.calls))
	}
}

func TestSettler_PostsBothSides(t *testing.T) {
	m := New(Spec{
		ID: "m1", Stake: 50,
		Players: [2]Player{{UserID: "u1"}, {UserID: "u2"}},
		JoinDeadline: time.Second, Duration: time.Second, GoalCap: 5,
		Now: time.Now(),
	})
	a := &fakeAcc{}
	b := &fakeBot{}
	s := NewSettler(a, b)
	s.Settle(context.Background(), m, "u1", "score")
	if a.calls != 1 || b.calls != 1 {
		t.Fatalf("calls: acc=%d bot=%d", a.calls, b.calls)
	}
	if a.last.Participants[0].Placement != 1 || a.last.Participants[1].Placement != 2 {
		t.Fatalf("placements: %+v", a.last.Participants)
	}
}
