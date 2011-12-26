
package main

import (
	"fmt"
)

const (	// TODO find a unified way to write this with stuff in edge_detectors
	LINE_EXPANSION = 1.3
	PROPORTION_KEEP = 0.7
	NUM_LINES = 20	// 9 cells in each dim, 10 lines in each dim X 2 dims
	MAX_ITER = 10
)

func main() {

	// read the image
	img := nil

	// initialize a whole bunch of lines randomly
	bounds := NewFloat64Rectangle(img.Bounds())
	lines := [NUM_LINES]Line
	for i,_ := range lines {
		lines[i] = RandomPointBetween(bounds.Min(), bounds.Max())
	}

	for iter := 0; iter < MAX_ITER; iter++ {

		// optimize positions of lines

		// take top 80% of the lines (by potential)
		
		// for remaining 20% randomly place lines on the grid
		//	- maybe utilize infomation about location of top 80%
		//	- maybe give these lines one round of optimization so they can compete	

		// increase each line length
		for _,l := range lines {
			l.ScaleLength(LINE_EXPANSION)
			l.ProjectInto(bounds)
	}

}

func OptimizePotential(lines []Lines, img image.Image) {
	// TODO maybe pass in a potential function
	// TODO this should be merged with Proposal() and Potenital() code in edge_detectors
}

