package main

import (
	"context"
	"fmt"
	"time"

	"github.com/veandco/go-sdl2/sdl"
	img "github.com/veandco/go-sdl2/sdl_image"
)

type scene struct {
	time int
	bg   *sdl.Texture
	bird []*sdl.Texture
}

func newScene(r *sdl.Renderer) (*scene, error) {
	bg, err := img.LoadTexture(r, "resources/images/background.png")
	if err != nil {
		return nil, fmt.Errorf("could not load the background image: %v", err)
	}

	// bird animation
	var birds []*sdl.Texture
	for i := 1; i <= 4; i++ {
		p := fmt.Sprintf("resources/images/bird_frame_%d.png", i)
		bird, err := img.LoadTexture(r, p)
		if err != nil {
			return nil, fmt.Errorf("could not load bird frame: %v", err)
		}
		birds = append(birds, bird)
	}

	return &scene{
		bg:   bg,
		bird: birds,
	}, nil
}

// bird painting will be done in a different goroutine
func (s *scene) run(ctx context.Context, r *sdl.Renderer) chan error {
	errc := make(chan error)
	go func() {
		defer close(errc)
		for range time.Tick(10 * time.Millisecond) {
			select {
			case <-ctx.Done():
				return
			default:
				if err := s.paint(r); err != nil {
					errc <- err
				}
			}
		}
	}()
	return errc
}

func (s *scene) paint(r *sdl.Renderer) error {
	s.time++
	r.Clear()
	if err := r.Copy(s.bg, nil, nil); err != nil {
		return fmt.Errorf("could not copy background: %v", err)
	}

	rect := &sdl.Rect{X: 10, Y: 300 - 43/2, W: 50, H: 29}

	i := s.time / 10 % len(s.bird)
	if err := r.Copy(s.bird[i], nil, rect); err != nil {
		return fmt.Errorf("could not copy bird: %v", err)
	}

	r.Present()
	return nil
}

func (s *scene) destroy() {
	s.bg.Destroy()
}
