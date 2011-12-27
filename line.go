
package main

import (
	"image"
	"image/draw"
	"math"
	"fmt"
)

/******************************************************************************************/

type Line struct {
	left, right Float64Point
	radius float64		// std deviation of gaussian off the normal of the line
}

func (l Line) Equals(o Line) bool {
	const ep = 1e-4
	return math.Fabs(l.radius-o.radius) < ep && l.left.Equals(o.left) && l.right.Equals(o.right)
}

func HorizontalLine() Line{
	return Line{Float64Point{0.0, 0.0}, Float64Point{1.0, 0.0}, 0.0}
}

func VerticalLine() Line{
	return Line{Float64Point{0.0, 0.0}, Float64Point{0.0, 1.0}, 0.0}
}

/******************************************************************************************/

func (l Line) Dx() float64 {
	return l.right.X - l.left.X
}

func (l Line) Dy() float64 {
	return l.right.Y - l.left.Y
}

func (l Line) Midpoint() (mid Float64Point) {
	v := PointMinus(l.right, l.left)
	v.Scale(0.5)
	return PointPlus(l.left, v)
}

func (l *Line) ScaleLength(scale float64) {
	m := l.Midpoint()
	v := PointMinus(l.right, m)
	v.Scale(scale)
	l.right = PointPlus(m, v)
	l.left = PointMinus(m, v)
}

func (l *Line) Shift(dx, dy float64) {
	l.left.Shift(dx, dy)
	l.right.Shift(dx, dy)
}

func (l *Line) ProjectInto(bounds Float64Rectangle) {
	l.left.ProjectInto(bounds)
	l.right.ProjectInto(bounds)
}

// this allows for stuff like anti-aliased drawing
type WeightedPoint struct {
	P image.Point
	W float64
}

// it would be nice to use goroutines for the interator, but
// it appears to be way too slow:
// http://groups.google.com/group/golang-nuts/browse_thread/thread/a717e1286a8736fd#
// for now i'll use a fully blocking producer-consumer model
func (l Line) UnweightedIterator() (pix []WeightedPoint) {

	cur := l.left
	iter := int(math.Fmax(math.Fabs(l.Dx()), math.Fabs(l.Dy())))
	if iter == 0 {
		p := image.Point{int(l.left.X), int(l.right.Y)}
		return append(pix, WeightedPoint{p, 1.0})
	}
	dx := l.Dx() / float64(iter); dy := l.Dy() / float64(iter)
	for i := 0; i<iter; i++ {
		p := image.Point{int(cur.X), int(cur.Y)}
		pix = append(pix, WeightedPoint{p, 1.0})
		cur.X += dx; cur.Y += dy
	}
	return pix
}

func (l Line) WeightedIterator() (pix []WeightedPoint) {

	max_delta := 2.0 * l.radius	// will miss prob mass outside of 2 std dev
	normalizer := 1.0 / math.Sqrt(2.0 * math.Pi * l.radius * l.radius)

	var p image.Point
	cur := l.left
	iter := int(math.Fmax(math.Fabs(l.Dx()), math.Fabs(l.Dy()))) + 1
	dx := l.Dx() / float64(iter); dy := l.Dy() / float64(iter)

	for i := 0; i<iter; i++ {
		for d := -max_delta; d < max_delta; d += 1.0 {

			if dx > dy {	// vertical sweeps
				p = image.Point{int(cur.X), int(cur.Y + d)}
			} else {	// horizontal sweeps
				p = image.Point{int(cur.X + d), int(cur.Y)}
			}

			fp := NewFloat64Point(p)
			dist := l.Distance(fp.X, fp.Y)
			dist = math.Fmin(dist, Distance(fp, l.left))
			dist = math.Fmin(dist, Distance(fp, l.right))
			weight := normalizer * math.Exp(-dist * dist / (2.0 * l.radius * l.radius))

			pix = append(pix, WeightedPoint{p, weight})
		}
		cur.X += dx; cur.Y += dy
	}
	return pix
}


/******************************************************************************************/

func (l Line) String() string {
	return fmt.Sprintf("[%s -> %s]", l.left.String(), l.right.String())
}

func (l Line) Draw(img draw.Image, c image.RGBAColor) {
	for _,wp := range l.WeightedIterator() {
		c.A = uint8(wp.W * 255.0)
		img.Set(wp.P.X, wp.P.Y, c)
	}
}

func (l Line) Angle(o Line) float64 {
	v1 := PointMinus(o.right, o.left)
	v2 := PointMinus(l.right, l.left)
	switch d := math.Acos(DotProduct(v1, v2) / v1.L2Norm() / v2.L2Norm()) * 180.0 / math.Pi; {
	case 0 <= d && d < 90.0:
		return d
	case 90 <= d && d < 180.0:
		return 180.0 - d
	default:
		panic(fmt.Sprintf("Line.Angle] wut?\td = %.2f\n", d))
	}
	return math.NaN()
}

func (l *Line) Rotate(theta float64) {
	v := PointMinus(l.right, l.left)
	v.Scale(0.5)
	v.Rotate(theta)
	m := l.Midpoint()
	l.left = PointMinus(m, v)
	l.right = PointPlus(m, v)
}

func (l Line) Distance(x, y float64) float64 {
	// http://paulbourke.net/geometry/pointline/
	u := (x - l.left.X) * (l.right.X - l.left.X)
	u += (y - l.left.Y) * (l.right.Y - l.left.Y)
	u /= math.Pow((l.right.X - l.left.X), 2.0) + math.Pow((l.right.Y - l.left.Y), 2.0)
	sx := l.left.X + u * (l.right.X - l.left.X)
	sy := l.left.Y + u * (l.right.Y - l.left.Y)
	return math.Sqrt((x-sx)*(x-sx) + (y-sy)*(y-sy))
}



