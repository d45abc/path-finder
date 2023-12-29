package main

import (
	"container/heap"
	"math"
)

type Item struct {
	node  *Node
	dist  float32
	index int
}

type PriorityQueue []*Item

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].dist < pq[j].dist
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x any) {
	n := len(*pq)
	item := x.(*Item)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	item.index = -1
	*pq = old[0 : n-1]
	return item
}

func (pq *PriorityQueue) tryToUpdate(node *Node, dist float32) {
	for _, n := range *pq {
		if n.node == node {
			n.dist = dist
			heap.Fix(pq, n.index)
			return
		}
	}
	heap.Push(pq, &Item{
		node: node,
		dist: dist,
	})
}

func findPath(nodes map[*Node]bool, start, end *Node) ([]*Node, bool) {
	dist := make(map[*Node]float32)
	prev := make(map[*Node]*Node)
	q := make(PriorityQueue, 0)
	heap.Init(&q)
	heap.Push(&q, &Item{
		node: start,
		dist: 0,
	})

	dist[start] = 0

	for len(q) > 0 {
		cur := heap.Pop(&q).(*Item)
		if cur.node == end {
			return reconstructPath(prev, start, end), true
		}
		for next := range cur.node.nextNodes {
			alt := dist[cur.node] + d(cur.node, next)
			if distToNext, ok := dist[next]; alt < distToNext || !ok {
				dist[next] = alt
				prev[next] = cur.node
				q.tryToUpdate(next, alt)
			}
		}
	}
	return nil, false
}

func d(from, to *Node) float32 {
	return float32(math.Sqrt(math.Pow(float64(from.y-to.y), 2) + math.Pow(float64(from.x-to.x), 2)))
}

func reconstructPath(prev map[*Node]*Node, start, end *Node) []*Node {
	path := []*Node{}
	cur := end
	for cur != start {
		path = append(path, cur)
		cur = prev[cur]
	}
	path = append(path, start)
	return path
}
