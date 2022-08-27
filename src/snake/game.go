package snake

import (
	"math/rand"
	"sort"
	"time"
)

const (
	GameW     = float64(720)
	GameH     = float64(720)
	SpeedStep = 10

	SinglePlayer = iota
	MultiPlayer  = iota
)

type Game struct {
	Snakes   []*Snake
	gameOver bool
	speed    int

	food         []Point
	playersCount int

	mode int
}

func NewGame(playerCounts int) *Game {
	var snakes []*Snake
	snakes = []*Snake{NewSnake("Player 1", "#fff", Point{1, 1}, Right, map[int]Dir{CodeLeft: Left, CodeRight: Right, CodeUp: Top, CodeDown: Bottom})}

	if playerCounts >= 2 {
		snakes = append(snakes, NewSnake("Player 2", "#ff0", Point{18, 18}, Left, map[int]Dir{CodeA: Left, CodeD: Right, CodeW: Top, CodeS: Bottom}))
	}
	if playerCounts >= 3 {
		snakes = append(snakes, NewSnake("Player 3", "#f00", Point{1, 18}, Top, map[int]Dir{CodeNum4: Left, CodeNum6: Right, CodeNum8: Top, CodeNum5: Bottom}))
	}
	if playerCounts == 4 {
		snakes = append(snakes, NewSnake("Player 4", "#f0f", Point{18, 1}, Bottom, map[int]Dir{CodeJ: Left, CodeL: Right, CodeI: Top, CodeK: Bottom}))
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

	return g
}

func (g *Game) HandleKeyDown(code int, gp *GamePage) {
	for _, snake := range g.Snakes {
		if snake.NeedMove {
			continue
		}
		newDir, ok := snake.Controllers[code]
		if !ok {
			continue
		}

		if newDir != snake.Dir && !snake.Dir.CheckReverse(newDir) {
			snake.Dir = newDir
			snake.NeedMove = true
		}
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

			newPoint := Point{float64(randX), float64(randY)}

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
	//var snakeLock sync.Mutex

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
