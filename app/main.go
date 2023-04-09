package main

import (
	"fmt"
	. "github.com/pkarpovich/esport-syncer/app/providers"
)

func getMatches(p Provider) (error, []Match) {
	return p.GetMatches()
}

func main() {
	dotaProvider := DotaProvider{TeamID: "6224"}
	err, matches := getMatches(&dotaProvider)
	if err != nil {
		fmt.Println(err)
		return
	}

	for i, match := range matches {
		fmt.Printf("%d. %s vs %s at %s (Tournament: %s | Result: %s)\n", i+1, match.Team1, match.Team2, match.Time, match.Tournament, match.Score)
	}
}
