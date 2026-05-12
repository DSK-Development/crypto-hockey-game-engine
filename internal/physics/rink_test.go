package physics

import "testing"

func TestClampMalletA_LeftHalf(t *testing.T) {
	r := DefaultRink()
	// A is left half, x in [Left+R, CenterX-R]
	got := r.ClampMallet(SlotA, Vec2{1000, 200}, 24)
	if got.X != r.CenterX-24 {
		t.Fatalf("right edge: %+v", got)
	}
	got = r.ClampMallet(SlotA, Vec2{-50, 200}, 24)
	if got.X != r.Left+24 {
		t.Fatalf("left edge: %+v", got)
	}
}

func TestClampMalletB_RightHalf(t *testing.T) {
	r := DefaultRink()
	got := r.ClampMallet(SlotB, Vec2{-100, 200}, 24)
	if got.X != r.CenterX+24 {
		t.Fatalf("left edge: %+v", got)
	}
	got = r.ClampMallet(SlotB, Vec2{1000, 200}, 24)
	if got.X != r.Right-24 {
		t.Fatalf("right edge: %+v", got)
	}
}

func TestClampMallet_VerticalBounds(t *testing.T) {
	r := DefaultRink()
	got := r.ClampMallet(SlotA, Vec2{100, -50}, 24)
	if got.Y != r.Top+24 {
		t.Fatalf("top: %+v", got)
	}
	got = r.ClampMallet(SlotA, Vec2{100, 1000}, 24)
	if got.Y != r.Bottom-24 {
		t.Fatalf("bottom: %+v", got)
	}
}
