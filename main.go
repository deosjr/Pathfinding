package main

import (
    "github.com/deosjr/Pathfinding/maps"
    "github.com/deosjr/Pathfinding/path"
)

func main() {
	grid := make([][]float64, 1000)
	for i := 0; i < 1000; i++ {
		grid[i] = make([]float64, 1000)
	}
	m := maps.NewGridMap(grid).WithPerlinNoise().SetWaterHeight(0.05)
	start, _ := m.Point(maps.NewPoint2D(0, 0))
	goal, _ := m.Point(maps.NewPoint2D(999, 999))
	route, _ := path.FindRoute(m, start, goal)
	m.Print(route)
}
