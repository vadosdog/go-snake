package snake

import (
	"fmt"
	"github.com/tfriedel6/canvas"
)

const (
	AreaW          = float64(720)
	AreaH          = float64(720)
	AreaCellCountW = 20
	AreaCellCountH = 20
)

type GamePage struct {
	Game       *Game
	gameAreaSP Point
	gameAreaEP Point
	cellW      float64
	cellH      float64
}

func getGamePage(payerCounts int) *GamePage {
	return &GamePage{
		Game: NewGame(payerCounts),

		gameAreaSP: Point{15, 15},
		gameAreaEP: Point{15 + AreaW, 15 + AreaH},

		cellW: AreaW / AreaCellCountW,
		cellH: AreaH / AreaCellCountH,
	}
}

func (gp *GamePage) Render(l *Launcher) {
	l.cv.BeginPath()

	l.cv.SetTextAlign(canvas.Left)

	//render area
	l.cv.SetFillStyle("#333")
	l.cv.FillRect(gp.gameAreaSP.X, gp.gameAreaSP.Y, AreaW, AreaH)

	l.cv.SetStrokeStyle("#FFF001")
	l.cv.SetLineWidth(1)
	for i := 0; i < 20+1; i++ {
		l.cv.MoveTo(gp.gameAreaSP.X+float64(i)*gp.cellW, gp.gameAreaSP.Y)
		l.cv.LineTo(gp.gameAreaSP.X+float64(i)*gp.cellW, gp.gameAreaEP.Y)

		l.cv.MoveTo(gp.gameAreaSP.X, gp.gameAreaSP.Y+float64(i)*gp.cellH)
		l.cv.LineTo(gp.gameAreaEP.X, gp.gameAreaSP.Y+float64(i)*gp.cellH)
	}

	// render snakes

	for _, snake := range gp.Game.Snakes {
		l.cv.SetFillStyle(snake.Color)
		for _, p := range snake.Parts {
			l.cv.FillRect(
				gp.gameAreaSP.X+p.X*gp.cellW+1,
				gp.gameAreaSP.Y+p.Y*gp.cellH+1,
				gp.cellW-1*2,
				gp.cellH-1*2,
			)
		}
	}

	// render food
	l.cv.SetFillStyle("#F15555")
	for _, p := range gp.Game.food {
		l.cv.FillRect(
			gp.gameAreaSP.X+p.X*gp.cellW+1,
			gp.gameAreaSP.Y+p.Y*gp.cellH+1,
			gp.cellW-1*2,
			gp.cellH-1*2,
		)
	}

	// render score
	lineX := float64(720 + 25)
	lineY := float64(50)
	l.cv.SetFont(l.font, 25)

	text := fmt.Sprintf("Food: %d", len(gp.Game.food))
	l.cv.FillText(text, lineX, lineY)
	text = fmt.Sprintf("Speed: %d", (500+SpeedStep-gp.Game.speed)/SpeedStep)
	lineY += 35
	l.cv.FillText(text, lineX, lineY)
	for _, snake := range gp.Game.getSortedSnakes() {
		l.cv.SetFillStyle(snake.Color)
		text = fmt.Sprintf("%s Score: %d", snake.Name, snake.Score)

		lineY += 35
		l.cv.FillText(text, lineX, lineY)
	}

	if gp.Game.gameOver {
		l.cv.SetFillStyle("#F15555")
		l.cv.SetFont(l.font, 30)
		text = "Game over"
		lineY += 50
		l.cv.FillText(text, lineX, lineY)

		l.cv.SetFillStyle("#FFAAFF")
		l.cv.SetFont(l.font, 30)
		lineY += 50
		l.cv.FillText("Enter to return", lineX, lineY)
	}

	l.cv.Stroke()
}

func (gp *GamePage) HandleKeyDown(code int, launcher *Launcher) {
	if gp.Game.gameOver {
		if code == CodeEnter {
			launcher.changePage(getWelcomePage())
		}

		return
	}

	gp.Game.HandleKeyDown(code, gp)
}

func (gp *GamePage) Run(l *Launcher) {
	gp.Game.Run()
}
