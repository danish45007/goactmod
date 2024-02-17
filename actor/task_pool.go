package actor

import (
	"sync"

	"github.com/danish45007/goactmod/entities"
)

type TaskActorPool struct {
	actorPoolLock *sync.Mutex
	actorPool     []entities.Actor
	wg            *sync.WaitGroup
}

func CreateTaskActorPool(wg *sync.WaitGroup) *TaskActorPool {
	return &TaskActorPool{
		actorPoolLock: &sync.Mutex{},
		actorPool:     make([]entities.Actor, 0),
		wg:            wg,
	}
}
