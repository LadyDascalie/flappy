package main

import (
	"log"
	"math/rand"
	"time"

	"github.com/veandco/go-sdl2/sdl"
	img "github.com/veandco/go-sdl2/sdl_image"
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

	img.Init(img.INIT_PNG)
	defer img.Quit()

	window, renderer, err := sdl.CreateWindowAndRenderer(windowWidth, windowHeight, sdl.WINDOW_SHOWN)
	if err != nil {
		log.Fatalf("could not create window: %v", err)
	}
	defer window.Destroy()

	trippyTitle(renderer)

	scene, err := newScene(renderer, 1)
	if err != nil {
		log.Fatalf("could not create scene: %v", err)
	}
	go scene.run(50)

loop:
	for {
		switch event := sdl.WaitEvent().(type) {
		case *sdl.QuitEvent:
			break loop
		case *sdl.KeyUpEvent:
			scene.bird.jump()
		default:
			log.Printf("igoring event of type %T", event)
		}
	}
}
