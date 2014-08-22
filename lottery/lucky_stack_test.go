package lottery

import (
	"math/rand"
	"testing"
)

// testing lucky stack logic

func TestNextBytesWithInit(t *testing.T) {
	// to have the same sequence each time
	random = rand.New(rand.NewSource(99))

	s := new(stack)

	first, second := s.nextBytes()

	if len(s.luckyBytes) == 0 {
		t.Log("stack must be initilized")
		t.FailNow()
	}
	if s.currentElement != 1 {
		t.Log("Current elemement must be equal to 1, after first call")
		t.FailNow()
	}
	// run one cycle on array:
	for i := 1; i < maxSize; i++ {
		pos := 2 * i
		rFirst, rSecond := s.luckyBytes[pos], s.luckyBytes[pos+1]
		first, second = s.nextBytes()

		if first != rFirst {
			t.Logf("First element is not equal ! first = %d, luckyBytes[%d] = %d", first, pos, rFirst)
			t.FailNow()
		}
		if second != rSecond {
			t.Logf("Second element is not equal ! second = %d, luckyBytes[%d] = %d", second, pos+1, rSecond)
			t.FailNow()
		}

		if s.currentElement != i+1 {
			t.Logf("Current elemement = %d, must be %d", s.currentElement, i+1)
			t.FailNow()
		}
	}
}
