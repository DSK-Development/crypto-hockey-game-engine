package physics

import (
	"math"
	"testing"
)

func eq(a, b float64) bool { return math.Abs(a-b) < 1e-9 }

func TestVec_AddSubScale(t *testing.T) {
	a := Vec2{1, 2}
	b := Vec2{3, 4}
	if got := a.Add(b); got != (Vec2{4, 6}) {
		t.Fatalf("add: %+v", got)
	}
	if got := b.Sub(a); got != (Vec2{2, 2}) {
		t.Fatalf("sub: %+v", got)
	}
	if got := a.Scale(2); got != (Vec2{2, 4}) {
		t.Fatalf("scale: %+v", got)
	}
}

func TestVec_LenAndNorm(t *testing.T) {
	v := Vec2{3, 4}
	if !eq(v.Len(), 5) {
		t.Fatalf("len: %v", v.Len())
	}
	n := v.Norm()
	if !eq(n.Len(), 1) {
		t.Fatalf("norm len: %v", n.Len())
	}
	if (Vec2{}).Norm() != (Vec2{}) {
		t.Fatal("zero norm must be zero")
	}
}

func TestVec_Dot(t *testing.T) {
	if (Vec2{1, 0}).Dot(Vec2{0, 1}) != 0 {
		t.Fatal("orthogonal dot must be 0")
	}
	if (Vec2{2, 3}).Dot(Vec2{4, -1}) != 5 {
		t.Fatal("dot")
	}
}
