package geom

type Area struct {
	Area   [][]*Cell
	Width  int
	Height int
}

type Cell struct {
	Coords  Coords
	Content int
}

const (
	EmptyCell = iota
	SnakeCell = iota
	FoodCell  = iota
)

func CreateArea(w int, h int) *Area {
	var rows [][]*Cell
	for y := 0; y < h; y++ {
		var row []*Cell
		for x := 0; x < w; x++ {
			row = append(row, &Cell{Content: EmptyCell, Coords: Coords{X: x, Y: y}})
		}
		rows = append(rows, row)
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

func (a *Area) GetNeighborCells(c *Cell) []*Cell {
	dirs := []Coords{
		Dir(Bottom).Exec(c.Coords),
		Dir(Top).Exec(c.Coords),
		Dir(Left).Exec(c.Coords),
		Dir(Right).Exec(c.Coords),
	}

	var neighbors []*Cell
	for _, dir := range dirs {
		if dir.X >= 0 &&
			dir.Y >= 0 &&
			dir.X < a.Width &&
			dir.Y < a.Height {
			neighbors = append(neighbors, a.Area[dir.Y][dir.X])
		}
	}

	return neighbors
}
func (a *Area) GetFreeNeighborCells(c *Cell) []*Cell {
	allNeighbors := a.GetNeighborCells(c)
	var neighbors []*Cell

	for _, neighbor := range allNeighbors {
		if neighbor.IsFree() {
			neighbors = append(neighbors, neighbor)
		}
	}

	return neighbors
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
