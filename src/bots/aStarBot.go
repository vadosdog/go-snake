package bots

import (
	"github.com/beefsack/go-astar"
	"snake/src/geom"
)

var AStarArea map[int]map[int]*Tile

type AStarBot struct {
	hasTarget bool
	path      AStarBotPath
}

func (bot *AStarBot) WhatsNext(area *geom.Area, head *geom.Cell, defaultDir geom.Dir) geom.Dir {
	fillAStarArea(area)

	if bot.hasTarget {
		if bot.path.isValid(area) {
			return getDir(head, bot.path.GetStep())
		}

		bot.hasTarget = false
	}

	bot.CalcPaths(area, head)

	if bot.hasTarget {
		return getDir(head, bot.path.GetStep())
	}

	// stay in bounds
	nextStep := defaultDir.Exec(head.Coords)
	if !area.IsFree(nextStep) {
		left := defaultDir.GetLeft()
		nextStep = left.Exec(head.Coords)
		if area.IsFree(nextStep) {
			return left
		}

		right := defaultDir.GetRight()
		nextStep = right.Exec(head.Coords)
		if area.IsFree(nextStep) {
			return right
		}
	}

	return defaultDir
}

func fillAStarArea(area *geom.Area) {
	AStarArea = make(map[int]map[int]*Tile, len(area.Area))

	for y, rows := range area.Area {
		AStarArea[y] = make(map[int]*Tile, len(area.Area[0]))
		for x, cell := range rows {
			AStarArea[y][x] = &Tile{Cell: cell, Area: area}
		}
	}
}

func (bot *AStarBot) CalcPaths(area *geom.Area, head *geom.Cell) {
	var foods []*geom.Cell
	area.ForEach(func(x int, y int, c *geom.Cell) {
		if c.IsFood() {
			foods = append(foods, c)
		}
	})

	var bestPath AStarBotPath
	// create new paths
	for _, food := range foods {
		newPath, ok := getPath(food, area, head)
		if !ok {
			continue
		}

		if bestPath.Len == 0 || newPath.Len < bestPath.Len {
			bestPath = newPath
		}
	}

	if bestPath.Len > 0 {
		bot.path = bestPath
		bot.hasTarget = true
	}
}

type AStarBotPath struct {
	Path   []*Tile
	Target *geom.Cell
	Len    float64
}

func (path *AStarBotPath) isValid(area *geom.Area) bool {
	if len(path.Path) <= 0 {
		return false
	}

	for _, tile := range path.Path {
		if !area.Area[tile.Cell.Coords.Y][tile.Cell.Coords.X].IsFree() {
			return false
		}
	}
	return true
}

func (path *AStarBotPath) GetStep() *geom.Cell {
	l := len(path.Path)
	tile := path.Path[l-1]
	path.Path = path.Path[:l-1]

	return tile.Cell
}

func getDir(head *geom.Cell, next *geom.Cell) geom.Dir {
	if next.Coords.X > head.Coords.X {
		return geom.Right
	}
	if next.Coords.X < head.Coords.X {
		return geom.Left
	}
	if next.Coords.Y > head.Coords.Y {
		return geom.Bottom
	}
	return geom.Top
}

func getPath(food *geom.Cell, area *geom.Area, head *geom.Cell) (AStarBotPath, bool) {
	path, l, found := astar.Path(AStarArea[head.Coords.Y][head.Coords.X], AStarArea[food.Coords.Y][food.Coords.X])

	if !found {
		return AStarBotPath{}, false
	}

	var tilePath []*Tile
	for i, pather := range path {
		if i == len(path)-1 {
			continue
		}

		p, ok := pather.(*Tile)
		if !ok {
			panic("AAAAAA!!!")
		}

		tilePath = append(tilePath, p)
	}

	return AStarBotPath{Path: tilePath, Len: l, Target: food}, true
}

type Tile struct {
	Cell *geom.Cell
	Area *geom.Area
}

func (t *Tile) PathNeighbors() []astar.Pather {
	neighborCells := t.Cell.Neighbors

	var neighbors []astar.Pather
	for _, cell := range neighborCells {
		if cell.IsFree() {
			neighbors = append(neighbors, AStarArea[cell.Coords.Y][cell.Coords.X])
		}
	}

	return neighbors
}

func (t *Tile) PathNeighborCost(to astar.Pather) float64 {
	return 1
}

func (t *Tile) PathEstimatedCost(to astar.Pather) float64 {
	return 1
}
