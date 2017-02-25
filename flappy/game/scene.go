package game

import (
	"image/color"
	"log"

	"engo.io/ecs"
	"engo.io/engo"
	"engo.io/engo/common"
)

var (
	blue   = color.RGBA{0x54, 0xc0, 0xc9, 0xff}
	green  = color.RGBA{0x86, 0xa9, 0x4b, 0xff}
	orange = color.RGBA{0xf7, 0xb6, 0x43, 0xff}
)

type Bird struct {
	ecs.BasicEntity
	common.CollisionComponent
	common.RenderComponent
	common.SpaceComponent
}

type Scene struct{}

func (*Scene) Preload() {
	if err := engo.Files.Load("bird.png", "pipe_down.png", "pipe_up.png"); err != nil {
		log.Println(err)
	}
}

func (*Scene) Setup(world *ecs.World) {
	engo.Input.RegisterButton("Exit", engo.Q)
	engo.Input.RegisterButton("Space", engo.Space)

	common.SetBackground(blue)

	world.AddSystem(&common.RenderSystem{})
	world.AddSystem(&common.CollisionSystem{})
	world.AddSystem(&FallingSystem{})
	world.AddSystem(&PipeSpawnSystem{})
	world.AddSystem(&EnemySystem{})

	texture, err := common.LoadedSprite("bird.png")
	if err != nil {
		log.Println(err)
	}

	// Create an entity
	bird := Bird{BasicEntity: ecs.NewBasic()}

	// Initialize the components, set scale to 2x
	bird.RenderComponent = common.RenderComponent{
		Drawable: texture,
		Scale:    engo.Point{2, 2},
	}

	bird.SpaceComponent = common.SpaceComponent{
		Position: engo.Point{engo.WindowWidth() / 5, engo.WindowHeight() / 2},
		Width:    texture.Width() * bird.RenderComponent.Scale.X,
		Height:   texture.Height() * bird.RenderComponent.Scale.Y,
	}

	bird.CollisionComponent = common.CollisionComponent{
		Solid: true,
		Main:  true,
	}

	// Add it to appropriate systems
	for _, system := range world.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			sys.Add(&bird.BasicEntity, &bird.RenderComponent, &bird.SpaceComponent)
		case *common.CollisionSystem:
			sys.Add(&bird.BasicEntity, &bird.CollisionComponent, &bird.SpaceComponent)
		case *FallingSystem:
			sys.Add(&bird.BasicEntity, &bird.SpaceComponent)
		}
	}
}

func (*Scene) Type() string { return "Scene" }
