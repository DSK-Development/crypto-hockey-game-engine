package physics

import "testing"

func TestCheckGoal_PuckPastLeft_ReturnsB(t *testing.T) {
	r := DefaultRink()
	p := Body{Pos: Vec2{r.Left - 20, 200}, R: 14}
	if got := CheckGoal(p, r); got != "B" {
		t.Fatalf("expected B, got %q", got)
	}
}

func TestCheckGoal_PuckPastRight_ReturnsA(t *testing.T) {
	r := DefaultRink()
	p := Body{Pos: Vec2{r.Right + 20, 200}, R: 14}
	if got := CheckGoal(p, r); got != "A" {
		t.Fatalf("expected A, got %q", got)
	}
}

func TestCheckGoal_OutsideMouth_NoGoal(t *testing.T) {
	r := DefaultRink()
	p := Body{Pos: Vec2{r.Left - 20, 50}, R: 14} // above goal mouth
	if got := CheckGoal(p, r); got != "" {
		t.Fatalf("expected no goal, got %q", got)
	}
}
