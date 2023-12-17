package path

import (
	"container/heap"
	"fmt"
	"math"
)

type Map[Node comparable] interface {
	// get neighbours for a node
	Neighbours(n Node) []Node

	// cost of the path from start to node n
	// G() only returns the cost of moving from n to
	// one of its neighbours
	G(n, neighbour Node) float64

	// heuristic cost estimate function:
	// cost of cheapest path from n to goal
	H(n, goal Node) float64
}

func FindRoute[Node comparable](m Map[Node], start, goal Node) ([]Node, error) {
    return FindRouteWithGoalFunc(m, start, goal, func(c, g Node) bool {
        return c == g
    })
}

func FindRouteWithGoalFunc[Node comparable](m Map[Node], start, goal Node, goalFunc func(c, g Node) bool) ([]Node, error) {
	openSet := map[Node]bool{
		start: true,
	}
	// closedSet := map[Node]bool{}

	cameFrom := map[Node]Node{}
    var prevToGoal Node

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
		current := item.node.(Node)

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
		    if goalFunc(n, goal) {
			    goalScore = gScore[current]
                prevToGoal = n
			    // return reconstructPath(cameFrom, current), nil
		    }
			cameFrom[n] = current
			gScore[n] = tentativeGscore
			fScore[n] = f
		}
	}

	if goalScore == math.MaxInt64 {
		return nil, fmt.Errorf("No path found")
	}
	return reconstructPath(cameFrom, prevToGoal, goal), nil
}

func reconstructPath[Node comparable](m map[Node]Node, prevToGoal, goal Node) []Node {
    current := prevToGoal
	path := []Node{goal, prevToGoal}
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
	node   any
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
