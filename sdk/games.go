package sdk

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type tag struct {
	Id           string
	Name         string
	ProductCount string
}

type productAvailability struct {
	IsAvailable          bool
	IsAvailableInAccount bool
}

type productTags []string

type productReleaseDate struct {
	Date          string
	Timezone      string
	Timezone_type int
}

type productWorksOn struct {
	Linux   bool
	Mac     bool
	Windows bool
}

type products struct {
	Id                   int
	IsNew                bool
	Updates              int
	IsHidden             bool
	Title                string
	Slug                 string
	Category             string
	Rating               int
	Image                string
	Url                  string
	DlcCount             int
	Tags                 productTags
	Availability         productAvailability
	IsInDevelopment      bool
	IsGalaxyCompatible   bool
	IsBaseProductMissing bool
	IsComingSoon         bool
	ReleaseDate          productReleaseDate
	WorksOn              productWorksOn
}

type OwnedGamesPage struct {
	Page            int
	TotalPages      int
	ProductsPerPage int
	TotalProducts   int
	Products        []products
	Tags            []tag
}

func (s Sdk) GetOwnedGames(page int) OwnedGamesPage {
	u := fmt.Sprintf("https://embed.gog.com/account/getFilteredProducts?mediaType=1&page=%d", page)
	c := s.getClient()

	r, err := c.Get(u)
	if err != nil {
		fmt.Println("Owned games retrieval request error:", err)
		os.Exit(1)
	}

	b, bErr := ioutil.ReadAll(r.Body)
	if bErr != nil {
		fmt.Println("Owned games retrieval body error:", bErr)
		os.Exit(1)
	}

	var o OwnedGamesPage
	sErr := json.Unmarshal(b, &o)
	if sErr != nil {
		fmt.Println("Responde deserialization error:", sErr)
	}

	return o
}
