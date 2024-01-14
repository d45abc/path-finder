package main

import (
	"fmt"
	"image/color"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Game struct {
	originX, originY float64
	zoom, rotation   int
	nodes            map[*Node]bool
	op               *ebiten.DrawImageOptions
	hovered          *Node
	start, end       *Node
	path             []*Node
	foundPath        bool
}

func (g *Game) Update() error {
	g.updateOnCameraMove()
	g.updateDrawOptions()
	g.updateOnCursorMove()
	g.updateOnLeftClick()
	if g.end != nil && g.start != nil && !g.foundPath {
		g.path, g.foundPath = findPath(g.nodes, g.start, g.end)
	}
	return nil
}

func (g *Game) updateOnCursorMove() {
	if g.start != nil && g.end != nil {
		return
	}
	x, y := ebiten.CursorPosition()
	var minDistance float64 = 5
	for n := range g.nodes {
		nx, ny := g.op.GeoM.Apply(float64(n.x), float64(n.y))
		distance := math.Sqrt(math.Pow(nx-float64(x), 2) + math.Pow(ny-float64(y), 2))
		if distance < minDistance {
			g.hovered = n
		}
	}
}

func (g *Game) updateOnLeftClick() {
	if !inpututil.IsMouseButtonJustPressed(ebiten.MouseButton0) {
		return
	}
	if g.hovered != nil {
		if g.start == nil {
			g.start = g.hovered
		} else if g.end == nil {
			g.end = g.hovered
		}
		g.hovered = nil
	}
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
		g.hovered = nil
		g.end = nil
		g.start = nil
		g.foundPath = false
	}
	if ebiten.IsKeyPressed(ebiten.KeyF11) {
		ebiten.SetFullscreen(!ebiten.IsFullscreen())
	}
}

func (g *Game) updateDrawOptions() {
	screenWidth, screenHeight := ebiten.WindowSize()
	g.op.GeoM.Reset()
	g.op.GeoM.Translate(float64(g.originX), float64(g.originY))
	g.op.GeoM.Translate(-float64(screenWidth)/2, -float64(screenHeight)/2) // Translate center of the image to the top-left corner
	g.op.GeoM.Rotate(float64(g.rotation) * math.Pi / 180)
	scale := math.Pow(1.01, float64(g.zoom))
	g.op.GeoM.Scale(
		scale,
		scale,
	)
	g.op.GeoM.Translate(float64(screenWidth)/2, float64(screenHeight)/2)
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.White)
	g.drawGraph(screen)
	g.drawInfo(screen)
	g.drawPoints(screen)
	g.drawPath(screen)
}

func (g *Game) drawPath(screen *ebiten.Image) {
	if g.foundPath {
		for i := 0; i < len(g.path)-1; i++ {
			g.path[i].drawLinkTo(g.path[i+1], screen, 3, color.RGBA{255, 0, 0, 255}, g.op)
		}
	}
}

func (g *Game) drawPoints(screen *ebiten.Image) {
	if g.hovered != nil {
		g.hovered.drawNode(screen, 5, color.RGBA{255, 0, 0, 255}, g.op)
	}
	if g.start != nil {
		g.start.drawNode(screen, 5, color.RGBA{0, 128, 128, 255}, g.op)
	}
	if g.end != nil {
		g.end.drawNode(screen, 5, color.RGBA{255, 128, 128, 255}, g.op)
	}
}

func (g *Game) drawInfo(screen *ebiten.Image) {
	info := fmt.Sprintf("FPS: %2.f \nOrigin: (%.2f;%.2f) \nZoom: %v \nRotation: %v\nFound path: %v",
		ebiten.ActualFPS(),
		g.originX,
		g.originY,
		g.zoom,
		g.rotation,
		g.foundPath,
	)
	ebitenutil.DebugPrint(screen, info)
}

func (g *Game) drawGraph(screen *ebiten.Image) {
	screenWidth, screenHeight := ebiten.WindowSize()
	for n := range g.nodes {
		n1x, n1y := g.op.GeoM.Apply(float64(n.x), float64(n.y))
		if n1x > 0 && n1x < float64(screenWidth) && n1y > 0 && n1y < float64(screenHeight) {
			width := math.Pow(1.01, float64(g.zoom))
			n.drawLinks(screen, float32(width), color.RGBA{0, 255, 0, 255}, g.op)
		}
	}
	fmt.Println(screenWidth, screenHeight)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

func main() {
	g := &Game{
		op:    &ebiten.DrawImageOptions{},
		nodes: make(map[*Node]bool),
		path:  []*Node{},
	}
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	g.loadMap("example.geojson")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
