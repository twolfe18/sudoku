
package main

import (
	"fmt"
	"math"
	"image"
	"rand"
)

type EdgeDetector struct {
	lines []RadialLine	// TODO make RadialLine (maybe)

	// potential += exp(-sq_dist(point,pixel) / radius)
	default_line_radius float64

	// potential -= exp(-(angle(a,b) % 90.0) * orientation_sensitivity)
	orientation_sensitivity float64

	// how many proposals to make at each hill climbing iteration
	num_proposals uint

	// proposals are chosen with prob: l1_normalize(potentials ^ greedyness).
	// 0 is uniform choice, infinity is perfectly greedy
	greedyness float64

	proposal_variance float64
}

func (ed EdgeDetector) AlignTo(img image.Image) {
	// TODO
	// i can just impelment each of these and see which is fastest (all derivative free)
	// option 1: draw K transforms, take the best point
	// option 2: draw K transforms, take the best point and do line search
	// option 3: draw K transforms, if best point isn't "good enough" then drak K more _smaller_ transforms
	bounds := NewFloat64Rectangle(img.Bounds())
	cur_ed := ed
	for iter := 0; iter < 10; iter++ {
		// propose some new edge detector positions
		s := 0.0
		proposals := make([]EdgeDetector, ed.num_proposals)
		potentials := make([]float64, ed.num_proposals)
		for i := uint(0); i < cur_ed.num_proposals; i++ {
			p := cur_ed.Proposal(bounds)
			pp := p.Potential(img)
			ppp := math.Pow(pp, cur_ed.greedyness)
			s += ppp
			proposals[i] = p
			potentials[i] = pp
		}

		// make a weighted random choice
		cutoff := s * rand.Float64()
		s = 0.0
		for i, pp := range potentials {
			s += math.Pow(pp, ed.greedyness)
			if cutoff < s {
				cur_ed = proposals[i]
				fmt.Printf("[EdgeDetector.AlignTo] accepting pot=%.1f\tfrom [ ")
				for _,v := range potentials { fmt.Printf("%.1f ", v) }
				fmt.Printf("]\n")
				break
			}
		}

		// test this on images to see how fast this should be decreased
		cur_ed.proposal_variance *= 0.9

		// print out ED for debugging
		outf := fmt.Sprintf("/Users/travis/Dropbox/code/sudoku/img/debug.%d.png", iter)
		SaveImage(cur_ed.Draw(img), outf)
	}
}

func NewEdgeDetector() EdgeDetector {
	ed := new(EdgeDetector)
	ed.default_line_radius = 10.0
	ed.orientation_sensitivity = 1.0
	ed.num_proposals = 10
	ed.greedyness = 1.0
	ed.proposal_variance = 1.0
	ed.lines = make([]RadialLine, ed.num_proposals)
	return *ed
}

func (ed EdgeDetector) CloneEdgeDetector() EdgeDetector {
	e := new(EdgeDetector)
	e.default_line_radius = ed.default_line_radius
	e.orientation_sensitivity = ed.orientation_sensitivity
	e.num_proposals = ed.num_proposals
	e.greedyness = ed.greedyness
	e.proposal_variance = ed.proposal_variance
	e.lines = ed.lines[:]
	return *e
}

func (ed EdgeDetector) Proposal(bounds Float64Rectangle) EdgeDetector {

	new_ed := ed.CloneEdgeDetector()

	for i, l := range ed.lines {

		new_midpoint := l.midpoint

		new_midpoint.X += (rand.Float64() * 2.0 - 1.0) * ed.proposal_variance
		new_midpoint.X = math.Fmax(new_midpoint.X, bounds.Min.X)
		new_midpoint.X = math.Fmin(new_midpoint.X, bounds.Max.X)

		new_midpoint.Y += (rand.Float64() * 2.0 - 1.0) * ed.proposal_variance
		new_midpoint.Y = math.Fmax(new_midpoint.Y, bounds.Min.Y)
		new_midpoint.Y = math.Fmin(new_midpoint.Y, bounds.Max.Y)

		new_rotation := l.rotation + (rand.Float64() * 2.0 - 1.0) * ed.proposal_variance

		// TODO make sure that this is inside bounds
		new_length := l.length

		new_radius := l.radius

		new_ed.lines[i] = RadialLine{new_midpoint, new_rotation, new_length, new_radius}
	}
	fmt.Printf("[EdgeDetector.Proposal] TODO make sure that new proposal is bounded!\n")
	return new_ed
}

func (ed EdgeDetector) Draw(img image.Image) image.Image {
	width := img.Bounds().Dx()
	height := img.Bounds().Dy()
	fmt.Printf("[EdgeDetector.draw] about to copy original image (%d, %d)...\n", width, height)
	output := image.NewRGBA(width, height)
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			output.Set(x, y, img.At(x, y))
		}
	}
	fmt.Printf("[EdgeDetector.draw] about to draw %d lines: ", len(ed.lines))
	for _, l := range ed.lines {
		l.Draw(output)
		fmt.Printf("*")
	}
	fmt.Printf("\n")
	return output
}

func (ed EdgeDetector) Potential(img image.Image) (p float64) {

	// put a "sparse prior" on random steps
		// steps should usually be mostly in one direction
		// compare to coordinate descent and no prior

	// does it make sense to have extra benefit for getting a cross at two intersecting lines?
		// this could get fooled on the numbers
		// probably not...

	// activation for each line and pixel
	add := 0.0
	b := img.Bounds()
	for x := b.Min.X; x < b.Max.X; x++ {
		for y := b.Min.Y; y < b.Max.Y; y++ {
			for _, line := range ed.lines {
				// TODO may need to play with this formula
				d := line.SquaredDistance(float64(x), float64(y)) * DarknessAt(img, x, y)
				fmt.Printf("[poten] d = %.2f\tsq_d = %.2f\tdarkness = %.2f\n", d, line.SquaredDistance(float64(x), float64(y)), DarknessAt(img, x, y))
				p += math.Exp(-1.0 * d / line.radius)
				add += math.Exp(-1.0 * d / line.radius)
			}
		}
	}

	// orientation of the lines
	remove := 0.0
	N := len(ed.lines)
	for i := 1; i < N; i++ {
		for j := 0; j < i; j++ {
			a := ed.lines[i].Angle(ed.lines[j])
			p -= math.Exp(-1.0 * math.Fmod(a, 90.0) * ed.orientation_sensitivity)
			remove += math.Exp(-1.0 * math.Fmod(a, 90.0) * ed.orientation_sensitivity)
		}
	}

	fmt.Printf("[EdgeDetector.Potential] potential = %.2f\t(+%.2f, -%.2f)\n", p, add, remove)
	return p
}

/**********************************************************************************************/

func main() {
	base := "/Users/travis/Dropbox/code/sudoku/img/"
	img := OpenImage(base + "clean_256_256.png")
	ed := NewEdgeDetector()
	ed.AlignTo(img)
	SaveImage(ed.Draw(img), base + "output.png")
}

func test_draw() {

	base := "/Users/travis/Dropbox/code/sudoku/img/"
	inf := base + "clean_256_256.png"
	outf := base + "output.png"
	img := OpenImage(inf)

	// convert to grayscale, make mutable
	m_gray_img := Convert2Grayscale(img)

	// draw a line on it
	ed := new(EdgeDetector)
	b := NewFloat64Rectangle(m_gray_img.Bounds())
	for i := 0; i < 500; i++ {
		mid := RandomPointBetween(b.Min, b.Max)
		lo := RandomPointBetween(b.Min, mid)
		hi := RandomPointBetween(mid, b.Max)
		radius := rand.Float64() * 5.0
		line := Point2Radial(PointLine{lo, hi, radius})
		ed.lines = append(ed.lines, line)
	}
	m_col_img := ed.Draw(m_gray_img)
	SaveImage(m_col_img, outf)
}



