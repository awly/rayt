package main

import (
	"testing"
)

const valMargin = 0.001

func equals(a, b float64) bool {
	return a+valMargin > b && a-valMargin < b
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
