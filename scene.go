package main

import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
	img "github.com/veandco/go-sdl2/sdl_image"
)

type scene struct {
	r  *sdl.Renderer
	bg *sdl.Texture
}

func newScene(r *sdl.Renderer) (*scene, error) {
	bg, err := img.LoadTexture(r, "res/imgs/background.png")
	if err != nil {
		return nil, fmt.Errorf("could not load background: %v", err)
	}

	return &scene{r: r, bg: bg}, nil
}

func (s *scene) drawFrame() {
	s.r.Clear()
	s.r.Copy(s.bg, nil, &sdl.Rect{
		X: 0, Y: 0, W: windowWidth, H: windowHeight,
	})
	s.r.Present()
}
