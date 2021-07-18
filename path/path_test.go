package path

import (
	"reflect"
	"testing"
)

//TODO fix import cycle
/*
func TestFindRoute(t *testing.T) {
	for i, tt := range []struct {
		m     GridMap
		start point2D
		goal  point2D
		want  []Node
	}{
		{
			m: NewGridMap([][]float64{
				{0, 0, 0},
				{0, 0, 0},
				{0, 0, 0},
			}),
			start: point2D{0, 0},
			goal:  point2D{2, 2},
			want: []Node{
				point{2, 2, 0},
				point{1, 1, 0},
				point{0, 0, 0},
			},
		},
		{
			m: NewGridMap([][]float64{
				{0, 1, 0},
				{0, 1, 0},
				{0, 1, 0},
			}),
			start: point2D{0, 0},
			goal:  point2D{2, 0},
			want: []Node{
				point{2, 0, 0},
				point{1, 0, 1},
				point{0, 0, 0},
			},
		},
		{
			m: NewGridMap([][]float64{
				{0, 9, 0},
				{0, 9, 0},
				{0, 0, 0},
			}),
			start: point2D{0, 0},
			goal:  point2D{2, 0},
			want: []Node{
				point{2, 0, 0},
				point{2, 1, 0},
				point{1, 2, 0},
				point{0, 1, 0},
				point{0, 0, 0},
			},
		},
		{
			m: NewGridMap([][]float64{
				{1, 1, 1},
				{1, 0, 1},
				{1, 0, 1},
			}).SetWaterHeight(0.5),
			start: point2D{0, 0},
			goal:  point2D{2, 2},
			want: []Node{
				point{2, 2, 1},
				point{2, 1, 1},
				point{1, 0, 1},
				point{0, 0, 1},
			},
		},
	} {
		start, _ := tt.m.point(tt.start)
		goal, _ := tt.m.point(tt.goal)
		route, err := FindRoute(tt.m, start, goal)
		if err != nil {
			t.Errorf("%v", err.Error())
		}
		if !reflect.DeepEqual(route, tt.want) {
			t.Errorf("%d): got %v want %v", i, route, tt.want)
		}
	}
}
*/
