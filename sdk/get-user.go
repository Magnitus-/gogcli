package sdk

import (
	"encoding/json"
	"errors"
	"fmt"
)

type currency struct {
	Code   string
	Symbol string
}

type updates struct {
	Messages              int
	PendingFriendRequests int
	UnreadChatMessages    int
	Products              int
	Total                 int
}

type language struct {
	Code string
	Name string
}

type wallet struct {
	Currency string
	Amount   float64
}

type purchasedItemsSummary struct {
	Games  int
	Movies int
}

type friend struct {
	Username string
	GalaxyId string
	Avatar   string
}

type User struct {
	Country           string
	SelectedCurrency  currency
	PreferredLanguage language
	Updates           updates
	UserId            string
	Username          string
	GalaxyUserId      string
	Email             string
	Avatar            string
	WalletBalance     wallet
	PurchasedItems    purchasedItemsSummary
	WishlistedItems   int
	Friends           []friend
}

func (u User) Print() {
	fmt.Println("Email:                  ", u.Email)
	fmt.Println("Username:               ", u.Username)
	fmt.Println("Avatar:                 ", u.Avatar)
	fmt.Println("UserId:                 ", u.UserId)
	fmt.Println("GalaxyUserId:           ", u.GalaxyUserId)
	fmt.Println("PreferredLanguage:      ", u.PreferredLanguage.Name)
	fmt.Println("SelectedCurrency:       ", u.SelectedCurrency.Code)
	fmt.Printf("WalletBalance:           %f %s\n", u.WalletBalance.Amount, u.WalletBalance.Currency)
	fmt.Println("WishlistedItems:        ", u.WishlistedItems)
	fmt.Println("PurchasedItems:")
	fmt.Println("  Games:                ", u.PurchasedItems.Games)
	fmt.Println("  Movies:               ", u.PurchasedItems.Movies)
	fmt.Println("Updates:")
	fmt.Println("  Messages:             ", u.Updates.Messages)
	fmt.Println("  PendingFriendRequests:", u.Updates.PendingFriendRequests)
	fmt.Println("  UnreadChatMessages:   ", u.Updates.UnreadChatMessages)
	fmt.Println("  Products:             ", u.Updates.Products)
	fmt.Println("Friends:")

	for _, f := range u.Friends {
		fmt.Println("  - Username:           ", f.Username)
		fmt.Println("    Avatar:             ", f.Avatar)
		fmt.Println("    GalaxyId:           ", f.GalaxyId)
	}
}

func (s *Sdk) GetUser(debug bool) (User, error) {
	var u User

	b, err := s.getUrl(
		"https://embed.gog.com/userData.json",
		"GetUser()",
		debug,
		true,
	)
	if err != nil {
		return u, err
	}

	sErr := json.Unmarshal(b, &u)
	if sErr != nil {
		msg := fmt.Sprintf("Responde deserialization error: %s", sErr.Error())
		return u, errors.New(msg)
	}

	return u, nil
}
