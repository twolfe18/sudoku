
package main

import (
	"image"
	"image/draw"
	"math"
	"fmt"
	"os"
)

/******************************************************************************************/

type Line struct {
	left, right Float64Point
	radius float64		// std deviation of gaussian off the normal of the line
}

/******************************************************************************************/

func (l Line) Midpoint() (mid Float64Point) {
	v := PointMinus(l.right, l.left)
	v.Scale(0.5)
	return PointPlus(l.left, v)
}

func (l Line) ScaleLength(scale float64) {
	m := Midpoint(l)
	v := PointMinus(l.right, m)
	v.Scale(scale)
	l.right = PointPlus(m, v)
	l.left = PointMinus(m, v)
}

func (l Line) Shift(dx, dy float64) {
	l.left.Shift(dx, dy)
	l.right.Shift(dx, dy)
}

func (l Line) ProjectInto(bounds Float64Rectangle) {
	l.left.ProjectInto(bounds)
	l.right.ProjectInto(bounds)
}

// you might think to have a function that returns a chan image.Point
// but i think this is cleaner because this way it is clear that the
// channel must be created and closed in the same place (as opposed to
// created in this function and closed by the caller)
func (l Line) Iterator(c chan image.Point) {
	cur := l.left
	iter := int(math.Fmax(math.Fabs(v.X), math.Fabs(v.Y)))
	if iter == 0 {
		c <- image.Point{int(l.left.X), int(l.right.Y)}
		return
	}
	dx := v.X / float64(iter); dy := v.Y / float64(iter)
	for i := 0; i<iter; i++ {
		c <- image.Point{int(cur.X), int(cur.Y)}
	}
}

// this allows for stuff like anti-aliased drawing
type WeightedPoint struct {
	P image.Point
	W float64
}

func (l Line) WeightedIterator(c chan WeightedPoint) {
	// TODO
	panic("[WeightedIterator] not implemented yet!")
}


/******************************************************************************************/

func (l Line) String() string {
	return fmt.Sprintf("[%s -> %s]", l.left.String(), l.right.String())
}


func (l Line) Draw(img draw.Image, c image.Color) {
	ch := make(chan image.Point)
	l.LineIter(ch)
	for p := <-ch { img.Set(p.X, p.Y, c) }
	close(ch)
}

func (l Line) Angle(o Line) float64 {
	v1 := PointMinus(o.right, o.left)
	v2 := PointMinus(l.right, l.left)
	//fmt.Printf("[Angle] dp = %.2f\n", DotProduct(v1, v2))
	switch d := math.Acos(DotProduct(v1, v2) / v1.L2Norm() / v2.L2Norm()) * 180.0 / math.Pi; {
	case 0 <= d && d < 90.0:
		return d
	case 90 <= d && d < 180.0:
		return 180.0 - d
	default:
		fmt.Printf("Line.Angle] wut?\td = %.2f\n", d)
		os.Exit(1)
	}
	return math.NaN()
}

func (l Line) SquaredDistance(x, y float64) float64 {
	// http://paulbourke.net/geometry/pointline/
	u := (x - l.left.X) * (l.right.X - l.left.X)
	u += (y - l.left.Y) * (l.right.Y - l.left.Y)
	u /= math.Pow((l.right.X - l.left.X), 2.0) + math.Pow((l.right.Y - l.left.Y), 2.0)
	sx := l.left.X + u * (l.right.X - l.left.X)
	sy := l.left.Y + u * (l.right.Y - l.left.Y)
	return math.Pow(x-sx, 2.0) + math.Pow(y-sy, 2.0)
}



