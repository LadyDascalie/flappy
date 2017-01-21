package main

import (
	"fmt"
	"log"
	"math/rand"

	"github.com/veandco/go-sdl2/sdl"
	ttf "github.com/veandco/go-sdl2/sdl_ttf"
)

func trippyTitle(renderer *sdl.Renderer) {
	for i := 0; i < 30; i++ {
		renderer.Clear()
		if err := drawText(renderer,
			"Flappy Gopher",
			&sdl.Rect{X: 10, Y: windowHeight / 4, W: windowWidth - 20, H: windowHeight / 2},
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
