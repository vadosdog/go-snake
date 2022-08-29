package bots

import (
	"snake/src/geom"
)

type Bot interface {
	WhatsNext(area *geom.Area, head *geom.Cell, defaultDir geom.Dir) geom.Dir
}
