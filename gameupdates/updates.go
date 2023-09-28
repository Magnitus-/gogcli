package gameupdates

type Updates struct {
	NewGames     []int64
	UpdatedGames []int64
}

func NewEmptyUpdates() *Updates {
	updates := Updates{[]int64{}, []int64{}}
	return &updates
}

func (u *Updates) AddNewGame(gameId int64) {
	(*u).NewGames = append((*u).NewGames, gameId)
}

func (u *Updates) AddUpdatedGame(gameId int64) {
	(*u).UpdatedGames = append((*u).UpdatedGames, gameId)
}
