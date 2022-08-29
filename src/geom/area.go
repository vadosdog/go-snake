package geom

type Area struct {
	Area   map[int]map[int]*Cell
	Width  int
	Height int
}

type Cell struct {
	Coords    Coords
	Content   int
	Neighbors []*Cell
}

const (
	EmptyCell = iota
	SnakeCell = iota
	FoodCell  = iota
)

func CreateArea(w int, h int) *Area {
	rows := make(map[int]map[int]*Cell, h)

	for y := 0; y < h; y++ {
		rows[y] = make(map[int]*Cell, w)
		for x := 0; x < w; x++ {
			cell := &Cell{Content: EmptyCell, Coords: Coords{X: x, Y: y}}
			rows[y][x] = cell
		}
	}
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			if x > 0 {
				rows[y][x-1].Neighbors = append(rows[y][x-1].Neighbors, rows[y][x])
				rows[y][x].Neighbors = append(rows[y][x].Neighbors, rows[y][x-1])
			}

			if y > 0 {
				rows[y-1][x].Neighbors = append(rows[y-1][x].Neighbors, rows[y][x])
				rows[y][x].Neighbors = append(rows[y][x].Neighbors, rows[y-1][x])
			}
		}
	}

	return &Area{Area: rows, Width: w, Height: h}
}

func (a *Area) ForEach(f func(x int, y int, c *Cell)) {
	for y, rows := range a.Area {
		for x, cell := range rows {
			f(x, y, cell)
		}
	}
}

func (a *Area) IsFree(coords Coords) bool {
	if a.OutOfBounds(coords) {
		return false
	}

	return a.Area[coords.Y][coords.X].IsFree()
}

func (a *Area) OutOfBounds(coords Coords) bool {
	return coords.Y >= a.Height ||
		coords.X >= a.Width ||
		coords.Y < 0 ||
		coords.X < 0
}

const (
	Top    = iota
	Right  = iota
	Bottom = iota
	Left   = iota
)

type Dir int

func (d Dir) Exec(p Coords) Coords {
	switch d {
	case Top:
		return Coords{X: p.X, Y: p.Y - 1}
	case Right:
		return Coords{X: p.X + 1, Y: p.Y}
	case Bottom:
		return Coords{X: p.X, Y: p.Y + 1}
	case Left:
		return Coords{X: p.X - 1, Y: p.Y}
	}

	return Coords{X: -1, Y: -1}
}

func (d Dir) GetLeft() Dir {
	switch d {
	case Top:
		return Left
	case Right:
		return Top
	case Bottom:
		return Right
	case Left:
		return Bottom
	}

	return -1
}

func (d Dir) GetRight() Dir {
	switch d {
	case Top:
		return Right
	case Right:
		return Bottom
	case Bottom:
		return Left
	case Left:
		return Top
	}

	return -1
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

func (c *Cell) IsFood() bool {
	return c.Content == FoodCell
}

func (c *Cell) IsFree() bool {
	return c.Content == EmptyCell || c.Content == FoodCell
}
