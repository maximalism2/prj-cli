package renderer

import (
	"github.com/gdamore/tcell/v2"
	"os"
	"strconv"
)

type View struct {
	Projects []string
}

var selectedProjectChan chan string
var focusedIdx int
var rowOffset int

var inputValue []rune

func GetSelectedProject() <-chan string {
	return selectedProjectChan
}

func cleanup() {
	// You have to catch panics in a defer, clean up, and
	// re-raise them - otherwise your application can
	// die without leaving any diagnostic trace.
	maybePanic := recover()
	screen.Fini()
	if maybePanic != nil {
		panic(maybePanic)
	}
}

func RenderView(view View) {
	defer cleanup()

	selectedProjectChan = make(chan string)
	inputValue = []rune{}
	focusedIdx = 0
	rowOffset = 0

	for {
		screen.Clear()
		render(view)
		renderInput(inputValue)

		renderDebugInfo(view)
		screen.Show()

		ev := screen.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventResize:
			screen.Sync()
		case *tcell.EventKey:
			r := ev.Rune()
			if ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyCtrlC {
				cleanup()
				os.Exit(0)
				return
			} else if ev.Key() == tcell.KeyCtrlL {
				screen.Sync()
			} else if ev.Key() == tcell.KeyEnter {
				selectedIdx := focusedIdx + rowOffset
				if selectedIdx >= len(view.Projects) {
					continue
				}

				selectedProjectChan <- view.Projects[selectedIdx]
				return
			} else if ev.Key() == tcell.KeyDown {
				_, maxRows := screen.Size()

				if focusedIdx+rowOffset+1 == len(view.Projects) {
					continue
				}

				if len(view.Projects)-rowOffset+3 < maxRows {
					focusedIdx += 1
					continue
				}

				if (focusedIdx + 10) > maxRows {
					rowOffset += 1
					continue
				}

				focusedIdx += 1
			} else if ev.Key() == tcell.KeyUp {
				if focusedIdx-1 < 3 && rowOffset > 0 {
					rowOffset -= 1
					continue
				}

				focusedIdx -= 1
				if focusedIdx < 0 {
					focusedIdx = 0
				}
			} else if ev.Key() == tcell.KeyBackspace || ev.Key() == tcell.KeyBackspace2 {
				if len(inputValue) > 0 {
					inputValue = inputValue[:len(inputValue)-1]
				}
				continue
			} else if isQuickSelectRune(r) {
				idx := int(r - '1')
				if r == '0' {
					idx = 9
				}

				if idx >= len(view.Projects) {
					continue
				}

				selectedProjectChan <- view.Projects[idx]
				return
			} else if isSupportedRune(r) {
				inputValue = append(inputValue, r)
				focusedIdx = 0
				rowOffset = 0
				continue
			}
		}
	}
}

func renderDebugInfo(view View) {
	maxWidth, maxHeight := screen.Size()
	start := maxWidth - 20
	style := tcell.StyleDefault.Foreground(tcell.ColorReset).Background(tcell.ColorReset)

	projectsLength := "Projects len: " + strconv.Itoa(len(view.Projects))
	focusedIdxLog := "Focused idx: " + strconv.Itoa(focusedIdx)
	rowOffsetLog := "Row offset: " + strconv.Itoa(rowOffset)
	screenSizeLog := "Screen: " + strconv.Itoa(maxWidth) + " ⨉ " + strconv.Itoa(maxHeight)

	for idx, r := range []rune(projectsLength) {
		screen.SetContent(idx+start, 2, r, nil, style)
	}

	for idx, r := range []rune(rowOffsetLog) {
		screen.SetContent(idx+start, 3, r, nil, style)
	}

	for idx, r := range []rune(focusedIdxLog) {
		screen.SetContent(idx+start, 4, r, nil, style)
	}

	for idx, r := range []rune(screenSizeLog) {
		screen.SetContent(idx+start, 5, r, nil, style)
	}
}

func renderInput(value []rune) {
	_, maxRows := screen.Size()
	prefix := []rune("Filter: ")
	prefixLen := len(prefix)
	valueLen := len(value)

	textStyle := tcell.StyleDefault.Foreground(tcell.ColorReset).Background(tcell.ColorReset)

	screen.ShowCursor(prefixLen+valueLen, maxRows-1)

	for idx, t := range prefix {
		screen.SetContent(idx, maxRows-1, t, nil, textStyle)
	}

	for idx, t := range value {
		screen.SetContent(idx+prefixLen, maxRows-1, t, nil, textStyle)
	}

	for idx := valueLen + prefixLen + 1; idx < maxRows; idx++ {
		screen.SetContent(idx, maxRows-1, ' ', nil, textStyle)
	}
}

func isQuickSelectRune(r rune) bool {
	return r >= '0' && r <= '9'
}

func isSupportedRune(r rune) bool {
	if r >= 'A' && r <= 'z' {
		return true
	}

	return r == '-' || r == '_' || r == '.' || r == '/' || r == ' '
}

func render(view View) {
	_, maxRows := screen.Size()

	for i := 0; i < maxRows-4; i++ {
		idxWithOffset := i + rowOffset
		if idxWithOffset == len(view.Projects) {
			return
		}

		renderProjectName(i, view.Projects[idxWithOffset])
	}
}

func renderProjectName(row int, text string) {
	textStyle := tcell.StyleDefault.Foreground(tcell.ColorReset).Background(tcell.ColorReset)

	var prefix rune
	if row < 9 {
		prefix = rune('0' + row + 1)
	} else if row == 9 {
		prefix = '0'
	} else {
		prefix = ' '
	}

	screen.SetContent(0, row, prefix, nil, textStyle)

	if row == focusedIdx {
		textStyle = tcell.StyleDefault.Foreground(tcell.ColorReset).Background(tcell.ColorGreen)
	}

	for idx, t := range []rune(text) {
		screen.SetContent(idx+2, row, t, nil, textStyle)
	}
}
