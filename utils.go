package main

import "math"

func distanceFloat[T float32 | float64](x1, y1, x2, y2 T) T {
	return T(math.Sqrt(math.Pow(float64(x1-x2), 2) + math.Pow(float64(y1-y2), 2)))
}

func distanceNode(n1, n2 *Node) float64 {
	return math.Sqrt(math.Pow(float64(n1.x-n2.x), 2) + math.Pow(float64(n1.y-n2.y), 2))
}

func invLerp(a, b, v float32) float32 {
	return (v - a) / (b - a)
}
