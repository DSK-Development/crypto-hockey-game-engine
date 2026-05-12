package physics

const puckFriction = 0.995

func StepPuck(p *Body, r Rink, dt float64) {
	p.Pos = p.Pos.Add(p.Vel.Scale(dt))
	p.Vel = p.Vel.Scale(puckFriction)
	resolveWalls(p, r)
}

func resolveWalls(p *Body, r Rink) {
	if p.Pos.Y-p.R < r.Top {
		p.Pos.Y = r.Top + p.R
		p.Vel.Y = -p.Vel.Y * r.Restitution
	}
	if p.Pos.Y+p.R > r.Bottom {
		p.Pos.Y = r.Bottom - p.R
		p.Vel.Y = -p.Vel.Y * r.Restitution
	}
	inGoal := p.Pos.Y >= r.GoalTop && p.Pos.Y <= r.GoalBottom
	if !inGoal {
		if p.Pos.X-p.R < r.Left {
			p.Pos.X = r.Left + p.R
			p.Vel.X = -p.Vel.X * r.Restitution
		}
		if p.Pos.X+p.R > r.Right {
			p.Pos.X = r.Right - p.R
			p.Vel.X = -p.Vel.X * r.Restitution
		}
	}
}

// ResolvePuckMallet resolves a circle-circle collision between the puck (movable)
// and a mallet (treated as kinematic with velocity malletVel). Returns true on contact.
func ResolvePuckMallet(p *Body, malletPos, malletVel Vec2, malletR float64) bool {
	delta := p.Pos.Sub(malletPos)
	dist := delta.Len()
	minDist := p.R + malletR
	if dist >= minDist {
		return false
	}
	if dist == 0 {
		// degenerate — push out along +X
		p.Pos = p.Pos.Add(Vec2{minDist, 0})
		return true
	}
	n := delta.Scale(1.0 / dist)
	overlap := minDist - dist
	p.Pos = p.Pos.Add(n.Scale(overlap))
	rel := p.Vel.Sub(malletVel)
	vn := rel.Dot(n)
	if vn > 0 {
		return true // separating already
	}
	const e = 0.92
	j := -(1 + e) * vn
	p.Vel = p.Vel.Add(n.Scale(j))
	return true
}

// CheckGoal returns "A" if slot A scored (puck went past right goal),
// "B" if slot B scored (past left goal), or "" if no goal.
func CheckGoal(p Body, r Rink) string {
	if p.Pos.Y < r.GoalTop || p.Pos.Y > r.GoalBottom {
		return ""
	}
	if p.Pos.X+p.R < r.Left {
		return "B"
	}
	if p.Pos.X-p.R > r.Right {
		return "A"
	}
	return ""
}
