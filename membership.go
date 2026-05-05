package autoscaler

type MembershipFunction interface {
	Membership(x float64) float64
}

type Trapezoid struct {
	A, B, C, D float64
}

func (t Trapezoid) Membership(x float64) float64 {
	if x <= t.A || x >= t.D {
		if x == t.A && t.A == t.B {
			return 1.0
		}
		if x == t.D && t.C == t.D {
			return 1.0
		}
		return 0.0
	}
	if x <= t.B {
		if t.B == t.A {
			return 1.0
		}
		return (x - t.A) / (t.B - t.A)
	}
	if x >= t.C {
		if t.D == t.C {
			return 1.0
		}
		return (t.D - x) / (t.D - t.C)
	}
	return 1.0
}

type Triangle struct {
	A, B, C float64
}

func (t Triangle) Membership(x float64) float64 {
	if x <= t.A || x >= t.C {
		return 0.0
	}
	if x == t.B {
		return 1.0
	}
	if x > t.A && x < t.B {
		return (x - t.A) / (t.B - t.A)
	}
	return (t.C - x) / (t.C - t.B)
}

func clamp(x, min, max float64) float64 {
	if x < min {
		return min
	}
	if x > max {
		return max
	}
	return x
}
