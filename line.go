
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

type Line struct {
	left, right Float64Point
	radius float64		// std deviation of gaussian off the normal of the line
}

/******************************************************************************************/

func (l Line) String() string {
	return fmt.Sprintf("[%s -> %s]", l.left.String(), l.right.String())
}

func (l Line) Draw(img draw.Image) {
	hl_color := image.RGBAColor{255, 0, 0, 255}
	v := PointMinus(l.right, l.left)
	cur := l.left
	iter := int(math.Fmax(math.Fabs(v.X), math.Fabs(v.Y)))
	if iter == 0 { img.Set(int(l.left.X), int(l.right.Y), hl_color) }
	dx := v.X / float64(iter); dy := v.Y / float64(iter)
	for i := 0; i<iter; i++ {
		img.Set(int(cur.X), int(cur.Y), hl_color)
		cur.X += dx; cur.Y += dy
	}
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



