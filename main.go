package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"sync"
	"time"

	"github.com/flxs/let-the-blocks-fall/field"
	"github.com/flxs/let-the-blocks-fall/gamestate"
	"github.com/gdamore/tcell"
)

var lock sync.Mutex

func main() {

	rand.Seed(time.Now().UnixNano())
	tcell.SetEncodingFallback(tcell.EncodingFallbackASCII)
	s, e := tcell.NewScreen()
	if e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e)
		os.Exit(1)
	}
	if e = s.Init(); e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e)
		os.Exit(1)
	}

	s.SetStyle(tcell.StyleDefault.
		Foreground(tcell.ColorWhite).
		Background(tcell.ColorBlack))
	s.Clear()

	quit := make(chan struct{})

	width, height := s.Size()
	pixelWidth := width / 2
	// leave room for an info bar
	height--

	//gs := gamestate.New(width/2, height)
	gs := loadOrNew(pixelWidth, height)
	lastField := gs.Field

	paused := false
	speed := 0
	baseSpeed := 21
	it := 0

	go func() {
		for {
			ev := s.PollEvent()
			switch ev := ev.(type) {
			case *tcell.EventKey:
				switch ev.Key() {
				case tcell.KeyEscape, tcell.KeyEnter, tcell.KeyCtrlC:
					close(quit)
					return
				case tcell.KeyRight:
					if !paused {
						nudge(1, false, &gs)
					}
				case tcell.KeyLeft:
					if !paused {
						nudge(-1, false, &gs)
					}
				case tcell.KeyDown:
					if !paused {
						nudge(1, true, &gs)
					}
				case tcell.KeyUp:
					if !paused {
						rotate(&gs)
					}
				case tcell.KeyRune:
					switch ev.Rune() {
					case ' ':
						paused = !paused
					case '+':
						if speed < baseSpeed-1 {
							speed++
						}
					case '-':
						if speed > -(baseSpeed - 1) {
							speed--
						}
					case 'c':
						gs = gamestate.New(width/2, height)
					}
				}
			case *tcell.EventResize:
				s.Sync()
			}
		}
	}()

	paint(&gs.Field, s)

loop:
	for {
		select {
		case <-quit:
			break loop
		case <-time.After(time.Millisecond * 20):
		}

		drawInfoBar(s, paused, gs.LinesCleared, speed)

		if paused {
			s.Show()
			continue
		}

		lock.Lock()
		if it%(baseSpeed-speed) == 0 {
			save(&gs)
			gs.Nudge(1, true)
		}

		gs.ClearCompleteLines()

		var f2 = gs.Field.Copy()
		f2.DrawBlock(gs.Block)

		paintDiff(&lastField, &f2, s)
		lock.Unlock()
		lastField = f2

		s.Show()

		it++

	}
	s.Fini()

}

func nudge(delta int, vertical bool, gs *gamestate.GameState) {
	lock.Lock()
	defer lock.Unlock()
	gs.Nudge(delta, vertical)
}

func rotate(gs *gamestate.GameState) {
	lock.Lock()
	defer lock.Unlock()
	gs.RotateBlock()
}

func save(gs *gamestate.GameState) {
	file, _ := json.Marshal(gs)
	_ = ioutil.WriteFile("state.json", file, 0644)
}

func loadOrNew(w int, h int) gamestate.GameState {
	_, err := os.Stat("state.json")
	if os.IsNotExist(err) {
		return gamestate.New(w, h)
	}

	file, _ := ioutil.ReadFile("state.json")

	var data gamestate.GameState

	_ = json.Unmarshal([]byte(file), &data)

	return data
}

func drawString(str string, x int, y int, style tcell.Style, s tcell.Screen) {
	for i, c := range str {
		s.SetCell(x+i, y, style, c)
	}
}

func drawInfoBar(s tcell.Screen, paused bool, lines int, speed int) {
	style := tcell.StyleDefault.
		Background(tcell.NewRGBColor(80, 80, 80)).
		Foreground(tcell.NewRGBColor(220, 220, 220)).
		Bold(true)

	width, _ := s.Size()

	for i := 0; i < width; i++ {
		s.SetCell(i, 0, style, ' ')
	}

	if paused {
		centerStr := "PAUSED"
		drawString(centerStr, width/2-(len(centerStr)/2), 0, style, s)
	} else {
		centerStr := "[Space] to pause, [C] to clear"
		drawString(centerStr, width/2-(len(centerStr)/2), 0, style, s)
	}

	linesStr := fmt.Sprintf("Lines: %d ", lines)
	drawString(linesStr, width-len(linesStr), 0, style, s)

	speedStr := fmt.Sprintf("Speed: %d", speed)
	drawString(speedStr, 0, 0, style, s)

}

func styleForColor(color int) tcell.Style {
	colors := []tcell.Color{
		tcell.NewRGBColor(50, 50, 50),
		tcell.ColorYellow,
		tcell.ColorRoyalBlue,
		tcell.ColorRed,
		tcell.ColorOrange,
		tcell.ColorLawnGreen,
		tcell.ColorMediumVioletRed,
		tcell.ColorMistyRose,
	}
	st := tcell.StyleDefault
	return st.Background(colors[color])
}

func paint(f *field.Field, s tcell.Screen) {
	for i := 0; i < len(f.Matrix); i++ {
		x := i % f.Width
		y := i/f.Width + 1
		style := styleForColor(f.Matrix[i])
		gl := ' '

		s.SetCell(2*x, y, style, gl)
		s.SetCell(2*x+1, y, style, gl)
	}
}

func paintDiff(fOld *field.Field, f *field.Field, s tcell.Screen) {
	for i := 0; i < len(f.Matrix); i++ {
		if fOld.Matrix[i] == f.Matrix[i] {
			continue
		}
		x := i % f.Width
		y := i/f.Width + 1 // +1 for info bar
		style := styleForColor(f.Matrix[i])
		gl := ' '

		s.SetCell(2*x, y, style, gl)
		s.SetCell(2*x+1, y, style, gl)
	}
}
