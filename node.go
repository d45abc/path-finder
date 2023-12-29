package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Node struct {
	x, y      float32
	nextNodes map[*Node]bool
}

func (n Node) drawNode(screen *ebiten.Image, radius float32, clr color.Color, op *ebiten.DrawImageOptions) {
	x, y := op.GeoM.Apply(float64(n.x), float64(n.y))
	vector.DrawFilledCircle(screen, float32(x), float32(y), radius, clr, true)

}

func (n Node) drawLinks(screen *ebiten.Image, width float32, clr color.Color, op *ebiten.DrawImageOptions) {
	n1x, n1y := op.GeoM.Apply(float64(n.x), float64(n.y))
	for next := range n.nextNodes {
		n2x, n2y := op.GeoM.Apply(float64(next.x), float64(next.y))
		vector.StrokeLine(screen, float32(n1x), float32(n1y), float32(n2x), float32(n2y), 3, color.RGBA{0, 255, 0, 255}, true)
	}
}

func (n Node) drawLinkTo(to *Node, screen *ebiten.Image, width float32, clr color.Color, op *ebiten.DrawImageOptions) {
	n1x, n1y := op.GeoM.Apply(float64(n.x), float64(n.y))
	n2x, n2y := op.GeoM.Apply(float64(to.x), float64(to.y))
	vector.StrokeLine(screen, float32(n1x), float32(n1y), float32(n2x), float32(n2y), 3, clr, true)
}
