package snake

const BaseLen = 5

type Snake struct {
	Parts       []Point
	Color       string
	Dir         Dir
	NeedMove    bool
	Controllers map[int]Dir
	IsLose      bool
	Score       int
	Name        string
}

func NewSnake(name string, color string, point Point, dir Dir, controllers map[int]Dir) *Snake {
	snake := &Snake{Color: color, Dir: dir, Controllers: controllers, Name: name}
	snake.Reset(point)
	return snake
}

func (s *Snake) Len() int {
	return len(s.Parts)
}

func (s *Snake) Head() Point {
	if s.Len() == 0 {
		return Point{-1, -1}
	}

	return s.Parts[0]
}

func (s *Snake) Tail() Point {
	if s.Len() == 0 {
		return Point{-1, -1}
	}

	return s.Parts[s.Len()-1]
}

func (s *Snake) Add(p Point) {
	s.Parts = append([]Point{p}, s.Parts...)
}

func (s *Snake) IsSnake(p Point) bool {
	for _, i := range s.Parts {
		if i == p {
			return true
		}
	}
	return false
}

func (s *Snake) IsSnakeTail(p Point) bool {
	for _, i := range s.Parts[1:] {
		if i == p {
			return true
		}
	}
	return false
}

func (s *Snake) CutIfSnake(p Point) bool {
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

func (s *Snake) Reset(point Point) {
	sx, sy, l := point.X, point.Y, BaseLen
	s.Parts = []Point{}
	for i := l - 1; i >= 0; i-- {
		switch s.Dir {
		case Left:
			s.Parts = append(s.Parts, Point{sx - float64(i), sy})
		case Top:
			s.Parts = append(s.Parts, Point{sx, sy - float64(i)})
		case Right:
			s.Parts = append(s.Parts, Point{sx + float64(i), sy})
		case Bottom:
			s.Parts = append(s.Parts, Point{sx, sy + float64(i)})
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
