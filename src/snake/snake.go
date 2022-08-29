package snake

import (
	"snake/src/bots"
	"snake/src/geom"
)

const BaseLen = 5

type Snake struct {
	Parts       []geom.Cell
	Color       string
	Dir         geom.Dir
	NeedMove    bool
	Controllers map[int]geom.Dir
	IsLose      bool
	Score       int
	Name        string
	bot         bots.Bot
	isAi        bool
}

func NewSnake(name string, color string, point geom.Cell, dir geom.Dir, controllers map[int]geom.Dir) *Snake {
	snake := &Snake{Color: color, Dir: dir, Controllers: controllers, Name: name}
	snake.Reset(point)
	return snake
}

func (s *Snake) ChangeDir(newDir geom.Dir) {
	if s.NeedMove {
		return
	}

	if newDir != s.Dir && !s.Dir.CheckReverse(newDir) {
		s.Dir = newDir
		s.NeedMove = true
	}
}

func (s *Snake) Len() int {
	return len(s.Parts)
}

func (s *Snake) Head() geom.Cell {
	if s.Len() == 0 {
		return geom.Cell{X: -1, Y: -1, Content: geom.SnakeCell}
	}

	return s.Parts[0]
}

func (s *Snake) Tail() geom.Cell {
	if s.Len() == 0 {
		return geom.Cell{X: -1, Y: -1, Content: geom.SnakeCell}
	}

	return s.Parts[s.Len()-1]
}

func (s *Snake) Add(p geom.Cell) {
	p.Content = geom.SnakeCell
	s.Parts = append([]geom.Cell{p}, s.Parts...)
}

func (s *Snake) IsSnake(p geom.Cell) bool {
	for _, i := range s.Parts {
		if i == p {
			return true
		}
	}
	return false
}

func (s *Snake) IsSnakeTail(p geom.Cell) bool {
	for _, i := range s.Parts[1:] {
		if i == p {
			return true
		}
	}
	return false
}

func (s *Snake) CutIfSnake(p geom.Cell) bool {
	i := 0
	for ; i < s.Len(); i++ {
		if s.Parts[i] == p {
			break
		}
	}

	if i >= s.Len() {
		return false
	}

	s.Parts = s.Parts[0:i]

	return true
}

func (s *Snake) Reset(point geom.Cell) {
	sx, sy, l := point.X, point.Y, BaseLen
	s.Parts = []geom.Cell{}
	for i := l - 1; i >= 0; i-- {
		switch s.Dir {
		case geom.Left:
			s.Parts = append(s.Parts, geom.Cell{X: sx - i, Y: sy, Content: geom.SnakeCell})
		case geom.Top:
			s.Parts = append(s.Parts, geom.Cell{X: sx, Y: sy - i, Content: geom.SnakeCell})
		case geom.Right:
			s.Parts = append(s.Parts, geom.Cell{X: sx + i, Y: sy, Content: geom.SnakeCell})
		case geom.Bottom:
			s.Parts = append(s.Parts, geom.Cell{X: sx, Y: sy + i, Content: geom.SnakeCell})
		}
	}
}

func (s *Snake) Move() {
	lastP := s.Head()
	s.Parts[0] = s.Dir.Exec(s.Head())
	for i := range s.Parts[1:] {
		s.Parts[i+1], lastP = lastP, s.Parts[i+1]
	}
}

func (s *Snake) SetBot(b bots.Bot) {
	s.isAi = true
	s.bot = b
}
