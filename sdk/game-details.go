package sdk

import (
	"fmt"
	"io/ioutil"
	"os"
)

type GameDetails struct {
}

func (s Sdk) GetGameDetails(gameId int) string {
	u := fmt.Sprintf("https://embed.gog.com/account//gameDetails/%d.json", gameId)

	c := s.getClient()

	r, err := c.Get(u)
	if err != nil {
		fmt.Println("Game details retrieval request error:", err)
		os.Exit(1)
	}

	b, bErr := ioutil.ReadAll(r.Body)
	if bErr != nil {
		fmt.Println("Owned games retrieval body error:", bErr)
		os.Exit(1)
	}

	return string(b)
}
