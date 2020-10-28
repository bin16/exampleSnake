package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/bin16/monster/snake"
)

func main() {
	game := snake.NewGame()

	ebiten.SetWindowTitle("🐍 贪吃蛇 🐍")
	ebiten.SetWindowSize(game.WindowWidth(), game.WindowHeight())
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
