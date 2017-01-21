package main

import (
	"fmt"
	"math/rand"

	"sync"

	"github.com/veandco/go-sdl2/sdl"
	img "github.com/veandco/go-sdl2/sdl_image"
)

type scene struct {
	r *sdl.Renderer

	bg    *sdl.Texture
	bgt   int32
	pipes pipes
	bird  bird
}

func newScene(r *sdl.Renderer, speed int32, gravity float64) (s *scene, err error) {
	s = &scene{r: r}

	s.bg, err = img.LoadTexture(r, "res/imgs/background.png")
	if err != nil {
		return nil, fmt.Errorf("could not load background texture: %v", err)
	}

	s.pipes = pipes{speed: speed}

	for i := 0; i < 10; i++ {
		s.pipes.pipes = append(s.pipes.pipes, &pipe{
			pos: windowWidth + int32(rand.Intn(2*windowWidth)),
			w:   52,
			h:   int32(rand.Intn(windowHeight / 2)),
			up:  rand.Intn(10) > 4,
		})
	}

	s.pipes.tex, err = img.LoadTexture(r, "res/imgs/pipe.png")
	if err != nil {
		return nil, fmt.Errorf("could not load pipe texture: %v", err)

	}

	s.bird = bird{
		x: 20, y: windowHeight / 2,
		gravity: gravity,
		w:       50, h: 43,
	}
	for i := 1; i <= 4; i++ {
		t, err := img.LoadTexture(r, fmt.Sprintf("res/imgs/bird_frame_%d.png", i))
		if err != nil {
			return nil, fmt.Errorf("could not load bird texture: %v", err)
		}
		s.bird.frames = append(s.bird.frames, t)
	}

	return s, nil
}

func (s *scene) restart() {
	s.pipes.pipes = nil
	for i := 0; i < 10; i++ {
		s.pipes.pipes = append(s.pipes.pipes, &pipe{
			pos: windowWidth + int32(rand.Intn(2*windowWidth)),
			w:   52,
			h:   int32(rand.Intn(windowHeight / 2)),
			up:  rand.Intn(10) > 4,
		})
	}
	s.bird.y = windowHeight / 2
	s.bird.dead = false
	s.bird.speed = 0
}

func (s *scene) run(fps float64) {
	for {
		if !s.bird.dead {
			s.update()
			s.draw()
		}
		sdl.Delay(uint32(1000 / fps))
	}
}

func (s *scene) update() {
	s.pipes.update()
	s.bird.update()
	s.bgt = (s.bgt + 1) % 2000

	if s.pipes.hits(&s.bird) {
		s.bird.dead = true
	}
}

func (s *scene) draw() {
	s.r.Clear()

	s.r.Copy(s.bg, &sdl.Rect{X: s.bgt, Y: 0, W: windowWidth, H: windowHeight}, nil)

	s.bird.draw(s.r)
	s.pipes.draw(s.r)

	if s.bird.dead {
		drawText(s.r, "YOU DIED!", &sdl.Rect{
			X: 100, Y: windowHeight / 4, W: windowWidth - 200, H: windowHeight / 2,
		}, sdl.Color{R: 255})
	}

	s.r.Present()
}

type pipes struct {
	pipes []*pipe
	tex   *sdl.Texture
	speed int32
}

func (pp *pipes) update() {
	for _, p := range pp.pipes {
		p.pos -= pp.speed
		if p.pos < -p.w {
			p.pos = windowWidth + int32(rand.Intn(windowWidth))
			p.h = int32(rand.Intn(windowHeight))
			p.up = rand.Intn(10) > 4
		}
	}
}

func (pp *pipes) draw(r *sdl.Renderer) {
	for _, p := range pp.pipes {
		rect := &sdl.Rect{X: p.pos, Y: windowHeight - p.h, W: p.w, H: p.h}
		flip := sdl.FLIP_NONE
		if !p.up {
			rect.H = p.h
			rect.Y = 0
			flip = sdl.FLIP_VERTICAL
		}

		r.CopyEx(pp.tex, nil, rect, 0, nil, flip)
	}
}

func (pp *pipes) hits(b *bird) bool {
	for _, p := range pp.pipes {
		if p.hits(b) {
			return true
		}
	}
	return false
}

type pipe struct {
	pos  int32
	w, h int32
	up   bool
}

func (p *pipe) hits(b *bird) bool {
	if b.x+b.w <= p.pos || b.x >= p.pos+p.w {
		return false
	}
	if !p.up {
		return b.y <= p.h
	}
	return b.y+b.h >= windowHeight-p.h
}

type bird struct {
	x, y, w, h     int32
	speed, gravity float64
	dead           bool
	frames         []*sdl.Texture
	frame          int
	mu             sync.Mutex
}

func (b *bird) update() {
	b.mu.Lock()
	b.y += int32(b.speed)
	b.mu.Unlock()

	b.speed += b.gravity
	if b.y > windowHeight {
		fmt.Println("bird is dead")
		b.dead = true
	}
	b.frame = (b.frame + 1) % len(b.frames)
}

func (b *bird) draw(r *sdl.Renderer) {
	r.Copy(b.frames[b.frame], nil, &sdl.Rect{X: b.x, Y: b.y, W: b.w, H: b.h})
}

func (b *bird) jump() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.speed = -10
}
