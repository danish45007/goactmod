package actor

import (
	"time"

	"github.com/danish45007/goactmod/entities"
	"github.com/danish45007/goactmod/tracker"
	"github.com/ian-kent/go-log/log"
)

/*
auto scaler is part task assigner actor
it scales the task actor pool based on the number of tasks in the queue
*/

type AutoScaler struct {
	*AssignerActor
	lastActorId   int
	ClosingSignal chan bool
	ClosedSignal  chan bool
}

// NewAutoScaler creates a new auto scaler
func GetAutoScaler(assignerActor *AssignerActor) *AutoScaler {
	scaler := &AutoScaler{
		AssignerActor: assignerActor,
		lastActorId:   0,
		ClosingSignal: make(chan bool),
		ClosedSignal:  make(chan bool),
	}
	return scaler
}

// Run starts the auto scaler
func (as *AutoScaler) Run(poolStarted chan bool) {
	log.Debug("auto scaler is started")
	// provision with the minimum number of actors
	as.provisionActor(as.Config.MinActors)
	// send the pool started signal
	poolStarted <- true
	processCompleted := false
	for !processCompleted {
		select {
		case <-as.ClosingSignal:
			processCompleted = true
		// auto scale by provisioning/deprovisioning actors
		case <-time.After(time.Duration(as.Config.ScalingInterval) * time.Microsecond):
			// in case provision
			// if the queue size is greater than the max queue size and the actor pool size is less than the max actors
			if (as.QueueSize() > as.Config.UpScaleQueueSize) && (len(as.actorPool) < as.Config.MaxActors) {
				as.provisionActor(as.Config.UpScaleFactor)
			} else if (as.QueueSize() < as.Config.DownScaleQueueSize) && (len(as.actorPool) > as.Config.MinActors) {
				as.deprovisionActor(as.Config.DownScaleFactor)
			}
		}
		// deprovision the remaining actors in case of a stop signal
		log.Debug("auto scaler is stopped")
		as.deprovisionActor(len(as.actorPool))
		// close the auto scaler
		as.ClosedSignal <- true
	}
	//

}

// stop the auto scaler
func (as *AutoScaler) Stop() {
	as.ClosingSignal <- true
	// wait for the auto scaler to close
	<-as.ClosedSignal
}

// deprovisionActor removes (delta) actors from the actor pool
// and stops the actors
func (as *AutoScaler) deprovisionActor(delta int) {
	log.Debug("deprovisioning actors in %s by %d", as.name, delta)
	as.actorPoolLock.Lock()
	// if delta is greater than the actor pool size, set delta to the actor pool size
	if delta > len(as.actorPool) {
		delta = len(as.actorPool)

	}
	deprovisionedActors := as.actorPool[:delta]
	// update the tracker with the new actors count
	as.tracker.GetTrackChannel() <- tracker.CreateTrack(tracker.ActorEngine, tracker.ActiveActors, -delta)
	// remove the actors from the actor pool
	as.actorPool = as.actorPool[delta:]
	// update the tracker with the new actors count
	as.tracker.GetTrackChannel() <- tracker.CreateTrack(tracker.ActorEngine, tracker.ActiveActors, -delta)
	// stop the actors
	for _, deprovisionActor := range deprovisionedActors {
		deprovisionActor.Stop()
	}
	as.actorPoolLock.Unlock()
}

// provisionActor add new (delta) actors to the actor pool
func (as *AutoScaler) provisionActor(delta int) []entities.Actor {
	log.Debug("provisioning actors in %s by %d", as.name, delta)
	actors := make([]entities.Actor, delta)

	// populate the actors with new task actors
	for i := 0; i < delta; i++ {
		actors[i] = CreateTaskActor(as.wg, i+as.lastActorId, as.tracker)
	}
	as.lastActorId += delta
	// update the tracker with the new actors count
	// add the new actors to the actor pool
	as.actorPoolLock.Lock()
	as.tracker.GetTrackChannel() <- tracker.CreateTrack(tracker.ActorEngine, tracker.ActiveActors, delta)
	as.actorPool = append(as.actorPool, actors...)
	as.actorPoolLock.Unlock()
	return actors
}
