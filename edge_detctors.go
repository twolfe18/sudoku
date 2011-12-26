
package main

import (
	"fmt"
	"math"
	"image"
	"rand"
	"os"
)

const (
	SudokuGridDimension = 9	// side of board (in squares, not lines)
)

type EdgeDetector struct {
	lines []Line

	// potential += exp(-sq_dist(point,pixel) / radius)
	default_line_radius float64

	// potential -= exp(-(angle(a,b) % 90.0) / orientation_sensitivity)
	orientation_sensitivity float64

	// how many proposals to make at each hill climbing iteration
	num_proposals uint

	// proposals are chosen with prob: l1_normalize(potentials ^ greedyness).
	// 0 is uniform choice, infinity is perfectly greedy
	greedyness float64

	proposal_variance float64
}

func (ed EdgeDetector) AlignTo(img image.Image) {
	// TODO i can just impelment each of these and see which is fastest (all derivative free)
	// option 1: draw K transforms, take the best point
	// option 2: draw K transforms, take the best point and do line search
	// option 3: draw K transforms, if best point isn't "good enough" then drak K more _smaller_ transforms
	bounds := NewFloat64Rectangle(img.Bounds())
	cur_ed := ed
	for iter := 0; iter < 15; iter++ {

		// propose some new edge detector positions
		minp := math.Inf(1)
		proposals := make([]EdgeDetector, ed.num_proposals)
		potentials := make([]float64, ed.num_proposals)
		for i := uint(0); i < cur_ed.num_proposals; i++ {
			proposals[i] = cur_ed.Proposal(bounds)
			potentials[i] = proposals[i].Potential(img)
			if potentials[i] < minp { minp = potentials[i] }
		}

		// make sure all potentials >= 0.0, calculate sum
		for i,_ := range potentials {
			c := potentials[i] - minp + 1	// smallest proposal will have potential = 1.0
			potentials[i] = math.Pow(c, ed.greedyness)
		}

		if len(potentials) != len(proposals) {
			fmt.Printf("[wtf] len(pot) = %d, len(pro) = %d\n", len(potentials), len(proposals))
			os.Exit(1)
		}
		i := WeightedChoice(potentials)
		cur_ed = proposals[i]
		fmt.Printf("[EdgeDetector.AlignTo] accepting pot=%.1f\tfrom [ ", potentials[i])
		for _,v := range potentials { fmt.Printf("%.1f ", v) }
		fmt.Printf("]\n")

		// test this on images to see how fast this should be decreased
		//cur_ed.proposal_variance *= 0.9

		// print out ED for debugging
		outf := fmt.Sprintf("/Users/travis/Dropbox/code/sudoku/img/debug.%d.png", iter)
		SaveImage(cur_ed.Draw(img), outf)
	}
}

func NewEdgeDetector(b Float64Rectangle) EdgeDetector {
	ed := new(EdgeDetector)
	ed.default_line_radius = 1.0
	ed.orientation_sensitivity = 3.0
	ed.num_proposals = 75
	ed.greedyness = 2.5
	ed.proposal_variance = 4.0	// in degrees

	// place some lines
	padding := 60.0	//2.0
	num_lines := 4	//SudokuGridDimension + 1
	dx := (b.Dx() - 2.0*padding) / float64(num_lines - 1)
	dy := (b.Dx() - 2.0*padding) / float64(num_lines - 1)
	x0 := b.Min.X + padding; xmax := b.Max.X - padding; x := x0
	y0 := b.Min.Y + padding; ymax := b.Max.Y - padding; y := y0
	//fmt.Printf("[NewEdgeDetector] x0 = %.2f, y0 = %.2f, xmax = %.2f, ymax = %.2f\n", x0, y0, xmax, ymax)
	for i := 0; i < num_lines; i++ {
		v := Line{Float64Point{x, y0}, Float64Point{x, ymax}, ed.default_line_radius}	// vertical
		h := Line{Float64Point{x0, y}, Float64Point{xmax, y}, ed.default_line_radius}	// horizontal
		ed.lines = append(ed.lines, v)
		ed.lines = append(ed.lines, h)
		x += dx; y += dy
		//fmt.Printf("[NewEdgeDetector] v = %s, h = %s\n", v.String(), h.String())
	}
	//fmt.Printf("[ned] ed.lines = %s\n", ed.lines)

	// random perturbation of "perfect"
	crappyness := 6.0
	ed.proposal_variance *= crappyness
	n := ed.Proposal(b)
	n.proposal_variance /= crappyness
	return n
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

	// rotations and shifts must be correlated
	independent_scale := 0.1
	mean_theta := (rand.Float64() * 2.0 - 1.0) * (math.Pi / 180.0 * ed.proposal_variance)
	mean_dx := (rand.Float64() * 2.0 - 1.0) * ed.proposal_variance
	mean_dy := (rand.Float64() * 2.0 - 1.0) * ed.proposal_variance

	for i, l := range ed.lines {

		nl := *new(Line)	// new line
		nl.radius = l.radius

		// first rotate the line
		theta := mean_theta + independent_scale * (rand.Float64() * 2.0 - 1.0) * (math.Pi / 180.0 * ed.proposal_variance)
		v := PointMinus(l.right, l.left)
		z := v.Rotate(theta)

		// scale back up to the correct length
		// new and old vecs share a midpoint, add/subtract half of the difference
		z.Scale(0.5)
		nl.left = PointMinus(Midpoint(l.left, l.right), z) 
		nl.right = PointPlus(Midpoint(l.left, l.right), z) 

		// now apply left-right and up-down shifts
		// TODO the indepented scale for dx dy shifts should be higher to allow for when
		// the original distance between lines is too great or small
		dx := mean_dx + independent_scale * (rand.Float64() * 2.0 - 1.0) * ed.proposal_variance	// left-right movement
		dy := mean_dy + independent_scale * (rand.Float64() * 2.0 - 1.0) * ed.proposal_variance	// up-down movement
		nl.left.X += dx
		nl.left.Y += dy
		nl.right.X += dx
		nl.right.Y += dy

		// now make sure it's in the bounds
		nl.ProjectInto(bounds)
		new_ed.lines[i] = nl
	}

	// scalings (shrinks and stretches) in x and y directions
	// TODO write variance struct that includes L/R, U/D shift amounts in (0,1)
	stretch := (rand.Float64() * 2.0 - 1.0) * ed.proposal_variance
	center := Float64Point{0.0, 0.0}	// find center of all lines, stretch to/from this point
	for _,l := range new_ed.lines {
		center := PointPlus(center, Midpoint(l))
	}
	center.Scale(1.0 / len(new_ed.lines))
	amount := 0.0	// TODO
	for _,l := range new_ed.lines {
		nmp := ShiftedMidpoint(l, center, amount)
		half := PointMinus(l.right, Midpoint(l))
		l.right = PointPlus(half, nmp)
		l.left = PointMinus(nmp, half)
	}

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
	var delta, dist float64
	add := 0.0
	b := img.Bounds()
	for x := b.Min.X; x < b.Max.X; x++ {
		for y := b.Min.Y; y < b.Max.Y; y++ {
			for _, line := range ed.lines {
				// TODO may need to play with this formula
				dist = line.SquaredDistance(float64(x), float64(y)) 
				delta = DarknessAt(img, x, y) * math.Exp(-dist / line.radius)
				p += delta; add += delta
				if math.IsInf(add, 1) {
					fmt.Printf("[Potential] hit inf!\n")
					os.Exit(1)
				}
			}
		}
	}
	p /= float64(len(ed.lines)); add /= float64(len(ed.lines))

	// orientation of the lines
	remove := 0.0
	num_pairs := 0	// man up: N * (N-1) / 2
	N := len(ed.lines)
	for i := 1; i < N; i++ {
		for j := 0; j < i; j++ {
			num_pairs += 1
			dist = ed.lines[i].Angle(ed.lines[j])
			delta = math.Exp(-math.Fmod(dist, 90.0)) * ed.orientation_sensitivity
			/*p -= delta;*/ remove += delta
		}
	}
	remove /= float64(num_pairs)
	p -= remove

	fmt.Printf("[EdgeDetector.Potential] potential = %.2f\t(+%.2f, -%.2f)\n", p, add, remove)
	return p
}

/**********************************************************************************************/

func main() {
	base := "/Users/travis/Dropbox/code/sudoku/img/"
	img := OpenImage(base + "clean_256_256.png")
	ed := NewEdgeDetector(NewFloat64Rectangle(img.Bounds()))

	// draw out ED right after creating it
	SaveImage(ed.Draw(img), base + "after_ed_init.png")

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
		ed.lines = append(ed.lines, Line{lo, hi, radius})
	}
	m_col_img := ed.Draw(m_gray_img)
	SaveImage(m_col_img, outf)
}



