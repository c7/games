package main

import (
	"engo.io/engo"

	"github.com/c7/games/flappy/game"
)

func main() {
	opts := engo.RunOptions{
		Title:    "Flappy",
		Width:    288,
		Height:   512,
		VSync:    true,
		FPSLimit: 60,
	}

	scene := &game.Scene{}
	gameOverScene := &game.GameOverScene{}

	engo.RegisterScene(gameOverScene)

	engo.Run(opts, scene)
}
