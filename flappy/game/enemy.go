package game

import (
	"fmt"

	"engo.io/ecs"
	"engo.io/engo"
	"engo.io/engo/common"
)

var score int

type enemyEntity struct {
	*ecs.BasicEntity
	*common.SpaceComponent
}

type EnemySystem struct {
	entities []enemyEntity
}

func (es *EnemySystem) Add(basic *ecs.BasicEntity, space *common.SpaceComponent) {
	es.entities = append(es.entities, enemyEntity{basic, space})
}

func (es *EnemySystem) Remove(basic ecs.BasicEntity) {
	delete := -1
	for index, e := range es.entities {
		if e.BasicEntity.ID() == basic.ID() {
			delete = index
			break
		}
	}
	if delete >= 0 {
		es.entities = append(es.entities[:delete], es.entities[delete+1:]...)
	}
}

func (es *EnemySystem) Update(dt float32) {
	engo.Mailbox.Listen("CollisionMessage", func(message engo.Message) {
		_, isCollision := message.(common.CollisionMessage)

		if isCollision {
			engo.SetSceneByName("GameOver", false)
		}
	})

	speed := 360 * dt

	for _, e := range es.entities {
		e.SpaceComponent.Position.X -= speed

		if int(e.SpaceComponent.Position.X) == 0 {
			score++
			fmt.Println("Score:", score)
		}
	}
}
