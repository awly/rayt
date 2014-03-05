// TODO:
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
	width   = 1000
	height  = 1000
	ambient = 0.1
)

func main() {
	log.Println("starting")

	sc, v, err := readInput()
	if err != nil {
		log.Fatalln(err)
	}
	draw(sc, &v)
	save(v, "out.png")

	log.Println("done")
}

func readInput() (scene, view, error) {
	sc := scene{
		objs: []obj{
			sphere{
				center: point{0, 0, 0},
				rad:    15,
				c:      color.RGBA{A: 255, R: 50, G: 200, B: 200},
			},
		},
		eye:   point{x: 0, y: 0, z: 50},
		light: point{x: 100, y: 100, z: 100},
	}

	v := view{
		c:  make([][]color.RGBA, height),
		e1: point{10, 10, 25},
		e2: point{10, -10, 25},
		e3: point{-10, -10, 25},
		e4: point{-10, 10, 25},
	}
	for i := range v.c {
		v.c[i] = make([]color.RGBA, width)
		for j := range v.c[i] {
			v.c[i][j].A = 255
			v.c[i][j].R = 50
			v.c[i][j].G = 50
			v.c[i][j].B = 50
		}
	}

	return sc, v, nil
}

func draw(sc scene, v *view) {
	v.foreach(func(x, y int, p point) {
		r := ray{start: sc.eye, vec: point{p.x - sc.eye.x, p.y - sc.eye.y, p.z - sc.eye.z}}
		var fobj obj
		var fp point
		mind := math.MaxFloat64
		for _, v := range sc.objs {
			hits := v.intersect(r)
			// TODO no intersections found. check scene and view creation.
			for _, p := range hits {
				d := distpp(p, r.start)
				if d < mind {
					fobj = v
					fp = p
					mind = d
				}
			}
		}
		if fobj == nil {
			// no intersections, keep default color
			return
		}
		v.c[y][x] = fobj.rayc(fp, sc.light)
	})
}

func save(v view, fname string) {
	out, err := os.Create(fname)
	if err != nil {
		log.Println(err)
		return
	}
	defer out.Close()

	if err = png.Encode(out, v); err != nil {
		log.Println(err)
	}
}

type scene struct {
	objs       []obj
	eye, light point
}

type view struct {
	e1, e2, e3, e4 point
	c              [][]color.RGBA
}

func (v view) sub(x1, x2, y1, y2 int) view {
	res := v
	res.c = res.c[y1:y2]
	for i := range res.c {
		res.c[i] = res.c[i][x1:x2]
	}
	return res
}

func (v view) foreach(f func(int, int, point)) {
	w := len(v.c) - 1
	h := len(v.c[0]) - 1
	for i := 0; i < h+1; i++ {
		c1 := float64(i) / float64(h)
		a := point{
			x: v.e2.x + (v.e1.x-v.e2.x)*c1,
			y: v.e2.y + (v.e1.y-v.e2.y)*c1,
			z: v.e2.z + (v.e1.z-v.e2.z)*c1,
		}
		b := point{
			x: v.e3.x + (v.e4.x-v.e3.x)*c1,
			y: v.e3.y + (v.e4.y-v.e3.y)*c1,
			z: v.e3.z + (v.e4.z-v.e3.z)*c1,
		}
		for j := 0; j < w+1; j++ {
			c2 := float64(j) / float64(w)
			f(i, j, point{
				x: a.x + (b.x-a.x)*c2,
				y: a.y + (b.y-a.y)*c2,
				z: a.z + (b.z-a.z)*c2,
			})
		}
	}
}

// implements image.Image
func (s view) ColorModel() color.Model { return color.RGBAModel }
func (s view) Bounds() image.Rectangle { return image.Rect(0, 0, len(s.c[0]), len(s.c)) }
func (s view) At(x, y int) color.Color { return s.c[y][x] }

type obj interface {
	intersect(l ray) []point
	rayc(p, l point) color.RGBA
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
	lv := math.Pow(l.vec.len(), 2)
	c := dpvu / lv

	return point{
		x: l.start.x + l.vec.x*c,
		y: l.start.y + l.vec.y*c,
		z: l.start.z + l.vec.z*c,
	}
}

type sphere struct {
	center point
	rad    float64
	c      color.RGBA
}

func (s sphere) intersect(l ray) []point {
	dcl := distrp(l, s.center)

	// no intersection exists
	if dcl > s.rad {
		return nil
	}
	// only one point, which is projection of sphere center onto ray
	if dcl == s.rad {
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

func (s sphere) rayc(p, l point) color.RGBA {
	lenr := distpp(l, p)
	lenn := distpp(p, s.center)

	rnorm := ray{start: p, vec: point{
		x: (l.x - p.x) / lenr,
		y: (l.y - p.y) / lenr,
		z: (l.z - p.z) / lenr,
	}}
	nnorm := ray{start: p, vec: point{
		x: (p.x - s.center.x) / lenn,
		y: (p.y - s.center.y) / lenn,
		z: (p.z - s.center.z) / lenn,
	}}

	shade := dotProd(rnorm, nnorm)
	if shade < 0 {
		shade = 0
	}

	res := s.c
	res.R = uint8(float64(res.R) * (ambient + (1-ambient)*shade))
	res.G = uint8(float64(res.G) * (ambient + (1-ambient)*shade))
	res.B = uint8(float64(res.B) * (ambient + (1-ambient)*shade))

	return res
}

func dotProd(a, b ray) float64 {
	return a.vec.x*b.vec.x + a.vec.y*b.vec.y + a.vec.z*b.vec.z
}
