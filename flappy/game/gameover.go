package game

import (
	"log"

	"engo.io/ecs"
	"engo.io/engo"
	"engo.io/engo/common"
)

type Image struct {
	ecs.BasicEntity
	common.RenderComponent
	common.SpaceComponent
}

type GameOverScene struct{}

func (*GameOverScene) Preload() {
	if err := engo.Files.Load("game_over.png"); err != nil {
		log.Println(err)
	}
}

func (*GameOverScene) Setup(world *ecs.World) {
	world.AddSystem(&common.RenderSystem{})
	world.AddSystem(&GameOverSystem{})

	// Retrieve a texture
	texture, err := common.LoadedSprite("game_over.png")
	if err != nil {
		log.Println(err)
	}

	// Create an entity
	gameOver := Image{BasicEntity: ecs.NewBasic()}

	// Initialize the components, set scale to 8x
	gameOver.RenderComponent = common.RenderComponent{
		Drawable: texture,
		Scale:    engo.Point{2, 2},
	}

	gameOver.SpaceComponent = common.SpaceComponent{
		Position: engo.Point{engo.GameWidth()/2 - texture.Width(), engo.GameHeight() / 3},
		Width:    texture.Width() * gameOver.RenderComponent.Scale.X,
		Height:   texture.Height() * gameOver.RenderComponent.Scale.Y,
	}

	// Add it to appropriate systems
	for _, system := range world.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			sys.Add(&gameOver.BasicEntity, &gameOver.RenderComponent, &gameOver.SpaceComponent)
		case *GameOverSystem:
			sys.Add(&gameOver.BasicEntity, &gameOver.RenderComponent)
		}
	}
}

func (*GameOverScene) Type() string { return "GameOver" }

func (*GameOverScene) Show() {
	common.SetBackground(orange)
}

func (*GameOverScene) Hide() {
	common.SetBackground(blue)
	score = 0
}

type hideEntity struct {
	*ecs.BasicEntity
	*common.RenderComponent
}

type GameOverSystem struct {
	entities []hideEntity
}

func (s *GameOverSystem) Add(basic *ecs.BasicEntity, render *common.RenderComponent) {
	s.entities = append(s.entities, hideEntity{basic, render})
}

func (s *GameOverSystem) Update(dt float32) {
	for _, e := range s.entities {
		e.RenderComponent.Hidden = false
	}

	if engo.Input.Button("Exit").JustReleased() {
		engo.Exit()
	}

	if engo.Input.Button("Space").JustReleased() {
		engo.SetSceneByName("Scene", true)
	}
}

func (s *GameOverSystem) Remove(basic ecs.BasicEntity) {
	delete := -1
	for index, e := range s.entities {
		if e.BasicEntity.ID() == basic.ID() {
			delete = index
			break
		}
	}
	if delete >= 0 {
		s.entities = append(s.entities[:delete], s.entities[delete+1:]...)
	}
}
