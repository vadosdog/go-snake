package bots

import (
	"math/rand"
	"snake/src/geom"
)

type RandomBot struct {
}

func (bot *RandomBot) WhatsNext(area geom.Area, head geom.Cell) geom.Dir {
	dirs := []int{geom.Top, geom.Bottom, geom.Left, geom.Right}
	i := rand.Intn(3)

	return geom.Dir(dirs[i])
}
