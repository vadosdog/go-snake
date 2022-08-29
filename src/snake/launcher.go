package snake

import (
	"github.com/tfriedel6/canvas"
	"github.com/tfriedel6/canvas/sdlcanvas"
)

const (
	wW = 1080
	wH = 750
)

type Launcher struct {
	cv  *canvas.Canvas
	wnd *sdlcanvas.Window

	page Page

	font *canvas.Font
}

type Page interface {
	HandleKeyDown(code int, launcher *Launcher)
	Render(l *Launcher)
	Run(l *Launcher)
}

func NewLauncher() *Launcher {
	wnd, cv, err := sdlcanvas.CreateWindow(wW, wH, "Go, snake, go!")

	if err != nil {
		panic(err)
	}

	font, err := cv.LoadFont("./tahoma.ttf")
	if err != nil {
		panic(err.Error())
	}

	l := &Launcher{
		cv:   cv,
		wnd:  wnd,
		page: getWelcomePage(),
		font: font,
	}

	return l
}

func (l *Launcher) Run() {
	go l.handleKeyUp()
	l.renderLoop()
}

func (l *Launcher) renderLoop() {
	l.wnd.MainLoop(func() {
		// clear
		l.cv.ClearRect(0, 0, 1080, 750)
		l.page.Render(l)
	})
}

func (l *Launcher) handleKeyUp() {
	l.wnd.KeyDown = func(code int, rn rune, name string) {
		l.page.HandleKeyDown(code, l)
	}
}

func (l *Launcher) changePage(page Page) {
	l.page = page
	l.page.Run(l)
}
