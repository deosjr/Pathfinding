package main

import (
	"container/heap"
	"fmt"
	"math"
)

type Point struct {
	x int
	y int
	z float64
}

type Point2D struct {
	x int
	y int
}

func FindRoute(m Map, start, goal Point) ([]Point, error) {
	openSet := map[Point]bool{
		start: true,
	}
	closedSet := map[Point]bool{}

	cameFrom := map[Point]Point{}

	gScore := map[Point]float64{}
	gScore[start] = 0

	fScore := map[Point]float64{}
	fScore[start] = h(start, goal)

	pq := priorityQueue{
		&pqItem{
			point:  start,
			fScore: fScore[start],
			index:  0,
		},
	}
	heap.Init(&pq)

	for pq.Len() != 0 {
		item := heap.Pop(&pq).(*pqItem)
		if item.fScore == math.MaxInt32 {
			break
		}
		current := item.point
		if current == goal {
			return reconstructPath(cameFrom, current), nil
		}

		// note: nodes can never be revisited atm
		closedSet[current] = true

		for _, neighbor := range m.Neighbors(current) {
			if closedSet[neighbor] {
				continue
			}

			gCurrent := getScore(gScore, current)
			gNeighbor := getScore(gScore, neighbor)
			//TODO: using maxfloat64 as 'notfound' causes overflow here
			// is there a better way to indicate actual infinity?
			tentativeGscore := gCurrent + g(m, current, neighbor)
			if tentativeGscore < gNeighbor {
				cameFrom[neighbor] = current
				gScore[neighbor] = tentativeGscore
				fScore[neighbor] = tentativeGscore + h(neighbor, goal)
			}

			if !openSet[neighbor] {
				openSet[neighbor] = true
				item := &pqItem{
					point:  neighbor,
					fScore: getScore(fScore, neighbor),
				}
				heap.Push(&pq, item)
			}
		}
	}

	return nil, fmt.Errorf("No path found")
}

func getScore(m map[Point]float64, p Point) float64 {
	x, ok := m[p]
	if !ok {
		return math.MaxInt32
	}
	return x
}

// TODO: cost function
func g(m Map, p, q Point) float64 {
	cost := 0.0
	if m.Water(q) > 0 {
		return math.MaxInt32
	}
	slope := math.Abs(q.z-p.z) / h(p, q)
	cost += 1000 * slope
	return cost
}

// heuristic cost estimate function
func h(p, q Point) float64 {
	return euclidian2d(p, q)
}

//euclidian distance in 2d
// geodesic, 'as the crow flies'
func euclidian2d(p, q Point) float64 {
	dx := float64(q.x - p.x)
	dy := float64(q.y - p.y)
	return math.Sqrt(dx*dx + dy*dy)
}

func reconstructPath(m map[Point]Point, current Point) []Point {
	path := []Point{current}
	for {
		prev, ok := m[current]
		if !ok {
			break
		}
		current = prev
		path = append(path, current)
	}
	return path
}

type pqItem struct {
	point  Point
	fScore float64
	index  int
}

type priorityQueue []*pqItem

func (pq priorityQueue) Len() int { return len(pq) }

func (pq priorityQueue) Less(i, j int) bool {
	return pq[i].fScore < pq[j].fScore
}

func (pq priorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *priorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*pqItem)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *priorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	item.index = -1
	*pq = old[0 : n-1]
	return item
}
