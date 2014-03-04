package main

import (
	"testing"
)

const valMargin = 0.001

func equals(a, b float64) bool {
	return a+valMargin > b && a-valMargin < b
}
func equalsp(a, b point) bool {
	return equals(a.x, b.x) && equals(a.y, b.y) && equals(a.z, b.z)
}

type distppcase struct {
	a, b point
	res  float64
}

var distppcases = []distppcase{
	distppcase{a: point{0, 0, 0}, b: point{0, 0, 0}, res: 0},
	distppcase{a: point{1, 0, 0}, b: point{0, 0, 0}, res: 1},
	distppcase{a: point{1, 1, 0}, b: point{0, 0, 0}, res: 1.41421},
	distppcase{a: point{1, 1, 1}, b: point{0, 0, 0}, res: 1.73205},
	distppcase{a: point{1, 0, 0}, b: point{0, 1, 0}, res: 1.41421},
}

func TestDistpp(t *testing.T) {
	for _, c := range distppcases {
		if got := distpp(c.a, c.b); !equals(got, c.res) {
			t.Errorf("distpp(%v, %v) = %f; got: %f", c.a, c.b, c.res, got)
		}
	}
}

type distrpcase struct {
	a   ray
	b   point
	res float64
}

var distrpcases = []distrpcase{
	distrpcase{a: ray{start: point{0, 0, 0}, vec: point{1, 0, 0}}, b: point{0, 0, 0}, res: 0},
	distrpcase{a: ray{start: point{0, 0, 0}, vec: point{2, 0, 0}}, b: point{1, 1, 0}, res: 1},
	distrpcase{a: ray{start: point{0, 0, 0}, vec: point{2, 0, 0}}, b: point{1, 1, 1}, res: 1.41421},
	distrpcase{a: ray{start: point{0, 0, 0}, vec: point{1, 0, 0}}, b: point{2, 1, 1}, res: 1.41421},
}

func TestDistrp(t *testing.T) {
	for _, c := range distrpcases {
		if got := distrp(c.a, c.b); !equals(got, c.res) {
			t.Errorf("distrp(%v, %v) = %f; got: %f", c.a, c.b, c.res, got)
		}
	}
}

type projpcase struct {
	a   ray
	b   point
	res point
}

var projpcases = []projpcase{
	projpcase{a: ray{start: point{0, 0, 0}, vec: point{1, 0, 0}}, b: point{0, 0, 0}, res: point{0, 0, 0}},
	projpcase{a: ray{start: point{0, 0, 0}, vec: point{2, 0, 0}}, b: point{1, 1, 0}, res: point{1, 0, 0}},
	projpcase{a: ray{start: point{0, 0, 0}, vec: point{2, 0, 0}}, b: point{1, 1, 1}, res: point{1, 0, 0}},
	projpcase{a: ray{start: point{0, 0, 0}, vec: point{1, 0, 0}}, b: point{3, 1, 1}, res: point{3, 0, 0}},
}

func TestProjp(t *testing.T) {
	for _, c := range projpcases {
		if got := c.a.projp(c.b); !equalsp(got, c.res) {
			t.Errorf("%v.projp(%v) = %v; got: %v", c.a, c.b, c.res, got)
		}
	}
}

type intersectcase struct {
	a   obj
	b   ray
	res []point
}

var intersectcases = []intersectcase{
	intersectcase{a: sphere{center: point{0, 0, 0}, rad: 1}, b: ray{start: point{1, 1, 0}, vec: point{-1, 0, 0}}, res: []point{point{0, 1, 0}}},
	intersectcase{a: sphere{center: point{0, 0, 0}, rad: 1}, b: ray{start: point{1, 1, 0}, vec: point{0, -1, 0}}, res: []point{point{1, 0, 0}}},
	intersectcase{a: sphere{center: point{0, 0, 0}, rad: 1}, b: ray{start: point{1, 1, 0}, vec: point{0, 0, 1}}, res: nil},
	intersectcase{a: sphere{center: point{0, 0, 0}, rad: 1}, b: ray{start: point{2, 0, 0}, vec: point{-1, 0, 0}}, res: []point{point{1, 0, 0}, point{-1, 0, 0}}},
}

func TestIntersect(t *testing.T) {
	for _, c := range intersectcases {
		got := c.a.intersect(c.b)
		if len(got) != len(c.res) {
			t.Errorf("%v.intersect(%v) = %v; got: %v", c.a, c.b, c.res, got)
			return
		}
		for i, v := range got {
			if !equalsp(v, c.res[i]) {
				t.Errorf("%v.intersect(%v) = %v; got: %v", c.a, c.b, c.res, got)
				break
			}
		}
	}
}
