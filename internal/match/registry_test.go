package match

import (
	"testing"
	"time"
)

func TestRegistry_GetAfterCreate(t *testing.T) {
	r := NewRegistry()
	m := New(testSpec(t, 5))
	if err := r.Add(m); err != nil {
		t.Fatal(err)
	}
	got, ok := r.Get("m1")
	if !ok || got.ID() != "m1" {
		t.Fatal("not found")
	}
}

func TestRegistry_DuplicateRejected(t *testing.T) {
	r := NewRegistry()
	m := New(testSpec(t, 5))
	_ = r.Add(m)
	if err := r.Add(New(testSpec(t, 5))); err == nil {
		t.Fatal("expected duplicate error")
	}
}

func TestRegistry_RemoveAndMissing(t *testing.T) {
	r := NewRegistry()
	_ = r.Add(New(testSpec(t, 5)))
	r.Remove("m1")
	if _, ok := r.Get("m1"); ok {
		t.Fatal("should be gone")
	}
}

func TestRegistry_Concurrent(t *testing.T) {
	r := NewRegistry()
	done := make(chan struct{}, 2)
	for i := 0; i < 2; i++ {
		go func() {
			for j := 0; j < 100; j++ {
				_, _ = r.Get("missing")
			}
			done <- struct{}{}
		}()
	}
	<-done
	<-done
	_ = time.Now()
}
