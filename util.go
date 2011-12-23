
package main

import (
	"os"
	"image"
	"image/png"
	"fmt"
	"rand"
)

type Float64Point struct {
	X, Y float64
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

func RandomPointBetween(lo, hi Float64Point) Float64Point {
	x := float64(hi.X) - rand.Float64() * float64(hi.X - lo.X)
	y := float64(hi.Y) - rand.Float64() * float64(hi.Y - lo.Y)
	return Float64Point{x, y}
}

func DarknessAt(img image.Image, x, y int) float64 {
	r, g, b, a := img.At(x, y).RGBA()
	lum := float64(a) * (0.21 * float64(r) + 0.71 * float64(g) + 0.07 * float64(b))
	return 255.0 - lum
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


