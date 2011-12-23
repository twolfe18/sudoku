
package main

import (
	"image"
	"image/draw"
	"math"
	"fmt"
	"os"
)

/******************************************************************************************/

// TODO not necessary now
/* type RadialLineVariance struct {
	Dx, Dy float64	// for midpoint
	Dtheta float64	// for RadialLine.rotation
	Dlength float64	// duh
} */

type Line interface {
	Draw(draw.Image)
	Angle(Line) float64
	SquaredDistance(float64, float64) float64
}

type PointLine struct {
	left, right Float64Point
	radius float64		// std deviation of gaussian off the normal of the line
}

type RadialLine struct {
	midpoint Float64Point
	rotation, length, radius float64
}

/******************************************************************************************/

func Radial2Point(line RadialLine) (pl PointLine) {
	a := line.length * math.Cos(line.rotation) / 2.0
	pl.left.X = line.midpoint.X - a
	pl.right.X = line.midpoint.X + a
	b := line.length * math.Sin(line.rotation) / 2.0
	pl.left.Y = line.midpoint.Y - b
	pl.right.Y = line.midpoint.Y + b
	pl.radius = line.radius
	return pl
}

func Point2Radial(line PointLine) (rl RadialLine) {
	rl.midpoint = Midpoint(line.left, line.right)
	rec := Float64Rectangle{line.left, line.right}	// Dx() always >= 0
	rl.length = math.Sqrt(math.Pow(rec.Dx(), 2.0) + math.Pow(rec.Dy(), 2.0))
	if rec.Dy() >= 0.0 {
		rl.rotation = math.Atan2(rec.Dy(), rec.Dx())
	} else {
		rl.rotation = 90.0 - math.Atan2(-1.0 * rec.Dy(), rec.Dx())
	}
	rl.radius = line.radius
	return rl
}

/******************************************************************************************/

func (l PointLine) Draw(img draw.Image) {
	hl_color := image.RGBAColor{255, 0, 0, 255}
	px := l.left.X
	py := l.right.Y
	it := math.Fmax(l.right.X - l.left.X, math.Fabs(l.left.Y - l.right.Y))
	dx := (l.right.X - l.left.X) / it
	dy := (l.right.Y - l.left.Y) / it
	for {
		img.Set(int(px), int(py), hl_color)
		px += dx
		py += dy
		if int(px) > int(l.right.X) { break }
	}
}

func (l PointLine) Angle(other Line) float64 {
	// TODO do a native version
	return Point2Radial(l).Angle(other)
}

// TODO test this!
func (l PointLine) SquaredDistance(x, y float64) float64 {
	// http://paulbourke.net/geometry/pointline/
	u := (x - l.left.X) * (l.right.X - l.left.X)
	u += (y - l.left.Y) * (l.right.Y - l.left.Y)
	u /= math.Pow((l.right.X - l.left.X), 2.0) + math.Pow((l.right.Y - l.left.Y), 2.0)
	sx := l.left.X + u * (l.right.X - l.left.X)
	sy := l.left.Y + u * (l.right.Y - l.left.Y)
	return math.Pow(x-sx, 2.0) + math.Pow(y-sy, 2.0)
}

/******************************************************************************************/

func (l RadialLine) Draw(img draw.Image) {
	// TODO do a native version
	Radial2Point(l).Draw(img)
}

func (l RadialLine) Angle(other Line) float64 {
	var o RadialLine
	switch other.(type) {
	case RadialLine:
		o = other.(RadialLine)
	case PointLine:
		o = Point2Radial(other.(PointLine))
	default:
		fmt.Printf("[RadialLine.Angle] unanticipated Line type\n")
		os.Exit(1)
	}
	return math.Fmod(l.rotation - o.rotation, 180.0)
}

func (l RadialLine) SquaredDistance(x, y float64) float64 {
	// TODO do a native version
	// TODO alternatively, could redefine in polar terms
		// degree of deflexion required to hit point
	return Radial2Point(l).SquaredDistance(x, y)
}




