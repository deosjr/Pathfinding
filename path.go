package main

import (
	"container/heap"
	"fmt"
	"math"
)

type Node interface{}

func FindRoute(m Map, start, goal Node) ([]Node, error) {
	openSet := map[Node]bool{
		start: true,
	}
	// closedSet := map[Node]bool{}

	cameFrom := map[Node]Node{}

	gScore := map[Node]float64{}
	gScore[start] = 0

	fScore := map[Node]float64{}
	fScore[start] = m.H(start, goal)

	pq := priorityQueue{
		&pqItem{
			node:   start,
			fScore: fScore[start],
			index:  0,
		},
	}
	heap.Init(&pq)

	// use goalScore when you want to be able to revisit goal
	// instead of returning the first path found
	goalScore := float64(math.MaxInt64)

	for pq.Len() != 0 {
		item := heap.Pop(&pq).(*pqItem)
		if item.fScore == math.MaxInt32 {
			break
		}
		current := item.node
		if current == goal {
			goalScore = gScore[current]
			// return reconstructPath(cameFrom, current), nil
		}

		// note: nodes can never be revisited if closedSet is used
		// closedSet[current] = true
		delete(openSet, current)

		for _, n := range m.Neighbours(current) {
			// if closedSet[n] {
			// 	continue
			// }
			tentativeGscore := gScore[current] + m.G(current, n)
			f := tentativeGscore + m.H(n, goal)

			v, ok := fScore[n]
			// only add to openSet if
			// - not in openSet yet
			// - has never been explored before or if it has, f < previous f score
			// - if goal has been found, only explore nodes for which f < current goal score
			if !openSet[n] && (!ok || f < v) && f < goalScore {
				openSet[n] = true
				item := &pqItem{
					node:   n,
					fScore: f,
				}
				heap.Push(&pq, item)
			} else if tentativeGscore >= gScore[n] {
				continue
			}
			cameFrom[n] = current
			gScore[n] = tentativeGscore
			fScore[n] = f
		}
	}

	if goalScore == math.MaxInt64 {
		return nil, fmt.Errorf("No path found")
	}
	return reconstructPath(cameFrom, goal), nil
}

func reconstructPath(m map[Node]Node, current Node) []Node {
	path := []Node{current}
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
	node   Node
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

func main() {
	grid := make([][]float64, 1000)
	for i := 0; i < 1000; i++ {
		grid[i] = make([]float64, 1000)
	}
	m := NewGridMap(grid).WithPerlinNoise().SetWaterHeight(0.02)
	start, _ := m.point(point2D{0, 0})
	goal, _ := m.point(point2D{999, 999})
	route, _ := FindRoute(m, start, goal)
	m.Print(route)
}
