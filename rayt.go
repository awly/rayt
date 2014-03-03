// TODO:
// - write tests for distrp
// - write tests for ray.projp
// - write tests for sphere.intersect
// - make surface "splitable"
// - create mechanism to produce rays (from eye trough surface)
// - create func to update ray color based on colision
// - fill surface values based on rays
// - write tests for all new funcs
// - parallelize
// - read scene info from json file
// - cleanup

package main

import (
	"image"
	"image/color"
	"image/png"
	"log"
	"math"
	"os"
)

const (
	width  = 1000
	height = 1000
)

func main() {
	log.Println("starting")

	sc, surf, err := readInput()
	if err != nil {
		log.Fatalln(err)
	}
	draw(sc, &surf)
	save(surf, "out.png")

	log.Println("done")
}

func readInput() (scene, surface, error) {
	sc := scene{
		objs: []obj{
			sphere{rad: 10},
		},
		eye:   point{x: 30, y: 30, z: 30},
		light: point{x: 100, y: 0, z: 0},
	}

	surf := surface{
		grid: make([][]color.RGBA, height),
	}
	for i := range surf.grid {
		surf.grid[i] = make([]color.RGBA, width)
		for j := range surf.grid[i] {
			surf.grid[i][j].A = 255
			surf.grid[i][j].R = uint8(i*i + j*j)
			surf.grid[i][j].G = uint8(i*i + j*j)
			surf.grid[i][j].B = uint8(i*i + j*j)
		}
	}

	return sc, surf, nil
}

func draw(sc scene, surf *surface) {}

func save(surf surface, fname string) {
	out, err := os.Create(fname)
	if err != nil {
		log.Println(err)
		return
	}
	defer out.Close()

	if err = png.Encode(out, surf); err != nil {
		log.Println(err)
	}
}

type scene struct {
	objs       []obj
	eye, light point
}

type surface struct {
	grid [][]color.RGBA
}

// implements image.Image
func (s surface) ColorModel() color.Model { return color.RGBAModel }
func (s surface) Bounds() image.Rectangle { return image.Rect(0, 0, len(s.grid[0]), len(s.grid)) }
func (s surface) At(x, y int) color.Color { return s.grid[y][x] }

type obj interface {
	intersect(l ray) []point
}

type point struct {
	x, y, z float64
}

// distance from zero
func (p point) len() float64 {
	return math.Sqrt(p.x*p.x + p.y*p.y + p.z*p.z)
}

// deistance between two points
func distpp(a, b point) float64 {
	return math.Sqrt(math.Pow(a.x-b.x, 2) + math.Pow(a.y-b.y, 2) + math.Pow(a.z-b.z, 2))
}

// distance between ray and point
func distrp(l ray, p point) float64 {
	return distpp(p, l.projp(p))
}

// vec - from point zero, not start
type ray struct {
	start, vec point
}

// point projected onto a ray
func (l ray) projp(p point) point {
	u := ray{
		start: l.start,
		vec: point{
			x: p.x - l.start.x,
			y: p.y - l.start.y,
			z: p.z - l.start.z,
		},
	}
	dpvu := dotProd(u, l)
	// point does not project on ray directly
	if dpvu < 0 {
		return l.start
	}
	lv := l.vec.len()
	c := dpvu / lv

	return point{
		x: p.x + l.vec.x*c,
		y: p.y + l.vec.y*c,
		z: p.z + l.vec.z*c,
	}
}

type sphere struct {
	center point
	rad    float64
}

func (s sphere) intersect(l ray) []point {
	dcl := distrp(l, s.center)

	// no intersection exists
	if dcl > s.rad {
		return nil
	}
	// only one point, which is projection of sphere center onto ray
	if dcl == 0 {
		// TODO avoid recomputation of projection?
		return []point{l.projp(s.center)}
	}

	// distance from sphere center projection to either one of intersection points
	di := math.Sqrt(math.Pow(s.rad, 2) - math.Pow(dcl, 2))
	// distance from ray start to sphere center projection point
	dlp := distpp(l.start, l.projp(s.center))
	// length of ray vector (in case it's not normalized)
	ll := l.vec.len()

	return []point{
		point{
			x: l.start.x + l.vec.x*(dlp-di)/ll,
			y: l.start.y + l.vec.y*(dlp-di)/ll,
			z: l.start.z + l.vec.z*(dlp-di)/ll,
		},
		point{
			x: l.start.x + l.vec.x*(dlp+di)/ll,
			y: l.start.y + l.vec.y*(dlp+di)/ll,
			z: l.start.z + l.vec.z*(dlp+di)/ll,
		},
	}
}

func dotProd(a, b ray) float64 {
	return a.vec.x*b.vec.x + a.vec.y*b.vec.y + a.vec.z*b.vec.z
}
