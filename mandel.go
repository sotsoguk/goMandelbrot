package main

import (
	// "fmt"
	// "image"
	// "image/color"
	// "image/png"
	// "math/cmplx"
	// "os"
	// "time"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math/cmplx"
	"os"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

func main() {
	fmt.Println("GoMandel!")
	start := time.Now()
	var z_ll complex128 = complex(-2.0, -1.2)
	var z_ur complex128 = complex(1.0, 1.2)
	// var z_ll complex128 = complex(0, 0)
	// var z_ur complex128 = complex(0.5, 0.6)

	xMin := real(z_ll)
	xMax := real(z_ur)
	yMin := imag(z_ll)
	yMax := imag(z_ur)

	nx := 1200
	ny := int((float64(nx) * (yMax - yMin)) / (xMax - xMin))
	rect := image.Rect(0, 0, nx, ny)
	img := image.NewRGBA(rect)
	var window *sdl.Window
	var renderer *sdl.Renderer

	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	defer sdl.Quit()

	window, err := sdl.CreateWindow("test", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		int32(nx), int32(ny), sdl.WINDOW_SHOWN)
	if err != nil {
		fmt.Println(os.Stderr, "Failed to create window: %s\n", err)
		panic(err)
	}
	defer window.Destroy()

	renderer, err = sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		fmt.Println(os.Stderr, "Failed to create renderer: %s\n", err)
		panic(err)
	}

	defer renderer.Destroy()

	maxIter := 200
	dX := float64(xMax-xMin) / float64(nx)
	dY := float64(yMax-yMin) / float64(ny)

	var c complex128 = z_ll
	for y := 0; y < ny; y++ {
		for x := 0; x < nx; x++ {
			c = z_ll + complex(float64(x)*dX, 0) + complex(0, float64(y)*dY)
			z := c
			isInside := true
			numToInfinity := 0

			for n := 0; n < maxIter; n++ {
				z = z*z + c
				if cmplx.Abs(z) > 2 {
					isInside = false
					numToInfinity = n
					break
				}
			}
			if isInside {
				img.Set(x, y, color.Black)
				renderer.SetDrawColor(0, 0, 0, 255)
				renderer.DrawPoint(int32(x), int32(y))
			} else {
				shade := uint8(float64(numToInfinity) / float64(maxIter) * 255)
				col := color.RGBA{shade, shade, shade, 0xff}
				img.Set(x, y, col)
				renderer.SetDrawColor(0, shade, 0, 0xff)
				renderer.DrawPoint(int32(x), int32(y))
			}

		}
	}
	elapsed := time.Since(start)

	f, _ := os.Create("mandel1.png")
	png.Encode(f, img)
	fmt.Println("Time: ", elapsed)

	// surface.FillRect(nil, 0)

	// rect := sdl.Rect{0, 0, 200, 200}
	// surface.FillRect(&rect, 0xffff0000)
	// window.UpdateSurface()
	renderer.Present()
	running := true
	for running {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				println("Quit")
				running = false
				break
			}
		}
	}
}
