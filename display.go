package main

import (
	"image"
	"image/color"
	"log"

	"golang.org/x/mobile/event/paint"

	"golang.org/x/exp/shiny/driver"
	"golang.org/x/exp/shiny/screen"
	"golang.org/x/image/draw"

	"golang.org/x/mobile/event/key"

	"golang.org/x/mobile/event/lifecycle"
)

// A sprite is a group of bytes which are a binary representation of the desired picture.
// Chip-8 sprites can be up to 15 bytes, for a possible sprite size of 8x15
const (
	EmuHeight = 32
	EmuWidth  = 64

	// Scale of the main window relative to Emu dimensions
	WinScale = 20

	WinHeight = EmuHeight * WinScale
	WinWidth  = EmuWidth * WinScale

	// Acutal Emu buffer scale
	EmuScale = 20
)

// Colors
var (
	Black = color.RGBA{0, 0, 0, 1.0}
	White = color.RGBA{255, 255, 255, 1.0}
	Blue  = color.RGBA{0, 0, 255, 1.0}
)

// Screen incapsulates our display arr, window
// and the backBuffer we're using to upload to display
type Screen struct {
	display    [EmuHeight][EmuWidth]int // y for height, x for row
	window     screen.Window
	backBuffer screen.Buffer
}

func (scr *Screen) clearDisplay() {

	for j := 0; j < EmuHeight; j++ {
		for i := 0; i < EmuWidth; i++ {
			scr.display[j][i] = 0
		}
	}
}

// NewDisplay returns Screen struct instance
// We pass a sender channel to the display to pass us the mouseEvents
// which obtain from the screen
func (vm *VM) NewDisplay(mouseEvents chan<- key.Event) *Screen {

	scr := &Screen{}
	vm.screen = scr

	// create a separate
	driver.Main(func(s screen.Screen) {
		opts := screen.NewWindowOptions{
			Height: WinHeight,
			Width:  WinWidth,
			Title:  "Chip-8 VM",
		}

		window, err := s.NewWindow(&opts)
		if err != nil {
			log.Print("Unable to create display window: ")
			log.Fatal(err)
			return
		}

		defer window.Release()

		dim := image.Point{X: EmuWidth, Y: EmuHeight}
		drawBuff, err := s.NewBuffer(dim)
		if err != nil {
			log.Fatal(err)
		}
		defer drawBuff.Release()

		scr.window = window
		scr.backBuffer = drawBuff

		log.Print("Window bounds: ", opts)
		log.Printf("Buffer bounds: %s", drawBuff.Bounds())
		log.Printf("Buffer size: %s", drawBuff.Size())

		// default draw to buffer on init
		defaultDrawToBuffer(drawBuff.RGBA())
		window.Send(paint.Event{})

		// Listening for window events
		for {
			e := window.NextEvent()
			switch e := e.(type) {

			case lifecycle.Event:
				if e.To == lifecycle.StageDead {
					return
				} else if e.To == lifecycle.StageFocused {
					log.Print("Focus back on the screen!")
				}

			case key.Event:
				log.Print("pressed key: ", e.Code)
				// todo: exit game,
				// currently only shuts off the screen window
				if e.Code == key.CodeEscape {
					return
				}
				mouseEvents <- e

			case paint.Event:
				log.Println("Paint event, re-painting the buffer..")

				scaledDim := image.Rectangle{
					Max: image.Point{X: EmuWidth * EmuScale, Y: EmuHeight * EmuScale}}

				drawBuff, err = s.NewBuffer(scaledDim.Max)

				// scale image
				src := scr.backBuffer.RGBA()
				dst := image.NewRGBA(scaledDim)
				draw.BiLinear.Scale(dst, dst.Bounds(), src, src.Bounds(), draw.Src, nil)

				copyImageToBuffer(&drawBuff, dst)

				window.Upload(image.Point{}, drawBuff, drawBuff.Bounds())
				window.Publish()

			case error:
				log.Print(e)
			}

		}
	})

	return scr
}

func copyImageToBuffer(b *screen.Buffer, i *image.RGBA) {

	buffImg := (*b).RGBA()
	dim := buffImg.Bounds().Max

	for y := 0; y < dim.Y; y++ {
		for x := 0; x < dim.X; x++ {
			col := i.At(x, y)
			if col == Black {
				buffImg.SetRGBA(x, y, Black)
			} else {
				buffImg.SetRGBA(x, y, White)
			}
		}
	}

}

// BufferToScreen puts the buffer to the window
func BufferToScreen(scr *Screen) {
	// This assumes that there has been updates to the current buffer
	// and now we are ready to refresh the display

	img := scr.backBuffer.RGBA()
	b := img.Bounds()

	log.Printf("Bounds: %s", b)

	for y := 0; y < EmuHeight; y++ {
		for x := 0; x < EmuWidth; x++ {
			if scr.display[y][x] == 1 {
				img.SetRGBA(x, y, White)
			} else {
				img.SetRGBA(x, y, Black)
			}
		}
	}

	scr.window.Send(paint.Event{})
}

func defaultDrawToBuffer(img *image.RGBA) {
	b := img.Bounds()

	log.Printf("Bounds: %s", b.String())

	for x := b.Min.X; x < b.Max.X; x++ {
		for y := b.Min.Y; y < b.Max.Y; y++ {
			img.SetRGBA(x, y, Black)
		}
	}
}
