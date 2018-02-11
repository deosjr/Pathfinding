package main

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
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type Map interface {
	// get x,y,z info for a 2d point
	Point(p Point2D) *Point
	// get neighbors in 2d grid
	Neighbors(p Point) []Point

	// INFORMATION FOR COST FUNCTIONS
	// get water height at point
	Water(p Point) float64
}

type GridMap struct {
	grid        [][]float64
	ySize       int
	xSize       int
	WaterHeight float64
}

func NewGridMap(grid [][]float64) GridMap {
	return GridMap{
		xSize:       len(grid[0]),
		ySize:       len(grid),
		grid:        grid,
		WaterHeight: math.MinInt64,
	}
}

func (gm GridMap) SetWaterHeight(x float64) GridMap {
	gm.WaterHeight = x
	return gm
}

func (gm GridMap) Point(p Point2D) *Point {
	if p.x < 0 || p.x > gm.xSize-1 || p.y < 0 || p.y > gm.ySize-1 {
		return nil
	}
	return &Point{p.x, p.y, gm.grid[p.y][p.x]}
}

func (gm GridMap) Neighbors(p Point) []Point {
	x, y := p.x, p.y
	points := []Point{}
	if p := gm.Point(Point2D{x - 1, y}); p != nil {
		points = append(points, *p)
	}
	if p := gm.Point(Point2D{x - 1, y - 1}); p != nil {
		points = append(points, *p)
	}
	if p := gm.Point(Point2D{x - 1, y + 1}); p != nil {
		points = append(points, *p)
	}
	if p := gm.Point(Point2D{x, y - 1}); p != nil {
		points = append(points, *p)
	}
	if p := gm.Point(Point2D{x, y + 1}); p != nil {
		points = append(points, *p)
	}
	if p := gm.Point(Point2D{x + 1, y - 1}); p != nil {
		points = append(points, *p)
	}
	if p := gm.Point(Point2D{x + 1, y}); p != nil {
		points = append(points, *p)
	}
	if p := gm.Point(Point2D{x + 1, y + 1}); p != nil {
		points = append(points, *p)
	}
	return points
}

// One water level over the entire map
func (gm GridMap) Water(p Point) float64 {
	if gm.WaterHeight == math.MinInt64 {
		return 0
	}
	z := float64(p.z)
	if z < gm.WaterHeight {
		return gm.WaterHeight - z
	}
	return 0
}

func (gm GridMap) WithPerlinNoise() GridMap {
	// alpha, beta, n iterations, random seed
	p := perlin.NewPerlin(2, 2, 3, rand.Int63())
	for y, row := range gm.grid {
		for x, _ := range row {
			nx := float64(x)/float64(gm.xSize) - 0.5
			ny := float64(y)/float64(gm.ySize) - 0.5
			noise := 0.5 * p.Noise2D(nx, ny)
			noise += 0.7 * p.Noise2D(2*nx, 2*ny)
			noise += 0.25 * p.Noise2D(4*nx, 4*ny)
			noise += 0.15 * p.Noise2D(8*nx, 8*ny)
			// normalize
			noise = noise / (0.5 + 0.7 + 0.25 + 0.15)
			// map from [-1,1] to [0,1]
			noise = (noise + 1) / 2
			noise = math.Pow(noise, 3.75)
			gm.grid[y][x] = noise
		}
	}
	return gm
}

func (gm GridMap) Print(route []Point) {
	m := image.NewRGBA(image.Rect(0, 0, gm.xSize, gm.ySize))
	for x := 0; x < gm.xSize; x++ {
		for y := 0; y < gm.ySize; y++ {
			pc := gm.grid[y][x]
			if pc > 1 {
				pc = 1
			}
			pcc := uint8(pc * 255)
			c := color.RGBA{pcc, pcc, pcc, 255}
			if pc < 0.02 {
				c = color.RGBA{0, 0, 255, 255}
			}
			m.Set(x, y, c)
		}
	}

	red := color.RGBA{255, 0, 0, 255}
	for _, v := range route {
		m.Set(v.x, v.y, red)
	}

	f, err := os.OpenFile("test.png", os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()
	png.Encode(f, m)
}
