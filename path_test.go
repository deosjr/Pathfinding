package main

import (
	"reflect"
	"testing"
)

func TestFindRoute(t *testing.T) {
	for i, tt := range []struct {
		m     Map
		start Point
		goal  Point
		want  []Point
	}{
		{
			m: NewGridMap([][]int{
				{0, 0, 0},
				{0, 0, 0},
				{0, 0, 0},
			}),
			start: Point{0, 0, 0},
			goal:  Point{2, 2, 0},
			want: []Point{
				{2, 2, 0},
				{1, 1, 0},
				{0, 0, 0},
			},
		},
		{
			m: NewGridMap([][]int{
				{0, 1, 0},
				{0, 1, 0},
				{0, 1, 0},
			}),
			start: Point{0, 0, 0},
			goal:  Point{2, 0, 0},
			want: []Point{
				{2, 0, 0},
				{1, 0, 1},
				{0, 0, 0},
			},
		},
		{
			m: NewGridMap([][]int{
				{0, 9, 0},
				{0, 9, 0},
				{0, 0, 0},
			}),
			start: Point{0, 0, 0},
			goal:  Point{2, 0, 0},
			want: []Point{
				{2, 0, 0},
				{2, 1, 0},
				{1, 2, 0},
				{0, 1, 0},
				{0, 0, 0},
			},
		},
		{
			m: NewGridMap([][]int{
				{0, 0, 0},
				{0, -1, 0},
				{0, -1, 0},
			}).SetWaterHeight(-0.5),
			start: Point{0, 0, 0},
			goal:  Point{2, 2, 0},
			want: []Point{
				{2, 2, 0},
				{2, 1, 0},
				{1, 0, 0},
				{0, 0, 0},
			},
		},
	} {
		route, err := FindRoute(tt.m, tt.start, tt.goal)
		if err != nil {
			t.Errorf("%v", err.Error())
		}
		if !reflect.DeepEqual(route, tt.want) {
			t.Errorf("%d): got %v want %v", i, route, tt.want)
		}
	}
}
