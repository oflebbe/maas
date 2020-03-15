// Package mandel is a package to generate mandelbrot images
package mandel

import (
	"image"
	"image/color"
	"math/cmplx"
)

// Mandel holds essential infos
type Mandel struct {
	maxIter    int
	colorTable []color.RGBA
	comm       chan []int
}

// NewMandel creates a new instance max is the iteration depth
func NewMandel(max int) (m Mandel) {
	m.maxIter = max
	m.colorTable = make([]color.RGBA, max+1)
	for i := 0; i < max; i++ {
		m.colorTable[i] = hSVToRGBA((float32(i)*10.*360.)/float32(max)+0.4,
			0.9-(float32(i)/float32(max))*0.5, 0.75+float32(i)/(4.3*float32(max)))
	}
	m.colorTable[max] = color.RGBA{0, 0, 0, 255}

	m.comm = make(chan []int)
	return
}

func (m *Mandel) point(c complex128) int {
	x := c
	counter := 0
	for counter < m.maxIter && real(x*cmplx.Conj(x)) < 4. {
		x = x*x + c
		counter++
	}
	return counter
}

func (m *Mandel) line(p complex128, d float64, steps int, ident int) {
	ret := make([]int, steps+1)
	ret[0] = ident
	for i := 0; i < steps; i++ {
		x := p + complex(float64(i)*d/float64(steps), 0)
		ret[i+1] = m.point(x)
	}
	m.comm <- ret
}

// Pic computes a mandelbrot tile a p size d with steps steps
func (m *Mandel) Pic(p complex128, d float64, steps int) image.Image {
	ret := image.NewRGBA(image.Rect(0, 0, steps, steps))

	for i := 0; i < steps; i++ {
		x := p + complex(0, float64(i)*d/float64(steps))
		go m.line(x, d, steps, i)
	}

	for i := 0; i < steps; i++ {
		l := <-m.comm
		for j := 1; j <= steps; j++ {
			ret.Set(j-1, steps-l[0]-1, m.colorTable[l[j]])
		}
	}

	return ret
}

func hSVToRGBFloat32(h float32, s float32, v float32) (float32, float32, float32) {
	h_i := int(h / 60.)
	f := h/60. - float32(h_i)
	p := v * (1. - s)
	q := v * (1. - s*f)
	t := v * (1. - s*(1.-f))
	switch h_i {
	case 0, 6:
		return v, t, p
	case 1:
		return q, v, p
	case 2:
		return p, v, t
	case 3:
		return p, q, v
	case 4:
		return t, p, v
	case 5:
		return v, p, q
	}
	return 0, 0, 0
}

func float32Int8(v float32) uint8 {
	v *= 256
	if int(v) >= 255 {
		return 255
	}
	return uint8(v)
}

func hSVToRGBA(h float32, s float32, v float32) color.RGBA {
	r, g, b := hSVToRGBFloat32(h, s, v)
	ru := float32Int8(r)
	gu := float32Int8(g)
	bu := float32Int8(b)
	return color.RGBA{ru, gu, bu, 255}
}
