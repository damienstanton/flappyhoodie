package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/veandco/go-sdl2/sdl"
	ttf "github.com/veandco/go-sdl2/sdl_ttf"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("could not init SDL: %v", err)
		os.Exit(2)
	}
}

func run() error {
	err := sdl.Init(sdl.INIT_EVERYTHING)
	if err != nil {
		return fmt.Errorf("could not init SDL: %v", err)
	}
	defer sdl.Quit()

	err = ttf.Init()
	if err != nil {
		return fmt.Errorf("could not load font: %v", err)
	}
	defer ttf.Quit()

	w, r, err := sdl.CreateWindowAndRenderer(800, 600, sdl.WINDOW_SHOWN)
	if err != nil {
		return fmt.Errorf("could not create window: %v", err)
	}
	defer w.Destroy()

	if err = drawTitle(r); err != nil {
		return fmt.Errorf("could not draw title: %v", err)
	}

	// how long the title card appears for
	time.Sleep(1 * time.Second)

	s, err := newScene(r)
	if err != nil {
		return fmt.Errorf("could not create scene: %v", err)
	}
	defer s.destroy()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	select {
	case <-s.run(ctx, r):
		return fmt.Errorf("could not paint scene: %v", err)
	case <-time.After(5 * time.Second):
		return nil
	}
}

func drawTitle(r *sdl.Renderer) error {
	r.Clear()

	f, err := ttf.OpenFont("resources/fonts/FiraCode-Retina.ttf", 20)
	if err != nil {
		return fmt.Errorf("could not load FiraCode font: %v", err)
	}
	defer f.Close()

	color := sdl.Color{
		R: 255,
		G: 100,
		B: 0,
		A: 255,
	}

	s, err := f.RenderUTF8_Solid("Flappy Hoodie", color)
	if err != nil {
		return fmt.Errorf("could not render the title: %v", err)
	}
	defer s.Free()

	t, err := r.CreateTextureFromSurface(s)
	if err != nil {
		return fmt.Errorf("could not init texture: %v", err)
	}
	defer t.Destroy()

	if err := r.Copy(t, nil, nil); err != nil {
		return fmt.Errorf("could not copy texture: %v", err)
	}

	r.Present()
	return nil
}
