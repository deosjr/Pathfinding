package main

type Map interface {
	// get x,y,z info for a 2d point
	Point(x, y int) *Point
	// get neighbors in 2d grid
	Neighbors(p Point) []Point

	// INFORMATION FOR COST FUNCTIONS
	// get water height at point
	Water(p Point) float64
}

type GridMap struct {
	grid        [][]int
	ySize       int
	xSize       int
	WaterHeight float64
}

func NewGridMap(grid [][]int) GridMap {
	return GridMap{
		xSize: len(grid[0]),
		ySize: len(grid),
		grid:  grid,
	}
}

func (gm GridMap) SetWaterHeight(x float64) GridMap {
	gm.WaterHeight = x
	return gm
}

func (gm GridMap) Point(x, y int) *Point {
	if x < 0 || x > gm.xSize-1 || y < 0 || y > gm.ySize-1 {
		return nil
	}
	return &Point{x, y, gm.grid[y][x]}
}

func (gm GridMap) Neighbors(p Point) []Point {
	x, y := p.x, p.y
	points := []Point{}
	if p := gm.Point(x-1, y); p != nil {
		points = append(points, *p)
	}
	if p := gm.Point(x-1, y-1); p != nil {
		points = append(points, *p)
	}
	if p := gm.Point(x-1, y+1); p != nil {
		points = append(points, *p)
	}
	if p := gm.Point(x, y-1); p != nil {
		points = append(points, *p)
	}
	if p := gm.Point(x, y+1); p != nil {
		points = append(points, *p)
	}
	if p := gm.Point(x+1, y-1); p != nil {
		points = append(points, *p)
	}
	if p := gm.Point(x+1, y); p != nil {
		points = append(points, *p)
	}
	if p := gm.Point(x+1, y+1); p != nil {
		points = append(points, *p)
	}
	return points
}

// One water level over the entire map
func (gm GridMap) Water(p Point) float64 {
	z := float64(p.z)
	if z < gm.WaterHeight {
		return gm.WaterHeight - z
	}
	return 0
}
