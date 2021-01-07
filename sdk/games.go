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

type product struct {
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
	Products        []product
	Tags            []tag
}

func (p product) StringifyOses() string {
	worksOn := "["
	if p.WorksOn.Windows {
		worksOn += "Windows, "
	}
	if p.WorksOn.Mac {
		worksOn += "Mac, "
	}
	if p.WorksOn.Linux {
		worksOn += "Linux, "
	}
	worksOn += "]"
	return worksOn
}

func (p product) StringifyTags(ts []tag) string {
	tags := "["
	for _, pTag := range p.Tags {
		for _, tag := range ts {
			if tag.Id == pTag {
				tags += fmt.Sprintf("%s, ", tag.Name)
			}
		}
	}
	tags += "]"
	return tags
}

func (o OwnedGamesPage) Print() {
	fmt.Println("Page:                  ", o.Page)
	fmt.Println("TotalPages:            ", o.TotalPages)
	fmt.Println("ProductsPerPage:       ", o.ProductsPerPage)
	fmt.Println("Products:")
	for _, p := range o.Products {
		fmt.Println("--------------------")
		fmt.Println("  - Title:             ", p.Title)
		fmt.Println("    Slug:              ", p.Slug)
		fmt.Println("    Id:                ", p.Id)
		fmt.Println("    Image:             ", p.Image)
		fmt.Println("    Url:               ", p.Url)
		fmt.Println("    Category:          ", p.Category)
		fmt.Println("    Tags:              ", p.StringifyTags(o.Tags))
		fmt.Println("    worksOn:           ", p.StringifyOses())
		fmt.Println("    IsNew:             ", p.IsNew)
		fmt.Println("    IsInDevelopment:   ", p.IsInDevelopment)
		fmt.Println("    IsComingSoon:      ", p.IsComingSoon)
		fmt.Println("    IsGalaxyCompatible:", p.IsGalaxyCompatible)
		fmt.Println("    Updates:           ", p.Updates)
		fmt.Println("    DlcCount:          ", p.DlcCount)
	}
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
