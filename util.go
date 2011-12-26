
package main

import (
	"os"
	"image"
	"image/png"
	"fmt"
	"rand"
	"math"
)

// TODO refactor all intstances of Float64Point in the file to two Float64Points

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

type Float64Rectangle struct {
	Min, Max Float64Point
}

func (r Float64Rectangle) Dx() (d float64) {
	return r.Max.X - r.Min.X
}

func (r Float64Rectangle) Dy() (d float64) {
	return r.Max.Y - r.Min.Y
}

func NewFloat64Rectangle(r image.Rectangle) (ret Float64Rectangle) {
	ret.Min.X = float64(r.Min.X)
	ret.Min.Y = float64(r.Min.Y)
	ret.Max.X = float64(r.Max.X)
	ret.Max.Y = float64(r.Max.Y)
	return ret
}

func max(a, b int) int {
	if a > b { return a }
	return b
}

func Midpoint(a, b Float64Point) (mid Float64Point) {
	mid = a
	mid.X += b.X
	mid.X /= 2.0
	mid.Y += b.Y
	mid.Y /= 2.0
	return mid
}

func WeightedChoice(weights []float64) int {
	s := 0.0
	for _,v := range weights {
		if v < 0.0 || v == math.NaN() || math.IsInf(v, 1) || math.IsInf(v, -1) {
			fmt.Printf("[WeightedChoice] illegal weight: %.2f\n", v)
			os.Exit(1)
		}
		s += v
	}
	if s == 0.0 {
		fmt.Printf("[WeightedChoice] all weights are 0!\n")
		os.Exit(1)
	}
	cutoff := rand.Float64() * s
	s = 0.0
	for i,v := range weights {
		if s >= cutoff { return i }
		s += v
	}
	fmt.Printf("[wtf] s = %.2f, cutoff = %.2f, weights = %s\n", s, cutoff, weights)
	return -1
}

func RandomPointBetween(lo, hi Float64Point) Float64Point {
	x := float64(hi.X) - rand.Float64() * float64(hi.X - lo.X)
	y := float64(hi.Y) - rand.Float64() * float64(hi.Y - lo.Y)
	return Float64Point{x, y}
}

func DarknessAt(img image.Image, x, y int) float64 {
	r, g, b, _ := img.At(x, y).RGBA()
	lum := 0.21 * float64(r) + 0.71 * float64(g) + 0.07 * float64(b)
	// TODO make this more flexible
	//return (255.0 - lum) / 255.0		// 8 bit
	return (65535.0 - lum) / 65535.0	// 16 bit
}

// converts an image to a mutable grayscale image
func Convert2Grayscale(input image.Image) *image.Gray16 {
	width := input.Bounds().Dx()
	height := input.Bounds().Dy()
	output := image.NewGray16(width, height)
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			r, g, b, _ := input.At(x, y).RGBA()
			lum := uint16(0.21 * float64(r) + 0.71 * float64(g) + 0.07 * float64(b))
			output.Set(x, y, image.Gray16Color{lum})
		}
	}
	return output
}

func SaveImage(img image.Image, outf string) {
	fmt.Printf("[SaveImage] saving to %s\n", outf)
	writer, err := os.OpenFile(outf, os.O_RDWR | os.O_CREATE, 0644)
	defer writer.Close()
	if err != nil {
		fmt.Printf("[SaveImage] could not open %s\n", outf)
		os.Exit(1)
	}
	err = png.Encode(writer, img)
	if err != nil {
		fmt.Printf("[SaveImage] problem saving to %s\n", outf)
		os.Exit(1)
	}
}

func OpenImage(img_name string) image.Image {
	file, err := os.Open(img_name)
	defer file.Close()
	if err != nil {
		fmt.Printf("could not find file: %s\n", img_name)
		os.Exit(1)
	}
	img, format, err := image.Decode(file)
	if err != nil {
		fmt.Printf("error while opening: %s\n", img_name)
		os.Exit(1)
	}
	fmt.Printf("loaded %s with format %s\n", img_name, format)
	return img
}


