package snake

import (
	"github.com/tfriedel6/canvas"
)

type Button struct {
	title        string
	playerCounts int
}

type WelcomePage struct {
	selectedButton int
	Buttons        []Button
}

func getWelcomePage() *WelcomePage {
	buttons := []Button{
		{title: "1 Player (Up, Left, Down, Right)", playerCounts: 1},
		{title: "2 Players (W, A, S, D)", playerCounts: 2},
		{title: "3 Players (Num8, Num4, Num5, Num6)", playerCounts: 3}, //96 92 93 94
		{title: "4 Players (I, J, K, L)", playerCounts: 4},             // 12 13 14 15
	}
	return &WelcomePage{
		Buttons: buttons,
	}
}

func (wp *WelcomePage) Render(l *Launcher) {
	l.cv.BeginPath()

	//Title
	l.cv.SetFont(l.font, 25)
	l.cv.SetFillStyle("#AAA")
	l.cv.SetTextAlign(canvas.Center)
	text := "Do you wanna play?"
	l.cv.FillText(text, 1080/2, 180)
	l.cv.SetFont(l.font, 52)
	l.cv.SetFillStyle("#FFAAFF")
	text = "GO PLAY!"
	l.cv.FillText(text, 1080/2, 230)

	//Buttons
	for i, b := range wp.Buttons {
		if i == wp.selectedButton {
			l.cv.SetFillStyle("#FFAAFF")
			l.cv.SetFont(l.font, 30)
		} else {
			l.cv.SetFillStyle("#333")
			l.cv.SetFont(l.font, 25)
		}

		l.cv.FillText(b.title, 1080/2, float64(500+i*50))
	}

	l.cv.Stroke()
}

func (wp *WelcomePage) HandleKeyDown(code int, launcher *Launcher) {
	if code == CodeUp {
		wp.selectedButton--
		if wp.selectedButton < 0 {
			wp.selectedButton = 3
		}
	}

	if code == CodeDown {
		wp.selectedButton++
		if wp.selectedButton > 3 {
			wp.selectedButton = 0
		}
	}

	if code == CodeEnter {
		launcher.changePage(getGamePage(wp.Buttons[wp.selectedButton].playerCounts))
	}
}

func (wp *WelcomePage) Run(l *Launcher) {
}
