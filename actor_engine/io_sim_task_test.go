package actor_engine_test

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/danish45007/goactmod/actor"
	"github.com/danish45007/goactmod/actor_engine"
)

// PrimeNumberTask implements entities.Task interface to generate and print prime numbers
type PrimeNumberTask struct {
}

func (t *PrimeNumberTask) Execute() {
	const n = 1000000 // Generate 10 million prime numbers
	count := 0
	num := 2
	for count < n {
		if isPrime(num) {
			fmt.Println(num) // Print the prime number
			count++
		}
		num++
	}
}

// isPrime checks if a number is prime
func isPrime(num int) bool {
	if num <= 1 {
		return false
	}
	for i := 2; i*i <= num; i++ {
		if num%i == 0 {
			return false
		}
	}
	return true
}

func TestActorEngine_ExecutePrimeNumberTask(t *testing.T) {

	// Create the ActorEngine with the real assigner actor
	engine := actor_engine.CreateActorEngine("PrimeNumberEngine", &actor.Config{})

	// Submit the PrimeNumberTask to the ActorEngine
	primeNumberTask := &PrimeNumberTask{}
	err := engine.SubmitTask(primeNumberTask)
	if err != nil {
		t.Errorf("SubmitTask() failed: %v", err)
	}

	// Simulate delay for task execution
	time.Sleep(1000 * time.Microsecond)

	// Stop the ActorEngine
	engine.Shutdown(&sync.WaitGroup{})
}
