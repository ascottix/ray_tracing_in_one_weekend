package main

import (
	"testing"
)

func TestConstructor(t *testing.T) {
	e := Vec3{}

	if e.X != 0.0 || e.Y != 0.0 || e.Z != 0.0 {
		t.Errorf("%v should be zero", e)
	}

	v := NewVec3(1.1, 2.2, 3.3)

	if v.X != 1.1 || v.Y != 2.2 || v.Z != 3.3 {
		t.Errorf("%v does not match initialization params", v)
	}

	p := NewPoint3(1, 1, 1)
	if p.LengthSquared() != 3 {
		t.Errorf("%v does not match initialization params", v)
	}

}

func TestBasicOperations(t *testing.T) {
	a := NewVec3(1, 2, 3)
	b := NewVec3(4, 5, 6)

	// Addition
	s := a.Add(b)
	if s.X != 5 || s.Y != 7 || s.Z != 9 {
		t.Errorf("%v + %v != %v", a, b, s)
	}

	// Multiplication
	f := 2.0
	m := a.Mul(f)
	if m.X != 2 || m.Y != 4 || m.Z != 6 {
		t.Errorf("%v * %f != %v", a, f, m)
	}

	// Length
	if a.LengthSquared() != 1+4+9 {
		t.Errorf("%v.LengthSquared() = %f mismatch", a, a.LengthSquared())
	}

	l := NewVec3(0, 3, 4).Length()
	if l != 5 {
		t.Errorf("Length() = %f mismatch", l)
	}

	// Dot product
	d := a.Dot(b)
	if d != 4+10+18 {
		t.Errorf("%v * %v = %f mismatch", a, b, d)
	}

	// Cross product
	c := a.Cross(b)
	if c.X != -3 || c.Y != 6 || c.Z != -3 {
		t.Errorf("%v x %v = %v mismatch", a, b, c)
	}

	x := NewVec3(-1, -2, 3).Cross(NewVec3(4, 0, -8))
	if x.X != 16 || x.Y != 4 || x.Z != 8 {
		t.Errorf("cross product mismatch: %v", x)
	}
}
