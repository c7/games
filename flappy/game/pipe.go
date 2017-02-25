package game

import (
	"fmt"
	"log"
	"math/rand"

	"engo.io/ecs"
	"engo.io/engo"
	"engo.io/engo/common"
)

type Pipe struct {
	ecs.BasicEntity
	common.CollisionComponent
	common.RenderComponent
	common.SpaceComponent
}

type PipeSpawnSystem struct {
	world *ecs.World
}

func (pss *PipeSpawnSystem) New(w *ecs.World) {
	pss.world = w
}

func (*PipeSpawnSystem) Remove(ecs.BasicEntity) {
	fmt.Println("remove")
}

func (pss *PipeSpawnSystem) Update(dt float32) {
	// 4% change of spawning a pipe each frame
	if rand.Float32() < .96 {
		return
	}

	position := engo.Point{
		X: engo.GameWidth()*2 + rand.Float32()*64,
		Y: rand.Float32()*engo.GameHeight() + engo.GameHeight()/1.5,
	}

	NewPipe(pss.world, position, "pipe_up.png")

	if position.Y < engo.GameHeight() {
		NewPipe(pss.world, engo.Point{position.X, rand.Float32()*160 - 280}, "pipe_down.png")
	}
}

func NewPipe(world *ecs.World, position engo.Point, spriteName string) {
	texture, err := common.LoadedSprite(spriteName)
	if err != nil {
		log.Println(err)
	}

	pipe := Pipe{BasicEntity: ecs.NewBasic()}
	pipe.RenderComponent = common.RenderComponent{
		Drawable: texture,
		Scale:    engo.Point{2, 2},
	}
	pipe.SpaceComponent = common.SpaceComponent{
		Position: position,
		Width:    texture.Width() * pipe.RenderComponent.Scale.X,
		Height:   texture.Height() * pipe.RenderComponent.Scale.Y,
	}
	pipe.CollisionComponent = common.CollisionComponent{Solid: true}

	for _, system := range world.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			sys.Add(&pipe.BasicEntity, &pipe.RenderComponent, &pipe.SpaceComponent)
		case *common.CollisionSystem:
			sys.Add(&pipe.BasicEntity, &pipe.CollisionComponent, &pipe.SpaceComponent)
		case *EnemySystem:
			sys.Add(&pipe.BasicEntity, &pipe.SpaceComponent)
		}
	}
}
