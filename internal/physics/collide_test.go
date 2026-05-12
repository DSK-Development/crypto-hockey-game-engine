package physics

import "testing"

func TestResolvePuckMallet_SeparatesAndReflects(t *testing.T) {
	p := &Body{Pos: Vec2{100, 200}, Vel: Vec2{-50, 0}, R: 14}
	// mallet overlapping puck, stationary
	hit := ResolvePuckMallet(p, Vec2{80, 200}, Vec2{}, 24)
	if !hit {
		t.Fatal("expected collision")
	}
	if p.Vel.X <= 0 {
		t.Fatalf("puck should reflect to +X, got %v", p.Vel.X)
	}
	// puck must no longer overlap mallet
	delta := p.Pos.Sub(Vec2{80, 200})
	if delta.Len() < 14+24-1e-6 {
		t.Fatalf("puck still overlapping: dist=%v", delta.Len())
	}
}

func TestResolvePuckMallet_NoCollisionWhenFar(t *testing.T) {
	p := &Body{Pos: Vec2{100, 200}, Vel: Vec2{0, 0}, R: 14}
	if ResolvePuckMallet(p, Vec2{500, 200}, Vec2{}, 24) {
		t.Fatal("should not collide")
	}
}

func TestResolvePuckMallet_AlreadySeparating(t *testing.T) {
	p := &Body{Pos: Vec2{100, 200}, Vel: Vec2{50, 0}, R: 14} // moving away
	startV := p.Vel
	ResolvePuckMallet(p, Vec2{80, 200}, Vec2{}, 24)
	// velocity should be unchanged on the separating axis (still positive X)
	if p.Vel.X < startV.X-1e-6 {
		t.Fatalf("vel decreased unexpectedly: %v -> %v", startV.X, p.Vel.X)
	}
}
