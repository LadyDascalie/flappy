package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/veandco/go-sdl2/sdl"
	ttf "github.com/veandco/go-sdl2/sdl_ttf"
)

const (
	windowHeight = 600
	windowWidth  = 800
)

func main() {
	rand.Seed(time.Now().Unix())

	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		log.Fatalf("could not initialize sdl: %v", err)
	}
	defer sdl.Quit()

	if err := ttf.Init(); err != nil {
		log.Fatalf("could not initialize ttf: %v", err)
	}
	defer ttf.Quit()

	window, renderer, err := sdl.CreateWindowAndRenderer(windowWidth, windowHeight, sdl.WINDOW_SHOWN)
	if err != nil {
		log.Fatalf("could not create window: %v", err)
	}
	defer window.Destroy()

	trippyTitle(renderer)

loop:
	for {
		switch event := sdl.WaitEvent().(type) {
		case *sdl.QuitEvent:
			break loop
		default:
			log.Printf("igoring event of type %T", event)
		}
	}
}

func trippyTitle(renderer *sdl.Renderer) {
	for i := 0; i < 30; i++ {
		renderer.Clear()
		if err := drawText(renderer,
			"Flappy Gopher",
			&sdl.Rect{X: 0, Y: windowHeight / 4, W: windowWidth, H: windowHeight / 2},
			randomColor(),
		); err != nil {
			log.Fatal(err)
		}
		sdl.Delay(100)
	}
}

func randomColor() sdl.Color {
	rand := func() uint8 { return uint8(rand.Intn(256)) }
	return sdl.Color{R: rand(), G: rand(), B: rand(), A: 0}
}

func drawText(renderer *sdl.Renderer, title string, rect *sdl.Rect, color sdl.Color) error {
	path := "res/fonts/Flappy.ttf"
	font, err := ttf.OpenFont(path, 30)
	if err != nil {
		return fmt.Errorf("could not load font: %v", err)
	}

	surface, err := font.RenderUTF8_Solid(title, color)
	if err != nil {
		return fmt.Errorf("could not render title: %v", err)
	}
	defer surface.Free()

	tex, err := renderer.CreateTextureFromSurface(surface)
	if err != nil {
		return fmt.Errorf("could not create texture: %v", err)
	}

	err = renderer.Copy(tex, nil, rect)
	if err != nil {
		return fmt.Errorf("could not copy texture: %v", err)
	}

	renderer.Present()
	return nil
}
