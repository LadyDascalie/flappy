package main

import (
	"fmt"
	"log"

	"sync"

	"github.com/veandco/go-sdl2/sdl"
	img "github.com/veandco/go-sdl2/sdl_image"
)

type scene struct {
	r *sdl.Renderer

	bg    *sdl.Texture
	pipes pipes
	bird  bird
}

func newScene(r *sdl.Renderer, speed int32, gravity float64) (s *scene, err error) {
	s = &scene{r: r}

	s.bg, err = img.LoadTexture(r, "res/imgs/background.png")
	if err != nil {
		return nil, fmt.Errorf("could not load background texture: %v", err)
	}

	s.pipes = pipes{
		speed: speed,
		pipe: pipe{
			pos: windowWidth,
			w:   52,
			h:   320,
		},
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

func (s *scene) run(fps float64) {
	for !s.bird.dead {
		s.update()
		s.draw()
		sdl.Delay(uint32(1000 / fps))
	}
}

func (s *scene) update() {
	s.pipes.update()
	s.bird.update()

	if s.pipes.hits(&s.bird) {
		s.bird.dead = true
	}
}

func (s *scene) draw() {
	s.r.Clear()

	s.r.Copy(s.bg, nil, &sdl.Rect{X: 0, Y: 0, W: windowWidth, H: windowHeight})

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
	pipe  pipe
	tex   *sdl.Texture
	speed int32
}

func (p *pipes) update() {
	p.pipe.pos -= p.speed
}

func (pp *pipes) draw(r *sdl.Renderer) {
	//for _, p := range pp.pipe {
	p := pp.pipe
	rect := &sdl.Rect{X: p.pos, Y: windowHeight - p.h, W: p.w, H: p.h}
	log.Printf("painting on %+v", rect)
	r.Copy(pp.tex, nil, rect)
	//}
}

func (p pipes) hits(b *bird) bool {
	return p.pipe.hits(b)
}

type pipe struct {
	pos  int32
	w, h int32
}

func (p pipe) hits(b *bird) bool {
	return b.y+b.h >= p.h && b.x+b.w > p.pos && b.x < p.pos+p.w
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
