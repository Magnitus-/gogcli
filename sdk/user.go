package sdk

import (
	"fmt"
	"io"
	"os"
	"strings"
)

func (s Sdk) GetUser() string {
	c := s.getClient()
	r, err := c.Get("https://embed.gog.com/userData.json")
	if err != nil {
		fmt.Println("User retrieval request error:", err)
		os.Exit(1)
	}
	var b strings.Builder
	_, bErr := io.Copy(&b, r.Body)
	if bErr != nil {
		fmt.Println("User retrieval body error:", bErr)
		os.Exit(1)
	}
	return b.String()
}
