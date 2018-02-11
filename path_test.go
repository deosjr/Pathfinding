package main

import (
	"reflect"
	"testing"
)

func TestFindRoute(t *testing.T) {
	for i, tt := range []struct {
		m     Map
		start Point2D
		goal  Point2D
		want  []Point
	}{
		{
			m: NewGridMap([][]float64{
				{0, 0, 0},
				{0, 0, 0},
				{0, 0, 0},
			}),
			start: Point2D{0, 0},
			goal:  Point2D{2, 2},
			want: []Point{
				{2, 2, 0},
				{1, 1, 0},
				{0, 0, 0},
			},
		},
		{
			m: NewGridMap([][]float64{
				{0, 1, 0},
				{0, 1, 0},
				{0, 1, 0},
			}),
			start: Point2D{0, 0},
			goal:  Point2D{2, 0},
			want: []Point{
				{2, 0, 0},
				{1, 0, 1},
				{0, 0, 0},
			},
		},
		{
			m: NewGridMap([][]float64{
				{0, 9, 0},
				{0, 9, 0},
				{0, 0, 0},
			}),
			start: Point2D{0, 0},
			goal:  Point2D{2, 0},
			want: []Point{
				{2, 0, 0},
				{2, 1, 0},
				{1, 2, 0},
				{0, 1, 0},
				{0, 0, 0},
			},
		},
		{
			m: NewGridMap([][]float64{
				{1, 1, 1},
				{1, 0, 1},
				{1, 0, 1},
			}).SetWaterHeight(0.5),
			start: Point2D{0, 0},
			goal:  Point2D{2, 2},
			want: []Point{
				{2, 2, 1},
				{2, 1, 1},
				{1, 0, 1},
				{0, 0, 1},
			},
		},
	} {

		grid := make([][]float64, 1000)
		for i := 0; i < 1000; i++ {
			grid[i] = make([]float64, 1000)
		}
		tt.m = NewGridMap(grid).WithPerlinNoise().SetWaterHeight(0.02)
		route, err := FindRoute(tt.m, *tt.m.Point(Point2D{0, 0}), *tt.m.Point(Point2D{999, 999}))
		tt.m.(GridMap).Print(route)
		break

		// route, err := FindRoute(tt.m, *tt.m.Point(tt.start), *tt.m.Point(tt.goal))
		if err != nil {
			t.Errorf("%v", err.Error())
		}
		if !reflect.DeepEqual(route, tt.want) {
			t.Errorf("%d): got %v want %v", i, route, tt.want)
		}
	}
}
