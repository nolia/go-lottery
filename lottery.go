package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/nolia/go-lottery/lottery"
)

func main() {
	luckyStack := lottery.NewLuckyStack()
	jackpot := lottery.NewJackpot(0.0)

	fmt.Println("Starting lottery server at port 8080 ...")

	http.HandleFunc("/play.json", func(rw http.ResponseWriter, req *http.Request) {
		decoder := json.NewDecoder(req.Body)
		var gameRequest lottery.Request
		err := decoder.Decode(&gameRequest)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Decoding err: %v\n", err)
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}
		err = gameRequest.Validate()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v; Invalid request: %v\n", err, gameRequest)
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}

		result := lottery.Play(luckyStack, jackpot, gameRequest)
		fmt.Printf("Request from %v, [%v,%v] : result %v \n", gameRequest.Id, gameRequest.First, gameRequest.Second, result)

		rw.WriteHeader(http.StatusOK)
		rw.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(rw).Encode(result)

	})

	http.ListenAndServe(":8080", nil)
}
