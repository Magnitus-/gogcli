package gameupdates

import (
	"testing"
)

func TestUpdates(t *testing.T) {
	u := NewEmptyUpdates()

	if len(u.NewGames) > 0 {
		t.Errorf("New empty updates creates non-empty array of new games")
	}

	if len(u.UpdatedGames) > 0 {
		t.Errorf("New empty updates creates non-empty array of updated games")
	}

	u.AddNewGame(1)
	if len(u.NewGames) != 1 || u.NewGames[0] != 1 {
		t.Errorf("Adding new game to empty new games array is not behaving as expected")
	}

	u.AddNewGame(2)
	if len(u.NewGames) != 2 || u.NewGames[1] != 2 {
		t.Errorf("Adding new game to non-empty new games array is not behaving as expected")
	}

	u.AddUpdatedGame(1)
	if len(u.UpdatedGames) != 1 || u.UpdatedGames[0] != 1 {
		t.Errorf("Adding new updated game to empty updated games array is not behaving as expected")
	}

	u.AddUpdatedGame(2)
	if len(u.UpdatedGames) != 2 || u.UpdatedGames[1] != 2 {
		t.Errorf("Adding new updated game to non-empty updated games array is not behaving as expected")
	}
}
