package snake

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

// Directions
const (
	North = iota
	East
	South
	West
)

// Snake ...
type Snake []int

func (s Snake) dir() int {
	h := s[0]
	n := s[1]
	switch n - h {
	case Width:
		return North
	case 1:
		return West
	case -Width:
		return South
	case -1:
		return East
	}

	return -1
}

func (s Snake) tailDir() int {
	l := s.len()
	t := s[l-1]
	b := s[l-2]
	switch t - b {
	case Width:
		return North
	case 1:
		return West
	case -Width:
		return South
	case -1:
		return East
	}

	return -1
}

func (s Snake) len() int {
	return len(s)
}

func getSnakeHeadGeoM(index, d, offset int) (float64, float64) {
	dx, dy := 0, 0
	x := index % Width * CellSize
	y := index / Width * CellSize
	switch d {
	case North:
		dx, dy = 0, -1
	case East:
		dx, dy = 1, 0
	case South:
		dx, dy = 0, 1
	case West:
		dx, dy = -1, 0
	}
	x += dx * offset
	y += dy * offset

	return float64(x), float64(y)
}

func getSnakeNeckGeoM(index, d, offset int) (float64, float64) {
	dx, dy := 0, 0
	x := index % Width * CellSize
	y := index / Width * CellSize
	switch d {
	case North:
		dx, dy = 0, -1
	case East:
		dx, dy = 0, 0
	case South:
		dx, dy = 0, 0
	case West:
		dx, dy = -1, 0
	}
	x += dx * offset
	y += dy * offset

	return float64(x), float64(y)
}

func getSnackNeckScale(d, offset int) (float64, float64) {
	switch d {
	case North:
		return 1, float64(CellSize+offset) / CellSize
	case East:
		return float64(CellSize+offset) / CellSize, 1
	case South:
		return 1, float64(CellSize+offset) / CellSize
	case West:
		return float64(CellSize+offset) / CellSize, 1
	}

	return 2, 2
}

func getSnakeNeckScale(d, offset int) (float64, float64) {
	if d == East || d == West {
		return float64(CellSize+offset) / CellSize, 1
	} else if d == North || d == South {
		return 1, float64(CellSize+offset) / CellSize
	}

	return 2, 2
}

func getSnakeGeoM(index, d, offset int) (float64, float64) {
	x := index % Width * CellSize
	y := index / Width * CellSize
	return float64(x), float64(y)
}

func mkSnakeHeadTileImage(d int) *ebiten.Image {
	tileImage := ebiten.NewImage(CellSize, CellSize)
	tileImage.Fill(SnakeColor)
	tileImage.Fill(color.RGBA{252, 220, 77, 255})

	return tileImage
}

func snakeEat(s Snake, index int) Snake {
	return append(Snake{index}, s...)
}

func snakeMove(s Snake, index int) Snake {
	return append(Snake{index}, s[:s.len()-1]...)
}

func nextIndex(index, d int) int {
	switch d {
	case North:
		return index - Width
	case East:
		return index + 1
	case South:
		return index + Width
	case West:
		return index - 1
	}

	log.Panicf("Failed, nextIndex(%d, %d)", index, d)
	return -1
}
