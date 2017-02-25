package game

import (
	"log"

	"engo.io/ecs"
	"engo.io/engo"
	"engo.io/engo/common"
)

type GameOverScene struct{}

func (*GameOverScene) Preload() {
	if err := engo.Files.Load("game_over.png"); err != nil {
		log.Println(err)
	}
}

func (*GameOverScene) Setup(world *ecs.World) {
	world.AddSystem(&GameOverSystem{})
}

func (*GameOverScene) Type() string { return "GameOver" }

func (*GameOverScene) Show() {
	common.SetBackground(orange)
}

func (*GameOverScene) Hide() {
	common.SetBackground(blue)
	score = 0
}

type GameOverSystem struct{}

func (s *GameOverSystem) Update(dt float32) {
	if engo.Input.Button("Exit").JustReleased() {
		engo.Exit()
	}

	if engo.Input.Button("Space").JustReleased() {
		engo.SetSceneByName("Scene", true)
	}
}

func (s *GameOverSystem) Remove(basic ecs.BasicEntity) {}
