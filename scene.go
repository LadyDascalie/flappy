package main

import (
	"fmt"

	"sync"

	"github.com/veandco/go-sdl2/sdl"
	img "github.com/veandco/go-sdl2/sdl_image"
)

type scene struct {
	r *sdl.Renderer

	bg *sdl.Texture

	bird bird
}

func newScene(r *sdl.Renderer, gravity float64) (s *scene, err error) {
	s = &scene{r: r}

	s.bg, err = img.LoadTexture(r, "res/imgs/background.png")
	if err != nil {
		return nil, fmt.Errorf("could not load background: %v", err)
	}

	s.bird = bird{
		x: 20, y: windowHeight / 2,
		gravity: gravity,
		w:       50, h: 43,
	}

	for i := 1; i <= 4; i++ {
		t, err := img.LoadTexture(r, fmt.Sprintf("res/imgs/bird_frame_%d.png", i))
		if err != nil {
			return nil, fmt.Errorf("could not load bird: %v", err)
		}
		s.bird.frames = append(s.bird.frames, t)
	}

	return s, nil
}

func (s *scene) run(fps float64) {
	for !s.bird.dead {
		s.bird.update()
		s.drawScene()
		sdl.Delay(uint32(1000 / fps))
	}
}

type bird struct {
	x, y    int32
	w, h    int32
	speed   float64
	gravity float64
	dead    bool
	frames  []*sdl.Texture
	frame   int
	mu      sync.Mutex
}

func (b *bird) update() {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.y += int32(b.speed)
	b.speed += b.gravity
	if b.y > windowHeight {
		fmt.Println("bird is dead")
		b.dead = true
	}
	b.frame = (b.frame + 1) % len(b.frames)
}

func (b *bird) jump() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.speed = -10
}

func (s *scene) drawScene() {
	s.r.Clear()

	s.r.Copy(s.bg, nil, &sdl.Rect{
		X: 0, Y: 0, W: windowWidth, H: windowHeight,
	})

	if !s.bird.dead {
		s.r.Copy(s.bird.frames[s.bird.frame], nil, &sdl.Rect{
			X: s.bird.x, Y: s.bird.y, W: s.bird.w, H: s.bird.h,
		})
	} else {
		drawText(s.r, "YOU DIED!", &sdl.Rect{
			X: 100, Y: windowHeight / 4, W: windowWidth - 200, H: windowHeight / 2,
		}, sdl.Color{R: 255})
	}

	s.r.Present()
}
