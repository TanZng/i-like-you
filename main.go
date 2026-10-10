package main

import (
	"image/color"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font/basicfont"
)

// ==============================
// CUSTOM SETTINGS
// ==============================

const (
	ScreenWidth  = 500
	ScreenHeight = 400

	TextMessage = "i like you"

	// Seconds per complete heart loop.
	SideDuration = 10.0

	// Distance between trail copies.
	TrailSpacing = 11.0

	// Number of trail copies.
	MaxTrail = 70

	// Trail opacity range: 0 = transparent, 255 = opaque.
	TrailStartOpacity = 100
	TrailEndOpacity   = 10

	// Heart position and size.
	PathX    = 250.0
	PathY    = 200.0
	PathSize = 9.0
)

var (
	BackgroundColor = color.RGBA{20, 20, 28, 255}
	TextColor       = color.RGBA{255, 190, 220, 255}
	TrailColor      = color.RGBA{255, 100, 170, 255}
)

type Point struct {
	x, y float64
}

type Game struct {
	t       float64
	trail   []Point
	lastX   float64
	lastY   float64
	hasLast bool
}

func (g *Game) Update() error {
	// One complete heart loop.
	g.t += 1.0 / (60.0 * SideDuration)
	if g.t >= 1 {
		g.t -= 1
	}

	px, py := position(g.t)

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

// Heart-shaped parametric path.
func position(t float64) (float64, float64) {
	angle := 2 * math.Pi * t

	// Classic mathematical heart curve.
	hx := 16 * math.Pow(math.Sin(angle), 3)
	hy := 13*math.Cos(angle) -
		5*math.Cos(2*angle) -
		2*math.Cos(3*angle) -
		math.Cos(4*angle)

	// Scale and center the heart.
	scale := PathSize
	x := PathX + hx*scale
	y := PathY - hy*scale

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

		c := color.RGBA{R: r, G: gc, B: b, A: 255}

		drawText(screen, p.x, p.y, c)
	}

	// Main text remains bright and readable.
	px, py := position(g.t)
	drawText(screen, px, py, TextColor)
}

func drawText(screen *ebiten.Image, x, y float64, c color.Color) {
	text.Draw(
		screen,
		TextMessage,
		basicfont.Face7x13,
		int(x)-30,
		int(y),
		c,
	)
}

func (g *Game) Layout(_, _ int) (int, int) {
	return ScreenWidth, ScreenHeight
}

func main() {
	ebiten.SetWindowSize(ScreenWidth*2, ScreenHeight*2)
	ebiten.SetWindowTitle("i like you")

	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
