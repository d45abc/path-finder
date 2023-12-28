package main

import (
	"fmt"
	"image/color"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Game struct {
	originX, originY          float64
	zoom, rotation            int
	nodes                     map[*Node]bool
	op                        *ebiten.DrawImageOptions
	screenWidth, screenHeight int
}

func (g *Game) Update() error {
	g.updateOnCameraMove()
	g.updateDrawOptions()
	return nil
}

func (g *Game) updateOnCameraMove() {
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		g.originY += 2 / math.Pow(1.01, float64(g.zoom))
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		g.originY -= 2 / math.Pow(1.01, float64(g.zoom))
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		g.originX += 2 / math.Pow(1.01, float64(g.zoom))
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		g.originX -= 2 / math.Pow(1.01, float64(g.zoom))
	}
	if ebiten.IsKeyPressed(ebiten.KeyQ) {
		g.rotation += 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyE) {
		g.rotation -= 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		g.zoom += 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		g.zoom -= 1
	}
	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		g.originX = 0
		g.originY = 0
		g.rotation = 0
		g.zoom = 0
	}
	if ebiten.IsKeyPressed(ebiten.KeyF11) {
		ebiten.SetFullscreen(!ebiten.IsFullscreen())
	}
}

func (g *Game) updateDrawOptions() {
	g.op.GeoM.Reset()
	g.op.GeoM.Translate(float64(g.originX), float64(g.originY))
	g.op.GeoM.Translate(-float64(g.screenWidth)/2, -float64(g.screenHeight)/2) // Translate center of the image to the top-left corner
	g.op.GeoM.Rotate(float64(g.rotation) * math.Pi / 180)
	scale := math.Pow(1.01, float64(g.zoom))
	g.op.GeoM.Scale(
		scale,
		scale,
	)
	g.op.GeoM.Translate(float64(g.screenWidth)/2, float64(g.screenHeight)/2)
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.White)
	g.drawGraph(screen)
	g.drawInfo(screen)
}

func (g *Game) drawInfo(screen *ebiten.Image) {
	info := fmt.Sprintf("FPS: %2.f \nOrigin: (%.2f;%.2f) \nZoom: %v \nRotation: %v",
		ebiten.ActualFPS(),
		g.originX,
		g.originY,
		g.zoom,
		g.rotation,
	)
	ebitenutil.DebugPrint(screen, info)
}

func (g *Game) drawGraph(screen *ebiten.Image) {
	for n := range g.nodes {
		n1x, n1y := g.op.GeoM.Apply(float64(n.x), float64(n.y))
		if n1x > 0 && n1x < float64(g.screenWidth) && n1y > 0 && n1y < float64(g.screenHeight) {
			width := math.Pow(1.01, float64(g.zoom))
			n.drawLinks(screen, float32(width), color.RGBA{0, 255, 0, 255}, g.op)
		}
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	g.screenWidth = outsideWidth
	g.screenHeight = outsideHeight
	return outsideWidth, outsideHeight
}

func main() {
	g := &Game{
		screenWidth:  600,
		screenHeight: 600,
		op:           &ebiten.DrawImageOptions{},
		nodes:        make(map[*Node]bool),
	}
	ebiten.SetWindowSize(g.screenWidth, g.screenHeight)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	g.loadMap("data.json")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}