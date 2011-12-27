
package main

import (
	"os"
	"image"
	"image/png"
	"image/draw"
	"fmt"
	"rand"
	"math"
)

func max(a, b int) int {
	if a > b { return a }
	return b
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

func DarknessAt(img image.Image, x, y int) float64 {
	r, g, b, _ := img.At(x, y).RGBA()
	lum := 0.21 * float64(r) + 0.71 * float64(g) + 0.07 * float64(b)
	// TODO make this more flexible
	//return (255.0 - lum) / 255.0		// 8 bit
	return (65535.0 - lum) / 65535.0	// 16 bit
}

// makes a mutable copy
func CopyImage(img image.Image) (cpy draw.Image) {
	b := img.Bounds()
	cpy = image.NewRGBA(b.Dx(), b.Dy())
	for x := b.Min.X; x < b.Max.X; x++ {
		for y := b.Min.Y; y < b.Max.Y; y++ {
			cpy.Set(x, y, img.At(x, y))
		}
	}
	return cpy
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


