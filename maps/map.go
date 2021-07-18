package maps

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"math/rand"
	"os"
	"time"

	perlin "github.com/aquilax/go-perlin"
    "github.com/deosjr/Pathfinding/path"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}


type GridMap struct {
	grid        [][]float64
	ySize       int
	xSize       int
	waterHeight float64
}

func NewGridMap(grid [][]float64) GridMap {
	return GridMap{
		xSize:       len(grid[0]),
		ySize:       len(grid),
		grid:        grid,
		waterHeight: math.MinInt64,
	}
}

type point struct {
	x, y int
	z    float64
}

type point2D struct {
	x, y int
}

func NewPoint2D(x, y int) point2D {
    return point2D{x, y}
}

func (gm GridMap) SetWaterHeight(x float64) GridMap {
	gm.waterHeight = x
	return gm
}

func (gm GridMap) Point(p point2D) (point, bool) {
	if p.x < 0 || p.x > gm.xSize-1 || p.y < 0 || p.y > gm.ySize-1 {
		return point{}, false
	}
	return point{p.x, p.y, gm.grid[p.y][p.x]}, true
}

func (gm GridMap) Neighbours(n path.Node) []path.Node {
	p := n.(point)
	x, y := p.x, p.y
	points := []path.Node{}
	points2d := []point2D{
		// cardinal directions
		{x - 1, y},
		{x, y - 1},
		{x, y + 1},
		{x + 1, y},

		// Mask M1
		{x - 1, y - 1},
		{x - 1, y + 1},
		{x + 1, y - 1},
		{x + 1, y + 1},

		// Mask M2
		{x + 1, y + 2},
		{x + 2, y + 1},
		{x + 2, y - 1},
		{x + 1, y - 2},
		{x - 1, y - 2},
		{x - 2, y - 1},
		{x - 2, y + 1},
		{x - 1, y + 2},

		// Mask M3
		{x + 1, y + 3},
		{x + 2, y + 3},
		{x + 3, y + 2},
		{x + 3, y + 1},
		{x + 3, y - 1},
		{x + 3, y - 2},
		{x + 2, y - 3},
		{x + 1, y - 3},
		{x - 1, y - 3},
		{x - 2, y - 3},
		{x - 3, y - 2},
		{x - 3, y - 1},
		{x - 3, y + 1},
		{x - 3, y + 2},
		{x - 2, y + 3},
		{x - 1, y + 3},
	}
	for _, p2d := range points2d {
		if p, ok := gm.Point(p2d); ok {
			points = append(points, p)
		}
	}
	return points
}

// TODO: cost function
func (m GridMap) G(pn, qn path.Node) float64 {
	p, q := pn.(point), qn.(point)
	cost := 0.0
	if m.Water(q) > 0 {
		return math.MaxInt32
	}
    cost += m.H(p, q)
	slope := math.Abs(q.z-p.z)
	cost += slope
	return cost
}

func (m GridMap) H(pn, qn path.Node) float64 {
	p, q := pn.(point), qn.(point)
	return euclidian2d(p, q)
}

// euclidian distance in 2d
// geodesic, 'as the crow flies'
func euclidian2d(p, q point) float64 {
	dx := float64(q.x - p.x)
	dy := float64(q.y - p.y)
	return math.Sqrt(dx*dx + dy*dy)
}

// INFORMATION FOR COST FUNCTIONS
// get water height at point
// One water level over the entire map
func (gm GridMap) Water(p point) float64 {
	if gm.waterHeight == math.MinInt64 {
		return 0
	}
	z := float64(p.z)
	if z < gm.waterHeight {
		return gm.waterHeight - z
	}
	return 0
}

func (gm GridMap) WithPerlinNoise() GridMap {
	// alpha, beta, n iterations, random seed
	p := perlin.NewPerlin(2, 4, 3, rand.Int63())
	for y, row := range gm.grid {
		for x, _ := range row {
			nx := float64(x)/float64(gm.xSize) - 0.5
			ny := float64(y)/float64(gm.ySize) - 0.5
			noise := 0.3 * p.Noise2D(nx, ny)
			noise += 0.8 * p.Noise2D(2*nx, 2*ny)
			noise += 0.25 * p.Noise2D(4*nx, 4*ny)
			noise += 0.15 * p.Noise2D(8*nx, 8*ny)
			// normalize
			noise = noise / (0.3 + 0.8 + 0.25 + 0.15)
			// map from [-1,1] to [0,1]
			noise = (noise + 1) / 2
			//noise = math.Pow(noise, 3.75)
			noise = math.Pow(noise, 3.5)
			gm.grid[y][x] = noise
		}
	}
	return gm
}

func (gm GridMap) Print(route []path.Node) {
	m := image.NewRGBA(image.Rect(0, 0, gm.xSize, gm.ySize))
	for x := 0; x < gm.xSize; x++ {
		for y := 0; y < gm.ySize; y++ {
			pc := gm.grid[y][x]
			if pc > 1 {
				pc = 1
			}
			pcc := uint8(pc * 255)
			c := color.RGBA{pcc, pcc, pcc, 255}
			if pc < 0.05 {
				c = color.RGBA{0, 0, 255, 255}
			}
			m.Set(x, y, c)
		}
	}

	red := color.RGBA{255, 0, 0, 255}
	for _, v := range route {
		p := v.(point)
		m.Set(p.x, p.y, red)
	}

	f, err := os.OpenFile("test.png", os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()
	png.Encode(f, m)
}

