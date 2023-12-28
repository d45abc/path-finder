package main

import (
	"encoding/json"
	"math"
	"os"
)

type GeoJSON struct {
	Features []struct {
		Properties struct {
			Oneway string `json:"oneway"`
		} `json:"properties"`
		Geometry struct {
			Coordinates [][2]float32 `json:"coordinates"`
		} `json:"geometry"`
	} `json:"features"`
}

func (g *Game) loadMap(name string) {
	f, err := os.ReadFile(name)
	if err != nil {
		panic(err)
	}
	var data GeoJSON
	json.Unmarshal([]byte(f), &data)

	for _, f := range data.Features {
		coord := f.Geometry.Coordinates
		prev := g.addNode(coord[0][0], coord[0][1])
		for i := 1; i < len(coord); i++ {
			cur := g.addNode(coord[i][0], coord[i][1])
			g.addLink(prev, cur)
			if f.Properties.Oneway == "no" || f.Properties.Oneway == "" {
				g.addLink(cur, prev)
			}
			prev = cur
		}
	}

	var minLatitude float32 = math.MaxFloat32
	var maxLatitude float32 = -math.MaxFloat32
	var minLongtitude float32 = math.MaxFloat32
	var maxLongtitude float32 = -math.MaxFloat32

	for n := range g.nodes {
		minLatitude = min(minLatitude, n.y)
		maxLatitude = max(maxLatitude, n.y)
		minLongtitude = min(minLongtitude, n.x)
		maxLongtitude = max(maxLongtitude, n.x)
	}

	for n := range g.nodes {
		n.x = invLerp(minLongtitude, maxLongtitude, n.x) * float32(g.screenWidth)
		n.y = invLerp(maxLatitude, minLatitude, n.y) * float32(g.screenHeight)
	}

}

func (g *Game) addNode(x, y float32) *Node {
	for n := range g.nodes {
		if n.x == x && n.y == y {
			return n
		}
	}
	newNode := &Node{x, y, make(map[*Node]bool)}
	g.nodes[newNode] = true
	return newNode
}

func (g *Game) addLink(from, to *Node) {
	for n := range from.nextNodes {
		if n == to {
			return
		}
	}
	from.nextNodes[to] = true
}

func invLerp(a, b, v float32) float32 {
	return (v - a) / (b - a)
}
