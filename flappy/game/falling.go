package game

import (
	"engo.io/ecs"
	"engo.io/engo"
	"engo.io/engo/common"
)

type fallingEntity struct {
	*ecs.BasicEntity
	*common.SpaceComponent
}

type FallingSystem struct {
	entities []fallingEntity
}

func (f *FallingSystem) Add(basic *ecs.BasicEntity, space *common.SpaceComponent) {
	f.entities = append(f.entities, fallingEntity{basic, space})
}

func (f *FallingSystem) Remove(basic ecs.BasicEntity) {
	delete := -1
	for index, e := range f.entities {
		if e.BasicEntity.ID() == basic.ID() {
			delete = index
			break
		}
	}
	if delete >= 0 {
		f.entities = append(f.entities[:delete], f.entities[delete+1:]...)
	}
}

func (f *FallingSystem) Update(dt float32) {
	if engo.Input.Button("Exit").JustReleased() {
		engo.Exit()
	}

	speed := 256 * dt

	for _, e := range f.entities {
		if engo.Input.Button("Space").JustPressed() {
			e.SpaceComponent.Position.Y -= speed + (e.SpaceComponent.Height * 2.4)
		}

		y := e.SpaceComponent.Position.Y

		if y < 0 || y > engo.WindowHeight()-e.SpaceComponent.Height {
			engo.Mailbox.Dispatch(common.CollisionMessage{})
		}

		e.SpaceComponent.Position.Y += speed
	}
}
