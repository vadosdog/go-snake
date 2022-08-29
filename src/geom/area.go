package geom

type Area struct {
	Area [][]Cell
}

type Cell struct {
	X, Y    int
	Content int
}

const (
	EmptyCell = iota
	SnakeCell = iota
	FoodCell  = iota
)

func (a Area) ForEach(f func(x int, y int, c Cell)) {
	for y, rows := range a.Area {
		for x, cell := range rows {
			f(x, y, cell)
		}
	}
}

const (
	Top    = iota
	Right  = iota
	Bottom = iota
	Left   = iota
)

type Dir int

func (d Dir) Exec(p Cell) Cell {
	switch d {
	case Top:
		return Cell{X: p.X, Y: p.Y - 1}
	case Right:
		return Cell{X: p.X + 1, Y: p.Y}
	case Bottom:
		return Cell{X: p.X, Y: p.Y + 1}
	case Left:
		return Cell{X: p.X - 1, Y: p.Y}
	}

	return Cell{X: -1, Y: -1}
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
