
package main

import (
	"os"
	"image"
	"image/png"
	"fmt"
	"rand"
	"math"
)

type Float64Point struct {
	X, Y float64
}

func (p Float64Point) String() string {
	return fmt.Sprintf("(%.1f, %.1f)", p.X, p.Y)
}

func (p *Float64Point) Scale(s float64) {
	p.X *= s
	p.Y *= s
}

func (p *Float64Point) Shift(dx, dy float64) {
	p.X += dx
	p.Y += dy
}

func (p Float64Point) ProjectInto(bounds Float64Rectangle) {
	p.X = math.Fmax(bounds.Min.X, p.X)
	p.X = math.Fmin(bounds.Max.X, p.X)
	p.Y = math.Fmax(bounds.Min.Y, p.Y)
	p.Y = math.Fmin(bounds.Max.Y, p.Y)
}

func (v Float64Point) Rotate(theta float64) Float64Point {
	// http://en.wikipedia.org/wiki/Rotation_(mathematics)#Matrix_algebra
	st := math.Sin(theta)
	ct := math.Cos(theta)
	xp := v.X * ct - v.Y * st
	yp := v.X * st + v.Y * ct
	return Float64Point{xp, yp}
}

func (v Float64Point) L2Norm() float64 {
	return math.Sqrt(v.X * v.X + v.Y * v.Y)
}

func PointMinus(a, b Float64Point) (r Float64Point) {
	r.X = a.X - b.X
	r.Y = a.Y - b.Y
	return r
}

func PointPlus(a, b Float64Point) (r Float64Point) {
	r.X = a.X + b.X
	r.Y = a.Y + b.Y
	return r
}

func DotProduct(a, b Float64Point) float64 {
	return a.X * b.X + a.Y * b.Y
}

func RandomPointBetween(lo, hi Float64Point) Float64Point {
	x := float64(hi.X) - rand.Float64() * float64(hi.X - lo.X)
	y := float64(hi.Y) - rand.Float64() * float64(hi.Y - lo.Y)
	return Float64Point{x, y}
}

