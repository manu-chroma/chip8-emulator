package main

import (
	"image"
	"image/color"
	"log"

	"golang.org/x/exp/shiny/driver"
	"golang.org/x/exp/shiny/screen"
	"golang.org/x/mobile/event/key"
	"golang.org/x/mobile/event/lifecycle"
)

// we draw graphics on screen through the use of sprites.
// A sprite is a group of bytes which are a binary representation of the desired picture.
// Chip-8 sprites can be up to 15 bytes, for a possible sprite size of 8x15
const (
	Row = 32
	Col = 64
)

// Colors
var (
	Black = color.RGBA{0, 0, 0, 1.0}
	White = color.RGBA{255, 255, 255, 1.0}
	Blue  = color.RGBA{0, 0, 255, 1.0}
)

// Screen ...
type Screen struct {
	display [Row][Col]bool
}

func (scr *Screen) clearDisplay() {

	for i := 0; i < Row; i++ {
		for j := 0; j < Col; j++ {
			scr.display[i][j] = false
		}
	}
}

// NewDisplay ...
// We pass a sender channel to the display to pass us the mouseEvents
// we obtain from the screen
func NewDisplay(mouseEvents chan<- key.Event) *Screen {
	// create a separate
	go driver.Main(func(s screen.Screen) {
		opts := screen.NewWindowOptions{
			Height: Row * 2,
			Width:  Col * 2,
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
		dim := image.Point{Col, Row}
		drawBuff, err := s.NewBuffer(dim)
		if err != nil {
			log.Fatal(err)
		}

		defaultDrawToBuffer(drawBuff.RGBA())

		log.Print("Window bounds: ", opts)
		log.Printf("Buffer bounds: %s", drawBuff.Bounds())
		log.Printf("Buffer size: %s", drawBuff.Size())

		window.Upload(image.Point{}, drawBuff, drawBuff.Bounds())
		window.Publish()

		// listening for window events
		// @discuss: should just this part be the go
		// function or the complete function?
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

			case error:
				log.Print(e)
			}

		}
	})

	log.Println("returning from NewDisplay method..")

	// return this dummy buffer for the time being
	return new(Screen)
}

// BufferToScreen puts the buffer to the window
// todo: where to put the collision detection
func BufferToScreen() {
	// This assumes that there has been updates to the current buffer
	// and now we are ready to refresh the display
	// todo: maybe we will need a back and front: separate buffers
	// for collision detection @discuss

}

// Bounds: (0,0)-(64,32)

func defaultDrawToBuffer(img *image.RGBA) {
	b := img.Bounds()

	log.Printf("Bounds: %s", b.String())

	for x := b.Min.X; x < b.Max.X; x++ {
		for y := b.Min.Y; y < b.Max.Y; y++ {
			_ = RandInRange(0, 2)

			img.SetRGBA(x, y, Blue)
			// if ran == 0 {
			// 	img.SetRGBA(x, y, White)
			// } else {
			// 	img.SetRGBA(x, y, Black)
			// }
		}
	}
}

func drawBackBuffer() {
}
