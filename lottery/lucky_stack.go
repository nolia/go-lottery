package lottery

import (
	"fmt"
	"math/rand"
	"time"
)

const (
	// maximum size of lucky pairs stack
	maxSize                = 100
	defaultDiscardDuration = time.Duration(10) * time.Second
)

var random *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

type LuckyStack interface {
	Check(first, second byte) bool
}

// creates and initilizes lucky stack
func NewLuckyStack() LuckyStack {
	s := new(stack)

	s.output = make(chan pair)
	go s.fill()
	s.startDiscardingTimer(defaultDiscardDuration)

	return s
}

type pair struct {
	first, second byte
}

type stack struct {
	luckyBytes     []byte
	currentElement int
	output         chan pair

	resetDuration time.Duration
	discardTimer  *time.Timer
}

// main logic of lucky stack:
func (s *stack) nextBytes() (first, second byte) {
	// first time or empty - init
	if len(s.luckyBytes) == 0 {
		s.luckyBytes = make([]byte, 2*maxSize)
	}
	// if reached last - go to beginning
	if s.currentElement > maxSize {
		s.currentElement = 0
	}
	// before getting values out of array - invalidate it
	if s.currentElement == 0 {
		for i := 0; i < 2*maxSize; i++ {
			s.luckyBytes[i] = byte(random.Int31n(256))
		}
	}
	pos := 2 * s.currentElement
	first, second = s.luckyBytes[pos], s.luckyBytes[pos+1]
	s.currentElement++
	return
}

func (s *stack) fill() {
	for {
		first, second := s.nextBytes()
		s.output <- pair{first, second}
	}
}

func (s *stack) startDiscardingTimer(d time.Duration) {
	s.resetDuration = d
	s.discardTimer = time.AfterFunc(d, func() {
		// poping up top pair and discarding bytes
		p := <-s.output
		fmt.Printf("No requests. Discarding pair %v\n", p)
		// placing timer again
		s.discardTimer.Reset(d)
	})
}

func (s *stack) Check(first, second byte) bool {
	p := <-s.output
	fmt.Printf("Checking %v, requested: [%x, %x]\n", p, first, second)

	// if discarding is enabled
	if s.discardTimer != nil {
		s.discardTimer.Reset(s.resetDuration)
	}

	return p.first == first && p.second == second
}
