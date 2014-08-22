package lottery

import (
	"math/rand"
	"strconv"
	"testing"
	"time"
)

// always returns false
type failStack int

func (f *failStack) Check(first, second byte) bool {
	return false
}

// always returns true
type winStack int

func (w *winStack) Check(first, second byte) bool {
	return true
}

func TestFail(t *testing.T) {
	lc := new(failStack)
	r := Request{
		Id:     "user1",
		Fee:    100,
		First:  1,
		Second: 1,
	}
	j := NewJackpot(0.0)

	res := Play(lc, j, r)
	t.Logf("result: %v\n", res)
	// should fail
	if res.ResultType != "lost" {
		t.Log("Should have lost, found ", res.ResultType)
		t.FailNow()
	}

	if j.amount != r.Fee {
		t.Logf("Amount is not write, %f", j.amount)
		t.FailNow()
	}
}

func TestBonus(t *testing.T) {
	lc := new(winStack)
	r := Request{
		Id:     "user1",
		Fee:    100,
		First:  1,
		Second: 1,
	}
	j := NewJackpot(0.0)

	res := Play(lc, j, r)
	t.Logf("result: %v\n", res)
	// should win with bonus
	if res.ResultType != "bonus" {
		t.Log("Should have bonus, found ", res.ResultType)
		t.FailNow()
	}
	// now we should have bonus game
	_, ok := bonusClients["user1"]
	if !ok {
		t.Log("Client id must be in bonusClients")
		t.FailNow()
	}

	// now we can send free game request
	r.Fee = 0.0

	res = Play(lc, j, r)
	// should win
	if res.ResultType != "win" {
		t.Log("Should have win, found ", res.ResultType)
		t.FailNow()
	}
}

func TestMultiClient(t *testing.T) {
	lc := new(failStack)
	j := NewJackpot(0.0)

	for i := 0; i < 5; i++ {
		go func() {
			time.Sleep(time.Duration(rand.Int63n(100)) * time.Millisecond)
			r := Request{
				Id:     "user #" + strconv.Itoa(i),
				Fee:    100,
				First:  1,
				Second: 1,
			}
			res := Play(lc, j, r)
			t.Logf("Result[%d] = %#v \n", i, res)
		}()
	}
	// wait till all finished
	time.Sleep(time.Duration(100) * time.Millisecond)

	//everyone has failed - jackpot has to be total:

	if j.amount != 500 {
		t.Log("Jackpot has to be 500, found ", j.amount)
		t.FailNow()
	}
}
