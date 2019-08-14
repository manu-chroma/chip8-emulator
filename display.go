package main

import (
	"image"
	"image/color"
	"log"

	"golang.org/x/mobile/event/paint"

	"golang.org/x/exp/shiny/driver"
	"golang.org/x/exp/shiny/screen"
	"golang.org/x/mobile/event/key"
	"golang.org/x/mobile/event/lifecycle"
)

// we draw graphics on screen through the use of sprites.
// A sprite is a group of bytes which are a binary representation of the desired picture.
// Chip-8 sprites can be up to 15 bytes, for a possible sprite size of 8x15
const (
	WinRow = 300
	WinCol = 600
	EmuRow = 32
	EmuCol = 64
)

// Colors
var (
	Black = color.RGBA{0, 0, 0, 1.0}
	White = color.RGBA{255, 255, 255, 1.0}
	Blue  = color.RGBA{0, 0, 255, 1.0}
)

// Screen ...
type Screen struct {
	display    [EmuRow][EmuCol]int
	window     screen.Window
	backBuffer screen.Buffer
}

func (scr *Screen) clearDisplay() {

	for i := 0; i < EmuRow; i++ {
		for j := 0; j < EmuCol; j++ {
			scr.display[i][j] = 0
		}
	}

	BufferToScreen(scr)
}

// NewDisplay ...
// We pass a sender channel to the display to pass us the mouseEvents
// we obtain from the screen
func NewDisplay(mouseEvents chan<- key.Event) *Screen {

	scr := &Screen{}

	// create a separate
	go driver.Main(func(s screen.Screen) {
		opts := screen.NewWindowOptions{
			Height: WinRow,
			Width:  WinCol,
			Title:  "Chip-8 VM",
		}

		window, err := s.NewWindow(&opts)
		if err != nil {
			log.Print("Unable to create display window: ")
			log.Fatal(err)
			return
		}

		defer window.Release()

		// create basic gradient
		// @bug why are we needing col * 2 rather than col?
		dim := image.Point{EmuCol, EmuRow}
		drawBuff, err := s.NewBuffer(dim)

		scr.window = window
		scr.backBuffer = drawBuff

		if err != nil {
			log.Fatal(err)
		}

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
				// exit game
				if e.Code == key.CodeEscape {
					return
				}
				mouseEvents <- e

			case paint.Event:
				log.Print("Paint event, re-painting the buffer..")
				// resize image
				// dim2 := image.Point{Col * 2, Row * 2}
				// newBuff, _ := s.NewBuffer(dim2)
				// newImage := resize.Resize(Row*2, Col*2, drawBuff.RGBA(), resize.Lanczos3)
				// _ = jpeg.Encode(newBuff, newImage, nil)
				window.Upload(image.Point{}, drawBuff, drawBuff.Bounds())
				window.Publish()

			case error:
				log.Print(e)
			}

		}
	})

	log.Println("returning from NewDisplay method..")

	// return this dummy buffer for the time being
	return scr
}

// BufferToScreen puts the buffer to the window
// todo: where to put the collision detection
func BufferToScreen(scr *Screen) {
	// This assumes that there has been updates to the current buffer
	// and now we are ready to refresh the display

	img := scr.backBuffer.RGBA()
	b := img.Bounds()

	log.Printf("Bounds: %s", b)

	for x := 0; x < EmuRow; x++ {
		for y := 0; y < EmuCol; y++ {
			if scr.display[x][y] == 1 {
				img.SetRGBA(x, y, White)
			} else {
				img.SetRGBA(x, y, Black)
			}
		}
	}

	// send screen paint event
	scr.window.Send(paint.Event{})
}

// Bounds: (0,0)-(64,32)

func defaultDrawToBuffer(img *image.RGBA) {
	b := img.Bounds()

	log.Printf("Bounds: %s", b.String())

	for x := b.Min.X; x < b.Max.X; x++ {
		for y := b.Min.Y; y < b.Max.Y; y++ {
			ran := RandInRange(0, 2)

			// img.SetRGBA(x, y, Blue)
			if ran == 0 {
				img.SetRGBA(x, y, White)
			} else {
				img.SetRGBA(x, y, Black)
			}
		}
	}
}

func drawBackBuffer() {
}
