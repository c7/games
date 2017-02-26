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
		case *PipeSystem:
			sys.Add(&pipe.BasicEntity, &pipe.SpaceComponent)
		}
	}
}

var score int

type pipeEntity struct {
	*ecs.BasicEntity
	*common.SpaceComponent
}

type PipeSystem struct {
	entities []pipeEntity
}

func (es *PipeSystem) Add(basic *ecs.BasicEntity, space *common.SpaceComponent) {
	es.entities = append(es.entities, pipeEntity{basic, space})
}

func (s *PipeSystem) Remove(basic ecs.BasicEntity) {
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

func (s *PipeSystem) Update(dt float32) {
	engo.Mailbox.Listen("CollisionMessage", func(message engo.Message) {
		_, isCollision := message.(common.CollisionMessage)

		if isCollision {
			engo.SetSceneByName("GameOver", false)
		}
	})

	speed := 360 * dt

	for _, e := range s.entities {
		e.SpaceComponent.Position.X -= speed

		if int(e.SpaceComponent.Position.X) == 0 {
			score++
			fmt.Println("Score:", score)
		}
	}
}
