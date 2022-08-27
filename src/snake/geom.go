package snake

type Point struct {
	X, Y float64
}

type AreaPoint struct {
	X, Y int
}

const (
	Top    = iota
	Right  = iota
	Bottom = iota
	Left   = iota
)

type Dir int

func (d Dir) Exec(p Point) Point {
	switch d {
	case Top:
		return Point{p.X, p.Y - 1}
	case Right:
		return Point{p.X + 1, p.Y}
	case Bottom:
		return Point{p.X, p.Y + 1}
	case Left:
		return Point{p.X - 1, p.Y}
	}

	return Point{-1, -1}
}

func (d Dir) CheckReverse(newDir Dir) bool {
	switch d {
	case Top:
		return newDir == Bottom
	case Right:
		return newDir == Left
	case Bottom:
		return newDir == Top
	case Left:
		return newDir == Right
	}

	return true
}
