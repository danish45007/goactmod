package actor

import (
	"errors"
	"sync"

	"github.com/danish45007/goactmod/entities"
	"github.com/danish45007/goactmod/tracker"
	"github.com/ian-kent/go-log/log"
)

const TaskActorQueueSize = 10

type TaskActor struct {
	id            int
	closingSignal chan bool
	wg            *sync.WaitGroup
	tasks         chan entities.Task
	tracker       *tracker.Tracker
}

// task_actor implements the Actor interface

// AddTask adds a task to the task actor queue
func (ta *TaskActor) AddTask(task entities.Task) error {
	// check if the actor task queue is full
	if len(ta.tasks) == TaskActorQueueSize {
		return errors.New("task actor queue is full")
	}
	// add the task to the actor task queue
	ta.tasks <- task
	return nil
}

// Start spawns a new task actor and executes the task
func (ta *TaskActor) Start() {
	defer ta.wg.Done()
	// increment the wait group counter
	ta.wg.Add(1)
	log.Debug("task actor %d is spawned", ta.id)
	for task := range ta.tasks {
		// execute the task
		task.Execute()
		// send the completed task to the tracker channel
		ta.tracker.GetTrackChannel() <- tracker.CreateCounterTrack(tracker.Task, tracker.Completed)
	}
	log.Debug("task actor %d is closed", ta.id)
	// close the closingSignal channel
	ta.closingSignal <- true
}

// Stop closes the task actor
func (ta *TaskActor) Stop() {
	// close the closingSignal channel
	close(ta.tasks)
	// wait for the task actor to close
	<-ta.closingSignal
}

// CreateTaskActor creates a new task actor
func CreateTaskActor(wg *sync.WaitGroup, id int, tracker *tracker.Tracker) entities.Actor {
	actor := &TaskActor{
		id:            id,
		closingSignal: make(chan bool),
		wg:            wg,
		tasks:         make(chan entities.Task, TaskActorQueueSize),
		tracker:       tracker,
	}
	go actor.Start()
	return actor
}
