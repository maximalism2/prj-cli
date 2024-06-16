package renderer

import (
	"github.com/gdamore/tcell/v2"
	"log"
)

var screen tcell.Screen

func Init() {
	defStyle := tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)

	// Initialize screen
	s, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("%+v", err)
	}
	if err := s.Init(); err != nil {
		log.Fatalf("%+v", err)
	}

	screen = s

	screen.SetStyle(defStyle)
	screen.Clear()
}
