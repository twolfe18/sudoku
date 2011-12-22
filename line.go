
package sudoku

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
	Draw(image.Image)
	Angle(Line) float64
	SquaredDistance(float64, float64) float64
}

type PointLine struct {
	// TODO image.Point is discrete, need continuous version
	tl, br Float64Point	// top left and bottom right
	radius float64		// std deviation of gaussian off the normal of the line
}

type RadialLine struct {
	midpoint Float64Point
	rotation, length, radius float64
}

/******************************************************************************************/

func Radial2Point(line RadialLine) PointLine {
	// TODO
}

func Point2Radial(line PointLine) RadialLine {
	// TODO
}

/******************************************************************************************/

func (l PointLine) Draw(img draw.Image) {
	hl_color := image.RGBAColor{255, 0, 0, 255}
	px := float64(l.tl.X)
	py := float64(l.tl.Y)
	it := float64(max(l.br.X - l.tl.X, l.br.Y - l.tl.Y))
	dx := (float64(l.br.X) - px) / it
	dy := (float64(l.br.Y) - py) / it
	for {
		img.Set(int(px), int(py), hl_color)
		px += dx
		py += dy
		if int(px) > l.br.X || int(py) > l.br.Y {
			break
		}
	}
}

func (l PointLine) Angle(other Line) float64 {
	// TODO do a native version
	return Point2Radial(l).Angle(other)
}

// TODO test this!
func (l PointLine) SquaredDistance(x, y float64) float64 {
	// http://paulbourke.net/geometry/pointline/
	u := (x - l.tl.X) * (l.br.X - l.tl.X)
	u += (y - l.tl.Y) * (l.br.Y - l.tl.Y)
	u /= float64((l.br.X - l.tl.X)^2 + (l.br.Y - l.tl.Y)^2)
	sx := float64(l.tl.X) + u * float64(l.br.X - l.tl.X)
	sy := float64(l.tl.Y) + u * float64(l.br.Y - l.tl.Y)
	return math.Pow((float64(x)-sx), 2.0) + math.Pow((float64(y)-sy), 2.0)
}

/******************************************************************************************/

func (l RadialLine) Draw(img draw.Image) {
	// TODO do a native version
	Radial2Point(l).Draw()
}

func (l RadialLine) Angle(other Line) float64 {
	var o RadialLine
	if other.(type) == RadialLine { o = other.(RadialLine) }
	else { o = Point2Radial(other) }
	return math.Fmod(l.rotation - o.rotation, 180.0)
}

func (l RadialLine) SquaredDistance(x, y float64) float64 {
	// TODO do a native version
	// TODO alternatively, could redefine in polar terms
		// degree of deflexion required to hit point
	return Radial2Point(l).SquaredDistance(x, y)
}




