
package main

import (
	"os"
	"image"
	"image/png"
	"fmt"
	"rand"
	"math"
)

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


