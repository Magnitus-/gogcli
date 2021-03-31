package sdk

import (
	"encoding/json"
	"errors"
	"fmt"
)

type ProductOsCompatibility struct {
	Windows bool
	Osx     bool
	Linux   bool
}

type ProductLinks struct {
	Purchase_link string
	Product_card  string
	Support       string
	Forum         string
}

type ProductDevelopmentInfo struct {
	Active bool
	Until  string
}

type ProductImages struct {
	Background          string
	Logo                string
	Logo2x              string
	Icon                string
	SidebarIcon         string
	SidebarIcon2x       string
	MenuNotificationAv  string
	MenuNotificationAv2 string
}

type ProductDescription struct {
	Lead                string
	Full                string
	Whats_cool_about_it string
}

type ProductScreenShotFormatedImage struct {
	Formatter_name string
	Image_url      string
}

type ProductScreenShot struct {
	Image_id               string
	Formatter_template_url string
	Formatted_images       []ProductScreenShotFormatedImage
}

type ProductVideo struct {
	Video_url     string
	Thumbnail_url string
	Provider      string
}

type Product struct {
	Id                           int64
	Title                        string
	Slug                         string
	Content_system_compatibility ProductOsCompatibility
	Links                        ProductLinks
	In_development               ProductDevelopmentInfo
    Is_pre_order                 bool
	Release_date                 string
	Images                       ProductImages
	Description                  ProductDescription
	Screenshots                  []ProductScreenShot
	Videos                       []ProductVideo
	Changelog                    string
}

func (s *Sdk) GetProduct(gameId int64) (Product, error) {
	var p Product

	fn := fmt.Sprintf("GetProduct(gameId=%d)", gameId)
	u := fmt.Sprintf("https://api.gog.com/products/%d?expand=downloads,expanded_dlcs,description,screenshots,videos,related_products,changelog", gameId)

	b, err := s.getUrl(
		u,
		fn,
		true,
	)
	if err != nil {
		return p, err
	}

	sErr := json.Unmarshal(b, &p)
	if sErr != nil {
		msg := fmt.Sprintf("Responde deserialization error: %s", sErr.Error())
		return p, errors.New(msg)
	}

	return p, nil
}