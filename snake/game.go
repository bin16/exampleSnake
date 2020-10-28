package snake

import (
	"errors"
	"fmt"
	"image/color"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	GameMainTitle = iota
	GamePlaying
	GameEnd
)

var (
	GameOverEatSelf  = errors.New("Eat Self")
	GameOverHitWall  = errors.New("Hit Wall")
	GameOverNoReason = errors.New("No Reason")
	GameOverYouWin   = errors.New("OK, You Win")

	grassTileImage, foodTileImage, wallTileImage *ebiten.Image

	GrassColor = color.RGBA{129, 174, 157, 255}
	FoodColor  = color.RGBA{251, 159, 137, 255}
	SnakeColor = color.RGBA{103, 89, 122, 255}
	WallColor  = color.RGBA{55, 55, 55, 255}
)

func init() {
	grassTileImage = ebiten.NewImage(CellSize, CellSize)
	grassTileImage.Fill(GrassColor)

	foodTileImage = ebiten.NewImage(CellSize-CellMargin*2, CellSize-CellMargin*2)
	foodTileImage.Fill(FoodColor)

	wallTileImage = ebiten.NewImage(CellSize, CellSize)
	wallTileImage.Fill(WallColor)
}

type Game struct {
	status   int
	width    int
	height   int
	cellSize int

	foods []int
	walls []int

	snake     *Snake
	speed     int // N pixels in frame
	offset    int
	tasks     []int
	keyDir    int
	lastError error
}

func (g *Game) init() {
	g.foods = []int{}
	g.makeFood()
	g.makeFood()
	g.makeFood()

	ci := Height/2*Width + Width/2

	*g.snake = Snake{ci, ci + 1}
	g.keyDir = -1
}

func (g *Game) takeFood(foodIndex int) {
	nFoods := []int{}
	for _, index := range g.foods {
		if index != foodIndex {
			nFoods = append(nFoods, index)
		}
	}

	g.foods = nFoods
}

func (g *Game) makeFood() error {
	allowList := []int{}
	for i := 0; i < Width*Height; i++ {
		if inIntList(i, g.foods) {
			continue
		} else if inIntList(i, []int(*g.snake)) {
			continue
		} else if isWall(i) {
			continue
		}

		allowList = append(allowList, i)
	}

	if len(allowList) == 0 {
		return GameOverYouWin
	}

	nIndex := rand.Intn(len(allowList))
	g.foods = append(g.foods, allowList[nIndex])

	return nil
}

func isWall(i int) bool {
	return i/Width == 0 || i/Width == Height-1 || i%Width == 0 || i%Width == Width-1
}

func inIntList(item int, list []int) bool {
	for _, a := range list {
		if a == item {
			return true
		}
	}

	return false
}

func notInIntList(item int, list []int) bool {
	for _, a := range list {
		if a == item {
			return false
		}
	}

	return true
}

func (g *Game) moveTo(index int) error {
	s := g.snake
	if isWall(index) {
		return GameOverHitWall
	} else if inIntList(index, []int(*g.snake)) {
		return GameOverEatSelf
	} else if inIntList(index, g.foods) {
		*g.snake = append(Snake{index}, (*s)...)
		g.takeFood(index)
		return g.makeFood()
	} else {
		*g.snake = append(Snake{index}, (*s)[:s.len()-1]...)
	}

	return nil
}

func (g *Game) WindowWidth() int {
	return g.width * g.cellSize
}

func (g *Game) WindowHeight() int {
	return g.height * g.cellSize
}

func NewGame() *Game {
	s := &Snake{124, 125, 126, 127, 128, 129, 130}
	g := &Game{
		status:   GameMainTitle,
		width:    32,
		height:   24,
		cellSize: 20,
		snake:    s,
		speed:    1,
		offset:   0,
		tasks:    []int{},
		keyDir:   -1,
	}
	g.init()
	return g
}

func (g *Game) Update() error {
	switch g.status {
	case GameMainTitle:
		if inpututil.IsKeyJustReleased(ebiten.KeySpace) {
			g.init()
			g.status = GamePlaying
		}
	case GamePlaying:
		h := (*g.snake)[0]
		if inpututil.IsKeyJustReleased(ebiten.KeyUp) {
			g.keyDir = North
		} else if inpututil.IsKeyJustReleased(ebiten.KeyRight) {
			g.keyDir = East
		} else if inpututil.IsKeyJustReleased(ebiten.KeyDown) {
			g.keyDir = South
		} else if inpututil.IsKeyJustReleased(ebiten.KeyLeft) {
			g.keyDir = West
		}

		g.offset += g.speed
		if g.keyDir != -1 { // with key event
			nextIndex := getNextIndex(h, g.keyDir)
			if err := g.moveTo(nextIndex); isGameOver(err) {
				g.lastError = err
				g.status = GameEnd
				return nil
			}
			g.offset = 0
			g.keyDir = -1
		} else if d := g.offset - CellSize; d >= 0 {
			nextIndex := getNextIndex(h, g.snake.dir())
			if err := g.moveTo(nextIndex); isGameOver(err) {
				g.lastError = err
				g.status = GameEnd
				return nil
			}
			g.offset -= CellSize
		}
	case GameEnd:
		if inpututil.IsKeyJustReleased(ebiten.KeyEscape) {
			g.status = GameMainTitle
		}
	}

	return nil
}

func isGameOver(err error) bool {
	if err == GameOverHitWall || err == GameOverEatSelf {
		return true
	}

	return false
}

func getNextIndex(index, d int) int {
	nextIndex := -1
	switch d {
	case North:
		nextIndex = index - Width
	case East:
		nextIndex = index + 1
	case South:
		nextIndex = index + Width
	case West:
		nextIndex = index - 1
	}

	return nextIndex
}

func (g *Game) Draw(screen *ebiten.Image) {
	switch g.status {
	case GameMainTitle:
		screen.Fill(SnakeColor)
		ebitenutil.DebugPrint(screen, fmt.Sprintf("Welcome to Snake!, [Space] to start."))
	case GamePlaying:
		screen.Fill(GrassColor)

		g.drawWall(screen)
		g.drawFoods(screen)
		g.drawSnake(screen)

		ebitenutil.DebugPrint(screen, fmt.Sprintf("Score: %d", (g.snake.len()-2)*100*g.snake.len()))
	case GameEnd:
		screen.Fill(WallColor)
		ebitenutil.DebugPrint(screen, fmt.Sprintf("GAME OVER - %s\n[Esc] to back to Main Title", g.lastError))
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	w := g.WindowWidth()
	h := g.WindowHeight()
	return w, h
}

func (g *Game) drawWall(screen *ebiten.Image) {
	for i := 0; i < Width*Height; i++ {
		if isWall(i) {
			x := i % Width
			y := i / Width

			op := ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(x*CellSize+CellMargin), float64(y*CellSize+CellMargin))
			screen.DrawImage(wallTileImage, &op)
		}
	}
}

func (g *Game) drawFoods(screen *ebiten.Image) {
	for _, i := range g.foods {
		x := i % Width
		y := i / Width

		op := ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(x*CellSize+CellMargin), float64(y*CellSize+CellMargin))
		screen.DrawImage(foodTileImage, &op)
	}
}

func (g *Game) drawSnake(screen *ebiten.Image) {
	s := *g.snake
	offset := g.offset
	snakeTileImage := ebiten.NewImage(CellSize, CellSize)
	snakeTileImage.Fill(SnakeColor)

	for i := 0; i < s.len(); i++ {
		index := s[i]
		op := &ebiten.DrawImageOptions{}
		if i == 0 {
			op.GeoM.Translate(getSnakeHeadGeoM(index, s.dir(), offset))
			screen.DrawImage(mkSnakeHeadTileImage(s.dir()), op)
		} else if i == s.len()-1 {
			op.GeoM.Translate(getSnakeHeadGeoM(index, s.tailDir(), offset))
			screen.DrawImage(snakeTileImage, op)
		} else if i == 1 {
			op.GeoM.Scale(getSnackNeckScale(s.dir(), offset))
			op.GeoM.Translate(getSnakeNeckGeoM(index, s.dir(), offset))
			screen.DrawImage(snakeTileImage, op)
		} else {
			op.GeoM.Translate(getSnakeGeoM(index, s.dir(), offset))
			screen.DrawImage(snakeTileImage, op)
		}
	}
}
