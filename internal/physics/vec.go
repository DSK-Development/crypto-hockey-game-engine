package physics

import "math"

type Vec2 struct {
	X, Y float64
}

func (a Vec2) Add(b Vec2) Vec2   { return Vec2{a.X + b.X, a.Y + b.Y} }
func (a Vec2) Sub(b Vec2) Vec2   { return Vec2{a.X - b.X, a.Y - b.Y} }
func (a Vec2) Scale(s float64) Vec2 { return Vec2{a.X * s, a.Y * s} }
func (a Vec2) Dot(b Vec2) float64   { return a.X*b.X + a.Y*b.Y }
func (a Vec2) Len() float64         { return math.Sqrt(a.X*a.X + a.Y*a.Y) }
func (a Vec2) Norm() Vec2 {
	l := a.Len()
	if l == 0 {
		return Vec2{}
	}
	return Vec2{a.X / l, a.Y / l}
}
