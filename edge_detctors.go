
package sudoku

import (
	"fmt"
	"image"
	"image/png"
	"image/draw"
	"os"
	"math"
	"rand"
)

type EdgeDetector struct {
	lines []Line

	// potential += exp(-sq_dist(point,pixel) / radius)
	default_line_radius float64

	// potential -= exp(-(angle(a,b) % 90.0) * orientation_sensitivity)
	orientation_sensitivity float64

	// how many proposals to make at each hill climbing iteration
	num_proposals uint

	// proposals are chosen with prob: l1_normalize(potentials ^ greedyness).
	// 0 is uniform choice, infinity is perfectly greedy
	greedyness float64
}

func NewEdgeDetector() *EdgeDetector {
	ed := new(EdgeDetector)
	ed.default_line_radius = 10.0
	ed.orientation_sensitivity = 1.0
	ed.num_proposals = 10
	ed.greedyness = 1.0
	// TODO add some lines!
}

func (ed *EdgeDetector) AlignTo(img image.Image) {
	// TODO
	// put hillclimbing in here :)
// i can just impelment each of these and see which is fastest (all derivative free)
// option 1: draw K transforms, take the best point
// option 2: draw K transforms, take the best point and do line search
// option 3: draw K transforms, if best point isn't "good enough" then drak K more _smaller_ transforms
	fmt.Printf("[EdgeDetector.Align] need to implement\n")
	os.Exit(1)
}

func (ed EdgeDetector) Proposal(delta float64) *EdgeDetector {
	// TODO
	// return a new ED that has expected diff of delta
	fmt.Printf("[EdgeDetector.Proposal] need to implement\n")
	os.Exit(1)
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
				d := line.SquaredDistance(x, y) * img.DarknessAt(x, y)
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
	ed.Align(img)
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
	b := m_gray_img.Bounds()
	for i := 0; i < 500; i++ {
		mid := RandomPointBetween(b.Min, b.Max)
		lo := RandomPointBetween(b.Min, mid)
		hi := RandomPointBetween(mid, b.Max)
		radius := rand.Float64() * 5.0
		line := Line{lo, hi, radius}
		ed.lines = append(ed.lines, line)
	}
	m_col_img := ed.Draw(m_gray_img)
	SaveImage(m_col_img, outf)
}



