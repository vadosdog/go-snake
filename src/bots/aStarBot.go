package bots

import (
	"github.com/beefsack/go-astar"
	"snake/src/geom"
)

type AStarBot struct {
	hasTarget bool
	path      AStarBotPath
}

func (bot *AStarBot) WhatsNext(area *geom.Area, head *geom.Cell, defaultDir geom.Dir) geom.Dir {
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

	return defaultDir
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
	tile := path.Path[0]
	path.Path = path.Path[1:]

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
		return geom.Top
	}
	return geom.Bottom
}

func getPath(food *geom.Cell, area *geom.Area, head *geom.Cell) (AStarBotPath, bool) {
	path, l, found := astar.Path(&Tile{Cell: head, Area: area}, &Tile{Cell: food, Area: area})

	if !found {
		return AStarBotPath{}, false
	}

	var tilePath []*Tile
	for _, pather := range path {
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
	neighborCells := t.Area.GetFreeNeighborCells(t.Cell)

	var neighbors []astar.Pather
	for _, cell := range neighborCells {
		neighbors = append(neighbors, &Tile{Cell: cell, Area: t.Area})
	}

	return neighbors
}

func (t *Tile) PathNeighborCost(to astar.Pather) float64 {
	return 1
}

func (t *Tile) PathEstimatedCost(to astar.Pather) float64 {
	return 1
}
