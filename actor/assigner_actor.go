package actor

import (
	"github.com/danish45007/goactmod/entities"
	"github.com/danish45007/goactmod/tracker"
)

const AssignerQueueSize = 10e2

type AssignerActor struct {
	name          string
	closeSignal   chan bool
	tasks         chan entities.Task
	assignerIndex int
	tracker       *tracker.Tracker
	scaler        *scaler
	*TaskActorPool
	*Config
}

func CreateAssignerActor(pool *TaskActorPool, tracker tracker.Track, config *Config) entities.Actor {
	return &assignerActor{
		name:          "assigner",
		closeSignal:   make(chan bool),
		tasks:         make(chan entities.Task, assignerQueueSize),
		assignerIndex: 0,
		tracker:       tracker,
		scaler:        scaler,
		TaskActorPool: pool,
		Config:        config,
	}
}
