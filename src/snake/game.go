package snake

import (
	"math/rand"
	"snake/src/bots"
	"snake/src/geom"
	"sort"
	"time"
)

const (
	SpeedStep = 10

	SinglePlayer = iota
	MultiPlayer  = iota

	AreaCellCountW = 20
	AreaCellCountH = 20
)

type Game struct {
	Snakes   []*Snake
	gameOver bool
	speed    int

	food         []geom.Cell
	playersCount int

	mode int

	Area geom.Area
}

func NewGame(playerCounts int) *Game {
	var snakes []*Snake
	snakes = []*Snake{NewSnake("Player 1", "#fff", geom.Cell{X: 1, Y: 1}, geom.Right, map[int]geom.Dir{CodeLeft: geom.Left, CodeRight: geom.Right, CodeUp: geom.Top, CodeDown: geom.Bottom})}
	// refactor
	snakes[0].SetBot(&bots.RandomBot{})

	if playerCounts >= 2 {
		snakes = append(snakes, NewSnake("Player 2", "#ff0", geom.Cell{X: 18, Y: 18}, geom.Left, map[int]geom.Dir{CodeA: geom.Left, CodeD: geom.Right, CodeW: geom.Top, CodeS: geom.Bottom}))
	}
	if playerCounts >= 3 {
		snakes = append(snakes, NewSnake("Player 3", "#f00", geom.Cell{X: 1, Y: 18}, geom.Top, map[int]geom.Dir{CodeNum4: geom.Left, CodeNum6: geom.Right, CodeNum8: geom.Top, CodeNum5: geom.Bottom}))
	}
	if playerCounts == 4 {
		snakes = append(snakes, NewSnake("Player 4", "#f0f", geom.Cell{X: 18, Y: 1}, geom.Bottom, map[int]geom.Dir{CodeJ: geom.Left, CodeL: geom.Right, CodeI: geom.Top, CodeK: geom.Bottom}))
	}

	gameMode := SinglePlayer
	if playerCounts > 1 {
		gameMode = MultiPlayer
	}

	g := &Game{
		Snakes:       snakes,
		speed:        500,
		gameOver:     false,
		playersCount: playerCounts,
		mode:         gameMode,
	}

	g.fillArea()

	return g
}

func (g *Game) HandleKeyDown(code int, gp *GamePage) {
	for _, snake := range g.Snakes {
		newDir, ok := snake.Controllers[code]
		if !ok {
			continue
		}

		snake.ChangeDir(newDir)
	}
}

func (g *Game) foodGeneration() {
	var foodTimer *time.Timer

	resetTimer := func() {
		foodTimer = time.NewTimer(time.Duration(3) * time.Second)
	}
	resetTimer()

	for {
		<-foodTimer.C

		if !g.gameOver {
			randX := rand.Intn(20)
			randY := rand.Intn(20)

			newPoint := geom.Cell{X: randX, Y: randY}

			check := true

			for _, snake := range g.Snakes {
				if snake.IsSnake(newPoint) {
					check = false
					break
				}
			}

			if !check {
				for _, p := range g.food {
					if p.X == newPoint.X && p.Y == newPoint.Y {
						check = false
						break
					}
				}
			}

			if check {
				g.food = append(g.food, newPoint)
			}
		}

		resetTimer()
	}
}

func (g *Game) Run() {
	go g.snakeMovement()
	go g.foodGeneration()
}

//
func (g *Game) snakeMovement() {
	var snakeTimer *time.Timer

	resetTimer := func() {
		snakeTimer = time.NewTimer(time.Duration(g.speed) * time.Millisecond)
	}

	resetTimer()

	//loop
	for {
		<-snakeTimer.C

		if g.gameOver {
			break
		}

		for _, snake := range g.Snakes {
			if snake.IsLose {
				continue
			}

			if snake.isAi {
				snake.ChangeDir(snake.bot.WhatsNext(g.Area, snake.Head()))
			}

			newPos := snake.Dir.Exec(snake.Head())

			//food
			isFood := false
			for i := range g.food {
				if newPos.X == g.food[i].X && newPos.Y == g.food[i].Y {
					g.food = append(g.food[:i], g.food[i+1:]...)
					snake.Add(newPos)
					if g.speed > SpeedStep {
						g.speed -= SpeedStep
					}
					isFood = true
					snake.Score++
					break
				}
			}

			if !isFood {
				snake.Move()
			}

			snake.NeedMove = false
		}

		// Check Lose
		loses := 0
		for _, snake := range g.Snakes {
			if snake.IsLose {
				continue
			}

			newPos := snake.Head()

			lose := false
			// Check walls
			if newPos.X < 0 || newPos.X >= 20 ||
				newPos.Y < 0 || newPos.Y >= 20 {
				lose = true
			}

			// Check snakes
			for _, checkSnake := range g.Snakes {
				if checkSnake == snake {
					if checkSnake.IsSnakeTail(newPos) {
						lose = true
					}
				} else {
					if checkSnake.IsSnake(newPos) {
						lose = true
					}
				}
			}

			if lose {
				snake.IsLose = true
				loses++
			}
		}

		// add score to other snakes
		continues := 0
		for _, snake := range g.Snakes {
			if !snake.IsLose {
				snake.Score += 5 * loses
				continues++
			}
		}

		if continues == 0 {
			g.gameOver = true
		}

		g.fillArea()
		resetTimer()
	}
}

func (g *Game) getSortedSnakes() []*Snake {
	sorted := make([]*Snake, len(g.Snakes))

	copy(sorted, g.Snakes)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Score > sorted[j].Score
	})

	return sorted
}

func (g *Game) fillArea() {
	if g.gameOver {
		return
	}
	var rows [][]geom.Cell
	for y := 0; y < AreaCellCountH; y++ {
		var row []geom.Cell
		for x := 0; x < AreaCellCountW; x++ {
			row = append(row, geom.Cell{Content: geom.EmptyCell, X: x, Y: y})
		}
		rows = append(rows, row)
	}

	for _, snake := range g.Snakes {
		for _, point := range snake.Parts {
			rows[int(point.Y)][int(point.X)] = geom.Cell{Content: geom.SnakeCell, X: int(point.X), Y: int(point.Y)}
		}
	}

	for _, point := range g.food {
		rows[int(point.Y)][int(point.X)] = geom.Cell{Content: geom.FoodCell, X: int(point.X), Y: int(point.Y)}
	}

	g.Area = geom.Area{Area: rows}
}
