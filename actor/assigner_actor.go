package actor

import (
	"errors"

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
	scaler        *AutoScaler
	*TaskActorPool
	*Config
}

func CreateAssignerActor(name string, pool *TaskActorPool, tracker *tracker.Tracker, config *Config) entities.Actor {
	return &AssignerActor{
		name:          name,
		closeSignal:   make(chan bool),
		tasks:         make(chan entities.Task, AssignerQueueSize),
		assignerIndex: 0,
		TaskActorPool: pool,
		tracker:       tracker,
		Config:        config,
	}
}

func (assigner *AssignerActor) QueueSize() int {
	return len(assigner.tasks)
}

func (assigner *AssignerActor) AddTask(task entities.Task) error {
	// check if the actor task queue is full
	if len(assigner.tasks) >= AssignerQueueSize {
		assigner.tracker.GetTrackChannel() <- tracker.CreateCounterTrack(tracker.Task, tracker.Rejected)
		return errors.New("assigner actor queue is full")
	}
	// add the task to the actor task queue
	assigner.tasks <- task
	assigner.tracker.GetTrackChannel() <- tracker.CreateCounterTrack(tracker.Task, tracker.Submitted)
	return nil
}

func (assigner *AssignerActor) Start() {
	poolStarted := make(chan bool)
	assigner.scaler = GetAutoScaler(assigner, poolStarted)
	// start the auto scaler
	<-poolStarted
	for task := range assigner.tasks {
		for {
			assigner.actorPoolLock.Lock()
			assigner.assignerIndex = assigner.assignerIndex % len(assigner.actorPool)
			actor := assigner.actorPool[assigner.assignerIndex]
			assigner.assignerIndex++
			assigner.actorPoolLock.Unlock()
			err := actor.AddTask(task)
			if err == nil {
				break
			}
		}
	}
	assigner.closeSignal <- true

}

func (assigner *AssignerActor) Stop() {
	close(assigner.tasks)
	<-assigner.closeSignal
	assigner.scaler.Stop()
}
