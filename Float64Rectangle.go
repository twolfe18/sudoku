
package main

import "image"

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
	ret.Min = NewFloat64Point(r.Min)
	ret.Max = NewFloat64Point(r.Max)
	return ret
}


