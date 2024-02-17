package tracker

import (
	"sync"
	"time"

	"github.com/ian-kent/go-log/log"
)

const QueueSize = 10e4

type TrackScope string

// TrackScope is a type for the scope of the metric to be tracked
const (
	ActorEngine TrackScope = "actor-engine"
	Task        TrackScope = "task"
)

type TrackMetric string

// TrackMetric is a type for the metric to be tracked
const (
	Submitted    TrackMetric = "submitted"
	Completed    TrackMetric = "completed"
	Rejected     TrackMetric = "rejected"
	ActiveActors TrackMetric = "active-actors"
)

// Track is a struct that holds the scope, metric and the new actors count
type Track struct {
	Scope          TrackScope
	Metric         TrackMetric
	NewActorsCount int
}

// Tracker is a struct that holds the metrics and the channels to send the metrics
type Tracker struct {
	closingSignal chan bool
	sysname       string
	tracker       chan Track
	metricLock    *sync.RWMutex
	metrics       map[TrackScope]map[TrackMetric]int
}

// collect metrics from the tracker
func (t *Tracker) collectMetrics() {
	for track := range t.tracker {
		// if the metric is not present, create a new map
		if t.metrics[track.Scope] == nil {
			t.metrics[track.Scope] = make(map[TrackMetric]int)
		}
		t.metricLock.Lock()
		// update the metric with the new actors count
		t.metrics[track.Scope][track.Metric] += track.NewActorsCount
		t.metricLock.Unlock()
	}
	// close the close_signal channel
	t.closingSignal <- true
}

// shutdown the tracker
func (t *Tracker) Shutdown() {
	// close the tracker channel
	close(t.tracker)
	// wait for the collectMetrics goroutine to finish
	<-t.closingSignal
}

func (t *Tracker) GetTrackChannel() chan Track {
	return t.tracker
}

func (t *Tracker) forEverPrintMetricsLoop() {
	for {
		// print the metrics with a delay of 1 seconds
		time.Sleep(1 * time.Second)
		t.printMetrics()
	}
}

func (t *Tracker) printMetrics() {
	// acquire a read lock
	t.metricLock.RLock()
	// log metrics
	log.Debug("system: %s, metrics: %v", t.sysname, t.metrics)
	// release the read lock
	t.metricLock.RUnlock()
}

// CreateTracker creates a new tracker
func CreateTracker(sysname string) *Tracker {
	t := &Tracker{
		closingSignal: make(chan bool),
		sysname:       sysname,
		tracker:       make(chan Track, QueueSize),
		metricLock:    &sync.RWMutex{},
		metrics:       make(map[TrackScope]map[TrackMetric]int),
	}
	// start the collectMetrics goroutine
	go t.collectMetrics()
	// start the forEverPrintMetricsLoop goroutine
	go t.forEverPrintMetricsLoop()
	return t
}

// CreateCounterTrack creates a default track with the scope and metric and 1 new actors count
func CreateCounterTrack(scope TrackScope, metric TrackMetric) Track {
	return Track{
		Scope:          scope,
		Metric:         metric,
		NewActorsCount: 1,
	}
}

// CreateTrack creates a new track with the scope, metric and the new actors count
func CreateTrack(scope TrackScope, metric TrackMetric, delta int) Track {
	return Track{
		Scope:          scope,
		Metric:         metric,
		NewActorsCount: delta,
	}
}
