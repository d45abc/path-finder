package main

import (
	"container/heap"
)

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
