package sdk

import (
	"encoding/json"
	"errors"
	"fmt"
)

type simpleGalaxyInstaller struct {
	Path string
	Os   string
}

type gameDetailsTags struct {
	Id           string
	Name         string
	ProductCount string
}

type gameDetailsExtra struct {
	ManualUrl string
	Name      string
	Type      string
	Info      int
	Size      string
}

type GameDetails struct {
	Title                  string
	BackgroundImage        string
	CdKey                  string
	TextInformation        string
	Downloads              gameDetailsDownloads
	Extras                 []gameDetailsExtra
	Dlcs                   []GameDetails
	Tags                   []gameDetailsTags
	IsPreOrder             bool
	ReleaseTimestamp       int
	Changelog              string
	ForumLink              string
	IsBaseProductMissing   bool
	Features               []string
	SimpleGalaxyInstallers []simpleGalaxyInstaller
	//messages
	//GalaxyDownloads
}

func (g GameDetails) Print() {
	fmt.Println("Title:           ", g.Title)
	fmt.Println("BackgroundImage: ", g.BackgroundImage)
	fmt.Println("CdKey:           ", g.CdKey)
	fmt.Println("ReleaseTimestamp:", g.ReleaseTimestamp)
	fmt.Println("ForumLink:       ", g.ForumLink)
	if len(g.Features) > 0 {
		fmt.Println("Features:")
		for _, f := range g.Features {
			fmt.Println("  -", f)
		}
	} else {
		fmt.Println("Features: []")
	}
	if len(g.Tags) > 0 {
		fmt.Println("Tags:")
		for _, t := range g.Tags {
			fmt.Println("  -", t.Name)
		}
	} else {
		fmt.Println("Tags: []")
	}
	if len(g.Downloads) > 0 {
		fmt.Println("Downloads:")
		for _, f := range g.Downloads {
			fmt.Println("  - Name:     ", f.Name)
			fmt.Println("    Language: ", f.Language)
			fmt.Println("    Os:       ", f.Os)
			fmt.Println("    Version:  ", f.Version)
			fmt.Println("    Size:     ", f.Size)
			fmt.Println("    Date:     ", f.Date)
			fmt.Println("    ManualUrl:", f.ManualUrl)
		}
	}
	if len(g.Dlcs) > 0 {
		fmt.Println("Dlcs:")
		for _, d := range g.Dlcs {
			fmt.Println("  - Title:          ", d.Title)
			fmt.Println("    BackgroundImage:", d.BackgroundImage)
			fmt.Println("    CdKey:          ", d.CdKey)
			if len(d.Downloads) > 0 {
				fmt.Println("    Downloads:")
				for _, f := range d.Downloads {
					fmt.Println("      - Name:     ", f.Name)
					fmt.Println("        Language: ", f.Language)
					fmt.Println("        Os:       ", f.Os)
					fmt.Println("        Version:  ", f.Version)
					fmt.Println("        Size:     ", f.Size)
					fmt.Println("        Date:     ", f.Date)
					fmt.Println("        ManualUrl:", f.ManualUrl)
				}
			} else {
				fmt.Println("    Downloads: []")
			}
			if len(d.Extras) > 0 {
				fmt.Println("    Extras:")
				for _, e := range d.Extras {
					fmt.Println("      - Name:     ", e.Name)
					fmt.Println("        Type:     ", e.Type)
					fmt.Println("        ManualUrl:", e.ManualUrl)
					fmt.Println("        Size:     ", e.Size)
				}
			} else {
				fmt.Println("    Extras: []")
			}
		}
	}
	if len(g.Extras) > 0 {
		fmt.Println("Extras:")
		for _, e := range g.Extras {
			fmt.Println("  - Name:     ", e.Name)
			fmt.Println("    Type:     ", e.Type)
			fmt.Println("    ManualUrl:", e.ManualUrl)
			fmt.Println("    Size:     ", e.Size)
			fmt.Println("")
		}
	} else {
		fmt.Println("Extras: []")
	}
}

func (s *Sdk) GetGameDetails(gameId int64) (GameDetails, error) {
	var g GameDetails

	fn := fmt.Sprintf("GetGameDetails(gameId=%d)", gameId)
	u := fmt.Sprintf("https://embed.gog.com/account/gameDetails/%d.json", gameId)

	b, _, err := s.getUrl(
		u,
		fn,
		true,
	)
	if err != nil {
		return g, err
	}

	sErr := json.Unmarshal(b, &g)
	if sErr != nil {
		msg := fmt.Sprintf("Responde deserialization error: %s", sErr.Error())
		return g, errors.New(msg)
	}

	return g, nil
}
