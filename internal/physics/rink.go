package physics

type Slot int

const (
	SlotA Slot = 0
	SlotB Slot = 1
)

type Rink struct {
	Left, Right, Top, Bottom float64
	CenterX                  float64
	GoalTop, GoalBottom      float64
	Restitution              float64
}

func DefaultRink() Rink {
	return Rink{
		Left: 0, Right: 800, Top: 0, Bottom: 400,
		CenterX:    400,
		GoalTop:    140,
		GoalBottom: 260,
		Restitution: 0.92,
	}
}

func (r Rink) ClampMallet(s Slot, target Vec2, radius float64) Vec2 {
	minX, maxX := r.Left+radius, r.CenterX-radius
	if s == SlotB {
		minX, maxX = r.CenterX+radius, r.Right-radius
	}
	x := clamp(target.X, minX, maxX)
	y := clamp(target.Y, r.Top+radius, r.Bottom-radius)
	return Vec2{x, y}
}

func clamp(v, lo, hi float64) float64 {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}
