package actor_engine

import (
	"sync"

	"github.com/danish45007/goactmod/actor"
	"github.com/danish45007/goactmod/entities"
	"github.com/danish45007/goactmod/tracker"
	"github.com/ian-kent/go-log/log"
)

type ActorEngine struct {
	name     string
	assigner entities.Actor
	wg       *sync.WaitGroup
	tracker  *tracker.Tracker
}

func (engine *ActorEngine) SubmitTask(task entities.Task) error {
	return engine.assigner.AddTask(task)
}

func (engine *ActorEngine) Run() {
	log.Debug("actor engine %s is started \n", engine.name)
	engine.assigner.Start()
}

func (engine *ActorEngine) Shutdown(wg *sync.WaitGroup) {
	defer wg.Done()
	// close the assigner actor
	engine.assigner.Stop()

	// wait for the actor engine to close
	engine.wg.Wait()

	// close the tracker
	engine.tracker.Shutdown()
	log.Debug("actor engine %s is closed \n", engine.name)
}

// CreateActorEngine creates a new actor engine
func CreateActorEngine(name string, config *actor.Config) *ActorEngine {
	wg := &sync.WaitGroup{}
	actorPool := actor.CreateTaskActorPool(wg)
	tracker := tracker.CreateTracker(name)
	assignerActor := actor.CreateAssignerActor(name, actorPool, tracker, config)
	engine := &ActorEngine{
		name:     name,
		wg:       wg,
		tracker:  tracker,
		assigner: assignerActor,
	}
	go engine.Run()
	return engine
}
