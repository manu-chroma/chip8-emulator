package main

import (
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
			Height: Col,
			Width:  Row,
			Title:  "Chip-8 VM",
		}

		w, err := s.NewWindow(&opts)
		if err != nil {
			log.Print("Unable to create dispaly window: ")
			log.Fatal(err)
			return
		}

		defer w.Release()

		// looping over to hear for window events
		// @discuss: should just this part be the go function or the complete function?
		for {
			e := w.NextEvent()
			switch e := e.(type) {
			case lifecycle.Event:
				if e.To == lifecycle.StageDead {
					return
				} else if e.To == lifecycle.StageFocused {
					log.Print("Focus back on the screen!")
				} else {

				}

			case key.Event:
				// exit game
				if e.Code == key.CodeEscape {
					return
				}

				log.Print("pressed key: ", e.Code)
				mouseEvents <- e
				log.Println("Sent the message to mouseEvents channel")

			case error:
				log.Print(e)

			}

		}
	})

	log.Println("returning from NewDisplay method..")

	// return this dummy buffer for the time being
	return new(Screen)
}
