package main

import (
	"image/color"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/font/opentype"
)

// ==============================
// CUSTOM SETTINGS
// ==============================

const (
	TextMessage = "i like you"
	TextSize    = 24.0

	SideDuration = 10.0
	TrailSpacing = 30.0
	MaxTrail     = 55
)

var (
	BackgroundColor = color.RGBA{20, 20, 28, 255}
	TextColor       = color.RGBA{255, 190, 220, 255}
	TrailColor      = color.RGBA{255, 100, 170, 255}

	textFace font.Face
)

type Point struct {
	x, y float64
}

type Game struct {
	t            float64
	trail        []Point
	lastX        float64
	lastY        float64
	hasLast      bool
	screenWidth  int
	screenHeight int
}

func (g *Game) Update() error {
	// One complete heart loop.
	g.t += 1.0 / (60.0 * SideDuration)
	if g.t >= 1 {
		g.t -= 1
	}

	px, py := position(g.t, g.screenWidth, g.screenHeight)

	distance := math.Hypot(px-g.lastX, py-g.lastY)

	if !g.hasLast || distance >= TrailSpacing {
		g.trail = append(g.trail, Point{px, py})
		g.lastX, g.lastY = px, py
		g.hasLast = true
	}

	// Retain only the most recent trail copies.
	if len(g.trail) > MaxTrail {
		g.trail = g.trail[len(g.trail)-MaxTrail:]
	}

	return nil
}

func position(t float64, w, h int) (float64, float64) {
	angle := 2 * math.Pi * t

	hx := 16 * math.Pow(math.Sin(angle), 3)
	hy := 13*math.Cos(angle) -
		5*math.Cos(2*angle) -
		2*math.Cos(3*angle) -
		math.Cos(4*angle)

	scale := math.Min(float64(w), float64(h)) / 45.0
	x := float64(w)/2 + hx*scale
	y := float64(h)/2 - hy*scale

	return x, y
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(BackgroundColor)

	count := len(g.trail)

	for i, p := range g.trail {
		// 0 = oldest trail copy, 1 = newest trail copy.
		progress := 0.0
		if count > 1 {
			progress = float64(i) / float64(count-1)
		}

		// Blend from the background color to the trail color.
		r := uint8(float64(BackgroundColor.R) +
			progress*(float64(TrailColor.R)-float64(BackgroundColor.R)))
		gc := uint8(float64(BackgroundColor.G) +
			progress*(float64(TrailColor.G)-float64(BackgroundColor.G)))
		b := uint8(float64(BackgroundColor.B) +
			progress*(float64(TrailColor.B)-float64(BackgroundColor.B)))

		a := uint8(float64(TrailColor.A) * progress)
		c := color.RGBA{R: r, G: gc, B: b, A: a}

		drawText(screen, p.x, p.y, c)
	}

	px, py := position(g.t, g.screenWidth, g.screenHeight)
	drawText(screen, px, py, TextColor)
}

func drawText(screen *ebiten.Image, x, y float64, c color.Color) {
	bounds := font.MeasureString(textFace, TextMessage)
	text.Draw(
		screen,
		TextMessage,
		textFace,
		int(x)-bounds.Round()/2,
		int(y),
		c,
	)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	if g.screenWidth != outsideWidth || g.screenHeight != outsideHeight {
		g.trail = nil
		g.hasLast = false
	}
	g.screenWidth = outsideWidth
	g.screenHeight = outsideHeight
	return outsideWidth, outsideHeight
}

func init() {
	tt, _ := opentype.Parse(goregular.TTF)
	textFace, _ = opentype.NewFace(tt, &opentype.FaceOptions{
		Size: TextSize,
		DPI:  72,
	})
}

func main() {
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowSize(1000, 800)
	ebiten.SetWindowTitle("i like you")

	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
