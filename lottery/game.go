package lottery

import (
	"errors"
	"fmt"
)

const GameMinimumFee = 100.0

var bonusClients map[string]bool = make(map[string]bool)

// Json-serializable client game request structure
type Request struct {
	Id     string  `json:"id"`
	First  byte    `json:"first"`
	Second byte    `json:"second"`
	Fee    float32 `json:"fee"`
}

// Returns true if the request is valid
func (r Request) Validate() error {
	if len(r.Id) == 0 {
		return errors.New("user id must not be empty")
	}
	if r.Fee < GameMinimumFee {
		if _, ok := bonusClients[r.Id]; ok {
			delete(bonusClients, r.Id)
		} else {
			return errors.New(fmt.Sprintf("game fee has to be >= %5.2f", GameMinimumFee))
		}
	}

	return nil
}

// Json-serializable game result structure
type GameResult struct {
	// result of the play: 'win', 'lost', or 'bonus'
	ResultType string  `json:"type"`
	WinAmount  float32 `json:"amount"`
}

func Play(lc LuckyStack, j *Jackpot, r Request) GameResult {
	res := GameResult{
		ResultType: "lost",
		WinAmount:  0.0,
	}

	j.Lock()
	defer j.Done()

	if lc.Check(r.First, r.Second) {
		if j.amount == 0.0 {
			bonusClients[r.Id] = true
			res.ResultType = "bonus"
		} else {
			res.ResultType = "win"
		}
		res.WinAmount = j.amount

		fmt.Printf("We have a winner !!! Jackpot is %f \n", j.amount)
	}
	j.amount += r.Fee

	return res
}

// Jackpot structure serves as account and synchronization point
// for all Play() requests
type Jackpot struct {
	amount float32
	sig    chan int
}

// Creates and initializes new jackpot structure
func NewJackpot(initialSize float32) *Jackpot {
	return &Jackpot{
		amount: initialSize,
		sig:    make(chan int, 1),
	}
}

// Locks current state of jackpot, if called by more than
// one client - holds until j.Done() is called
func (j *Jackpot) Lock() {
	j.sig <- 1
}

// Releases current state of jackpot
func (j *Jackpot) Done() {
	<-j.sig
}
