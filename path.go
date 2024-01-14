package main

import (
	"container/heap"
	"math"
)

func findPathDijkstra(nodes map[*Node]bool, start, end *Node) ([]*Node, bool) {
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
			alt := dist[cur.node] + float32(distanceNode(cur.node, next))
			if distToNext, ok := dist[next]; alt < distToNext || !ok {
				dist[next] = alt
				prev[next] = cur.node
				q.tryToUpdate(next, alt)
			}
		}
	}
	return nil, false
}

func findPathAStar(nodes map[*Node]bool, start, end *Node) ([]*Node, bool) {
	prev := make(map[*Node]*Node)

	gScore := make(map[*Node]float32)
	gScore[start] = 0

	fScore := make(map[*Node]float32)
	fScore[start] = float32(distanceNode(start, end))

	openSet := make(PriorityQueue, 0)
	heap.Init(&openSet)
	heap.Push(&openSet, &Item{
		node: start,
		dist: fScore[start],
	})

	for len(openSet) > 0 {
		current := heap.Pop(&openSet).(*Item)
		if current.node == end {
			return reconstructPath(prev, start, end), true
		}
		for n := range current.node.nextNodes {
			tentative_gScore := gScore[current.node] + float32(distanceNode(current.node, n))

			if _, ok := gScore[n]; !ok {
				gScore[n] = float32(math.Inf(1))
			}
			if tentative_gScore < gScore[n] {
				prev[n] = current.node
				gScore[n] = tentative_gScore
				fScore[n] = tentative_gScore + float32(distanceNode(n, end))
				openSet.tryToUpdate(n, fScore[n])
			}
		}
	}
	return nil, false
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
