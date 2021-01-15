package sdk

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"time"
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
		fmt.Println("")
	}
}

func (s *Sdk) GetOwnedGames(page int, search string, debug bool) (OwnedGamesPage, error) {
	var o OwnedGamesPage

	fn := fmt.Sprintf("GetOwnedGames(page=%d, search=%s)", page, search)
	u := fmt.Sprintf("https://embed.gog.com/account/getFilteredProducts?mediaType=1&page=%d", page)
	if search != "" {
		u += fmt.Sprintf("&search=%s", url.QueryEscape(search))
	}

	b, err := s.getUrl(
		u,
		fn,
		debug,
		true,
	)
	if err != nil {
		return o, err
	}

	sErr := json.Unmarshal(b, &o)
	if sErr != nil {
		msg := fmt.Sprintf("Responde deserialization error: %s", sErr.Error())
		return o, errors.New(msg)
	}

	return o, nil
}

type OwnedGamesPageReturn struct {
	page OwnedGamesPage
	err  error
}

func (s *Sdk) GetOwnedGamesPageAsync(page int, search string, debug bool, returnVal chan OwnedGamesPageReturn) {
	o, err := s.GetOwnedGames(page, search, debug)
	returnVal <- OwnedGamesPageReturn{page: o, err: err}
}

func (s *Sdk) GetAllOwnedGamesPages(search string, concurrency int, pause int, debug bool) ([]OwnedGamesPage, []error) {
	fnCall := fmt.Sprintf(
		"GetAllOwnedGamesPages(search=%s, concurrency=%d, pause=%d",
		search,
		concurrency,
		pause,
	)
	var pageCount int
	var currentPage int
	var pages []OwnedGamesPage
	var errs []error
	var callVal OwnedGamesPageReturn
	c := make(chan OwnedGamesPageReturn)

	if debug {
		(*s).logger.Println(
			fmt.Sprintf("%s -> fetching first page", fnCall),
		)
	}
	go s.GetOwnedGamesPageAsync(1, search, debug, c)
	callVal = <-c
	if callVal.err != nil {
		errs = append(errs, callVal.err)
		return pages, errs
	}

	if callVal.page.TotalPages == 0 {
		return pages, errs
	}

	pages = append(pages, callVal.page)
	pageCount = callVal.page.TotalPages
	currentPage = callVal.page.Page + 1

	if debug {
		(*s).logger.Println(
			fmt.Sprintf("%s -> total pages to fetch: %d", fnCall, pageCount),
		)
	}

	for currentPage < pageCount {
		maxPage := min(currentPage+concurrency, pageCount)

		if debug {
			(*s).logger.Println(
				fmt.Sprintf("%s -> fetching page %d to %d", fnCall, currentPage, maxPage),
			)
		}

		for i := currentPage + 1; i <= maxPage; i++ {
			go s.GetOwnedGamesPageAsync(i, search, debug, c)
		}

		for i := currentPage + 1; i <= maxPage; i++ {
			callVal = <-c
			if callVal.err != nil {
				errs = append(errs, callVal.err)
			} else {
				pages = append(pages, callVal.page)
			}
		}

		if len(errs) > 0 {
			return pages, errs
		}

		currentPage = maxPage
		time.Sleep(time.Duration(pause) * time.Millisecond)
	}
	return pages, errs
}

func min(x int, y int) int {
	if x < y {
		return x
	}
	return y
}
