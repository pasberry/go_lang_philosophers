//Implement the dining philosopher’s problem with the following constraints/modifications.

//1. There should be 5 philosophers sharing chopsticks, with one chopstick between each adjacent pair of philosophers.
//2. Each philosopher should eat only 3 times (not in an infinite loop as we did in lecture)
//3. The philosophers pick up the chopsticks in any order, not lowest-numbered first (which we did in lecture).
//4. In order to eat, a philosopher must get permission from a host which executes in its own goroutine.
//5. The host allows no more than 2 philosophers to eat concurrently.
//6. Each philosopher is numbered, 1 through 5.
//7. When a philosopher starts eating (after it has obtained necessary locks) it prints “starting to eat <number>” on a
//   line by itself, where <number> is the number of the philosopher.
//8. When a philosopher finishes eating (before it has released its locks) it prints “finishing eating <number>” on a
//   line by itself, where <number> is the number of the philosopher.
package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type ChopStick struct {
	mutex sync.Mutex
}

type Host struct {
	firstEatInvite, secondEatInvite EatInvitation
}

type EatInvitation struct {
	mutex sync.Mutex
}

func (h *Host) getInvite() EatInvitation {
	t := rand.Intn(10)

	var invitation EatInvitation

	if t % 2 == 0 {
		invitation =  h.firstEatInvite
	} else {
		invitation =  h.secondEatInvite
	}

	return invitation
}

type Philosopher struct {
	id, requiredMeals, mealsHad int
	leftChopstick, rightChopstick ChopStick
	host Host
}

func (p Philosopher) think() {
	t := rand.Intn(10)
	time.Sleep(time.Duration(t) * time.Second)
}

func (p *Philosopher) isFed() bool{
	return p.mealsHad == p.requiredMeals
}

func (p *Philosopher) finishedEating() {
	p.mealsHad++
}

func (p *Philosopher) eat() {

	p.think()

	invitation := p.host.getInvite()
	for !p.isFed() {

		invitation.mutex.Lock()
		p.think()

		p.leftChopstick.mutex.Lock()
		p.rightChopstick.mutex.Lock()

		p.think()
		fmt.Println("starting to eat", p.id)
		p.think()
		fmt.Println("finishing eating", p.id)

		p.leftChopstick.mutex.Unlock()
		p.rightChopstick.mutex.Unlock()
		invitation.mutex.Unlock()

		p.finishedEating()

	}

}

func main() {

	host := Host{EatInvitation{},EatInvitation{}}

	chopsticks := [5]ChopStick{{},{},{},{},{}}

	philosophers := []Philosopher{
		{id:1,mealsHad: 0,requiredMeals:3, leftChopstick:chopsticks[0],rightChopstick:chopsticks[1],host:host},
		{id:2,mealsHad: 0,requiredMeals:3, leftChopstick:chopsticks[1],rightChopstick:chopsticks[2],host:host},
		{id:3,mealsHad: 0,requiredMeals:3, leftChopstick:chopsticks[2],rightChopstick:chopsticks[3],host:host},
		{id:4,mealsHad: 0,requiredMeals:3, leftChopstick:chopsticks[3],rightChopstick:chopsticks[4],host:host},
		{id:5,mealsHad: 0,requiredMeals:3, leftChopstick:chopsticks[0],rightChopstick:chopsticks[4],host:host},
	}
	var waitGroup sync.WaitGroup
	waitGroup.Add(5)

	go haveDinner(&waitGroup, &philosophers[0])
	go haveDinner(&waitGroup, &philosophers[1])
	go haveDinner(&waitGroup, &philosophers[2])
	go haveDinner(&waitGroup, &philosophers[3])
	go haveDinner(&waitGroup, &philosophers[4])
	
	waitGroup.Wait()
}

func haveDinner(waitGroup *sync.WaitGroup, philospher *Philosopher) {

	defer waitGroup.Done()
	philospher.eat()
}