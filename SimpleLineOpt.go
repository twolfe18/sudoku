
package main

import (
	"fmt"
	"image"
	"math"
)

const (	// TODO find a consistent way to write this with stuff in edge_detectors
	LINE_EXPANSION = 1.3
	PROPORTION_KEEP = 0.7
	NUM_LINES = 1	//20	// 9 cells in each dim, 10 lines in each dim X 2 dims
	MAX_ITER = 10
)

/*func later() {

	// TODO have lines *mask* out the dark parts that they are covering

	// TODO first try line optimization with one line

	// read the image
	img := nil

	// initialize a whole bunch of lines randomly
	bounds := NewFloat64Rectangle(img.Bounds())
	lines := nil//[NUM_LINES]Line
	for i,_ := range lines {
		lines[i] = RandomPointBetween(bounds.Min(), bounds.Max())
	}

	for iter := 0; iter < MAX_ITER; iter++ {

		// optimize positions of lines

		// take top 80% of the lines (by potential)
		// i think this way is better than somehow specifying how long
		// grid lines should be. this is less dependent on the picture (maybe)


		// for remaining 20% randomly place lines on the grid
		//	- maybe utilize infomation about location of top 80%
		//	- maybe give these lines one round of optimization so they can compete	

		// increase each line length
		for _,l := range lines {
			l.ScaleLength(LINE_EXPANSION)
			l.ProjectInto(bounds)
		}
	}

} */

func main() {
	base := "/Users/travis/Dropbox/code/sudoku/img/"
	img := OpenImage(base + "clean_256_256.png")
	SaveImage(img, "after_init.png")

	p := DefaultParams()
	lines := make([]Line, 0)
	for len(lines) < 20 {

		// randomly place a line on the board
		b := NewFloat64Rectangle(img.Bounds())
		left := RandomPointBetween(b.Min, b.Max)
		right := RandomPointBetween(b.Min, b.Max)
		i := len(lines)
		lines := append(lines, Line{left, right, 1.0})

		// TODO mask off current lines

		// see where it goes to
		for iter := 0; iter < 10; iter++ {
			newline := LocalOptimizePotential(lines[i], img, p)
			if lines[i].Equals(newline) {
				fmt.Printf("[main] converged at iter %d\n", iter)
				break
			} else {
				lines[i] = newline
				cpy := CopyImage(img)
				for _,l := range lines {
					l.Draw(cpy, image.RGBAColor{255, 0, 0, 255})
				}
				SaveImage(cpy, fmt.Sprintf("%sdebug.%d.%d.png", base, i, iter))
			}
		}
	}
}

type Params struct {
	lambda_dtheta, delta_dtheta, max_dtheta float64
	lambda_dx, delta_dx, max_dx float64
	lambda_dy, delta_dy, max_dy float64
}

func DefaultParams() (p Params) {
	p.lambda_dtheta = 1.0
	p.lambda_dx = 1.0
	p.lambda_dy = 1.0
	p.delta_dtheta = 1
	p.delta_dx = 1.0
	p.delta_dy = 1.0
	p.max_dtheta = 15.0
	p.max_dx = 10.0
	p.max_dy = 10.0
	return p
}

func LocalOptimizePotential(line Line, img image.Image, p Params) (bestline Line) {
	// TODO do some kind of branch and bound
	var newline Line
	bestpot := math.Inf(-1)
	fmt.Printf("[LocalOptimizePotential] about to go through %d *s: ", int(2.0 * p.max_dtheta / p.delta_dtheta))
	best_theta := 0.0
	for dtheta := -p.max_dtheta; dtheta <= p.max_dtheta; dtheta += p.delta_dtheta {
		fmt.Printf("*")
		for dx := -p.max_dx; dx <= p.max_dx; dx += p.delta_dx {
			for dy := -p.max_dy; dy <= p.max_dy; dy += p.delta_dy {
				newline = line
				newline.Rotate(dtheta)
				newline.Shift(dx, dy)
				p := LinePotential(newline, img) - p.lambda_dtheta * dtheta - p.lambda_dx * dx - p.lambda_dy * dy
				if p > bestpot {
					best_theta = dtheta
					bestpot = p
					bestline = newline
				}
			}
		}
	}
	fmt.Printf("\n")
	fmt.Printf("[LocalOptimizePotential] dtheta = %.2f opt %s \t to \t %s\n", best_theta, line, bestline)
	return bestline
}

func LinePotential(line Line, img image.Image) (pot float64) {
	for _,wp := range line.WeightedIterator() {
		darkness := DarknessAt(img, wp.P.X, wp.P.Y)
		if wp.W < 0.0 || wp.W > 1.0 {
			panic(fmt.Sprintf("[LinePotential] weight must be in [0,1]: %.2f", wp.W))
		}
		pot += darkness * wp.W
	}
	return pot
}



