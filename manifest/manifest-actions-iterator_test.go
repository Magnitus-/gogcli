package manifest

import (
	"testing"
)

func getFixtures() (Manifest, GameActions) {
	manifest := Manifest{
		Games: []ManifestGame{
			ManifestGame{
				Id:           1,
				Title:        "Croom",
				VerifiedSize: 10,
			},
			ManifestGame{
				Id:           2,
				Title:        "Broom",
				VerifiedSize: 30,
			},
			ManifestGame{
				Id:           3,
				Title:        "Aroom",
				VerifiedSize: 20,
			},
		},
	}

	actions := GameActions{
		1: GameAction{
			Title:  "Croom",
			Id:     1,
			Action: "update",
			InstallerActions: map[string]FileAction{
				"CroomInstaller1": FileAction{
					Title:  "CroomInstaller1",
					Name:   "CroomInstaller1",
					Url:    "CroomInstaller1",
					Kind:   "installer",
					Action: "add",
				},
				"CroomInstaller2": FileAction{
					Title:  "CroomInstaller2",
					Name:   "CroomInstaller2",
					Url:    "CroomInstaller2",
					Kind:   "installer",
					Action: "remove",
				},
			},
			ExtraActions: map[string]FileAction{},
		},
		2: GameAction{
			Title:            "Broom",
			Id:               2,
			Action:           "add",
			InstallerActions: map[string]FileAction{},
			ExtraActions: map[string]FileAction{
				"BroomExtra1": FileAction{
					Title:  "BroomExtra1",
					Name:   "BroomExtra1",
					Url:    "BroomExtra1",
					Kind:   "extra",
					Action: "add",
				},
				"BroomExtra2": FileAction{
					Title:  "BroomExtra2",
					Name:   "BroomExtra2",
					Url:    "BroomExtra2",
					Kind:   "extra",
					Action: "remove",
				},
			},
		},
		3: GameAction{
			Title:  "Aroom",
			Id:     3,
			Action: "remove",
			InstallerActions: map[string]FileAction{
				"AroomInstaller1": FileAction{
					Title:  "AroomInstaller1",
					Name:   "AroomInstaller1",
					Url:    "AroomInstaller1",
					Kind:   "installer",
					Action: "add",
				},
			},
			ExtraActions: map[string]FileAction{
				"AroomExtra1": FileAction{
					Title:  "AroomExtra1",
					Name:   "AroomExtra1",
					Url:    "AroomExtra1",
					Kind:   "extra",
					Action: "remove",
				},
			},
		},
	}

	return manifest, actions
}

func equivalentFileActions(fileActionPtr1 *FileAction, fileActionPtr2 *FileAction) bool {
	if fileActionPtr1 == nil && fileActionPtr2 == nil {
		return true
	}
	if (fileActionPtr1 == nil && fileActionPtr2 != nil) || (fileActionPtr1 != nil && fileActionPtr2 == nil) {
		return false
	}
	return (*fileActionPtr1).Name == (*fileActionPtr2).Name && (*fileActionPtr1).Title == (*fileActionPtr2).Title && (*fileActionPtr1).Url == (*fileActionPtr2).Url && (*fileActionPtr1).Kind == (*fileActionPtr2).Kind && (*fileActionPtr1).Action == (*fileActionPtr2).Action
}

func expect(i ActionsIterator, actions []Action, t *testing.T) {
	for _, action := range actions {
		if !i.ShouldContinue() {
			t.Errorf("Iterator should indicate that it wants to continue when not finished")
		}

		next, _ := i.Next()
		if next.IsFileAction != action.IsFileAction || next.GameId != action.GameId || next.GameAction != action.GameAction || (!equivalentFileActions(next.FileActionPtr, action.FileActionPtr)) {
			t.Errorf("Expected action %v and got %v", next, action)
		}
	}

	if i.ShouldContinue() {
		t.Errorf("Iterator should not indicate that it wants to continue when finished")
	}
}

func TestFullIdSort(t *testing.T) {
	manifest, actions := getFixtures()
	iterator := NewActionsIterator(actions, -1)
	iterator.Sort(ActionsIteratorSort{
		[]int64{},
		"id",
		true,
	}, &manifest)

	fileAction1 := FileAction{
		Title:  "CroomInstaller1",
		Name:   "CroomInstaller1",
		Url:    "CroomInstaller1",
		Kind:   "installer",
		Action: "add",
	}
	fileAction2 := FileAction{
		Title:  "CroomInstaller2",
		Name:   "CroomInstaller2",
		Url:    "CroomInstaller2",
		Kind:   "installer",
		Action: "remove",
	}
	fileAction3 := FileAction{
		Title:  "BroomExtra1",
		Name:   "BroomExtra1",
		Url:    "BroomExtra1",
		Kind:   "extra",
		Action: "add",
	}
	fileAction4 := FileAction{
		Title:  "BroomExtra2",
		Name:   "BroomExtra2",
		Url:    "BroomExtra2",
		Kind:   "extra",
		Action: "remove",
	}
	fileAction5 := FileAction{
		Title:  "AroomInstaller1",
		Name:   "AroomInstaller1",
		Url:    "AroomInstaller1",
		Kind:   "installer",
		Action: "add",
	}
	fileAction6 := FileAction{
		Title:  "AroomExtra1",
		Name:   "AroomExtra1",
		Url:    "AroomExtra1",
		Kind:   "extra",
		Action: "remove",
	}
	expectedActions := []Action{
		Action{
			1,
			true,
			&fileAction1,
			"",
		},
		Action{
			1,
			true,
			&fileAction2,
			"",
		},
		Action{
			2,
			false,
			nil,
			"add",
		},
		Action{
			2,
			true,
			&fileAction3,
			"",
		},
		Action{
			2,
			true,
			&fileAction4,
			"",
		},
		Action{
			3,
			true,
			&fileAction5,
			"",
		},
		Action{
			3,
			true,
			&fileAction6,
			"",
		},
		Action{
			3,
			false,
			nil,
			"remove",
		},
	}
	expect(*iterator, expectedActions, t)
}

func TestFullIdSortWithPreferredId(t *testing.T) {
	manifest, actions := getFixtures()
	iterator := NewActionsIterator(actions, -1)
	iterator.Sort(ActionsIteratorSort{
		[]int64{3},
		"id",
		true,
	}, &manifest)

	fileAction1 := FileAction{
		Title:  "CroomInstaller1",
		Name:   "CroomInstaller1",
		Url:    "CroomInstaller1",
		Kind:   "installer",
		Action: "add",
	}
	fileAction2 := FileAction{
		Title:  "CroomInstaller2",
		Name:   "CroomInstaller2",
		Url:    "CroomInstaller2",
		Kind:   "installer",
		Action: "remove",
	}
	fileAction3 := FileAction{
		Title:  "BroomExtra1",
		Name:   "BroomExtra1",
		Url:    "BroomExtra1",
		Kind:   "extra",
		Action: "add",
	}
	fileAction4 := FileAction{
		Title:  "BroomExtra2",
		Name:   "BroomExtra2",
		Url:    "BroomExtra2",
		Kind:   "extra",
		Action: "remove",
	}
	fileAction5 := FileAction{
		Title:  "AroomInstaller1",
		Name:   "AroomInstaller1",
		Url:    "AroomInstaller1",
		Kind:   "installer",
		Action: "add",
	}
	fileAction6 := FileAction{
		Title:  "AroomExtra1",
		Name:   "AroomExtra1",
		Url:    "AroomExtra1",
		Kind:   "extra",
		Action: "remove",
	}
	expectedActions := []Action{
		Action{
			3,
			true,
			&fileAction5,
			"",
		},
		Action{
			3,
			true,
			&fileAction6,
			"",
		},
		Action{
			3,
			false,
			nil,
			"remove",
		},
		Action{
			1,
			true,
			&fileAction1,
			"",
		},
		Action{
			1,
			true,
			&fileAction2,
			"",
		},
		Action{
			2,
			false,
			nil,
			"add",
		},
		Action{
			2,
			true,
			&fileAction3,
			"",
		},
		Action{
			2,
			true,
			&fileAction4,
			"",
		},
	}
	expect(*iterator, expectedActions, t)
}

func TestMax2TitleSort(t *testing.T) {
	manifest, actions := getFixtures()
	iterator := NewActionsIterator(actions, 2)
	iterator.Sort(ActionsIteratorSort{
		[]int64{},
		"title",
		true,
	}, &manifest)

	fileAction3 := FileAction{
		Title:  "BroomExtra1",
		Name:   "BroomExtra1",
		Url:    "BroomExtra1",
		Kind:   "extra",
		Action: "add",
	}
	fileAction4 := FileAction{
		Title:  "BroomExtra2",
		Name:   "BroomExtra2",
		Url:    "BroomExtra2",
		Kind:   "extra",
		Action: "remove",
	}
	fileAction5 := FileAction{
		Title:  "AroomInstaller1",
		Name:   "AroomInstaller1",
		Url:    "AroomInstaller1",
		Kind:   "installer",
		Action: "add",
	}
	fileAction6 := FileAction{
		Title:  "AroomExtra1",
		Name:   "AroomExtra1",
		Url:    "AroomExtra1",
		Kind:   "extra",
		Action: "remove",
	}
	expectedActions := []Action{
		Action{
			3,
			true,
			&fileAction5,
			"",
		},
		Action{
			3,
			true,
			&fileAction6,
			"",
		},
		Action{
			3,
			false,
			nil,
			"remove",
		},
		Action{
			2,
			false,
			nil,
			"add",
		},
		Action{
			2,
			true,
			&fileAction3,
			"",
		},
		Action{
			2,
			true,
			&fileAction4,
			"",
		},
	}
	expect(*iterator, expectedActions, t)
}

func TestMax2SizeSort(t *testing.T) {
	manifest, actions := getFixtures()
	iterator := NewActionsIterator(actions, 2)
	iterator.Sort(ActionsIteratorSort{
		[]int64{},
		"size",
		true,
	}, &manifest)

	fileAction1 := FileAction{
		Title:  "CroomInstaller1",
		Name:   "CroomInstaller1",
		Url:    "CroomInstaller1",
		Kind:   "installer",
		Action: "add",
	}
	fileAction2 := FileAction{
		Title:  "CroomInstaller2",
		Name:   "CroomInstaller2",
		Url:    "CroomInstaller2",
		Kind:   "installer",
		Action: "remove",
	}
	fileAction5 := FileAction{
		Title:  "AroomInstaller1",
		Name:   "AroomInstaller1",
		Url:    "AroomInstaller1",
		Kind:   "installer",
		Action: "add",
	}
	fileAction6 := FileAction{
		Title:  "AroomExtra1",
		Name:   "AroomExtra1",
		Url:    "AroomExtra1",
		Kind:   "extra",
		Action: "remove",
	}
	expectedActions := []Action{
		Action{
			1,
			true,
			&fileAction1,
			"",
		},
		Action{
			1,
			true,
			&fileAction2,
			"",
		},
		Action{
			3,
			true,
			&fileAction5,
			"",
		},
		Action{
			3,
			true,
			&fileAction6,
			"",
		},
		Action{
			3,
			false,
			nil,
			"remove",
		},
	}
	expect(*iterator, expectedActions, t)
}

func TestMax2DescSizeSort(t *testing.T) {
	manifest, actions := getFixtures()
	iterator := NewActionsIterator(actions, 2)
	iterator.Sort(ActionsIteratorSort{
		[]int64{},
		"size",
		false,
	}, &manifest)

	fileAction3 := FileAction{
		Title:  "BroomExtra1",
		Name:   "BroomExtra1",
		Url:    "BroomExtra1",
		Kind:   "extra",
		Action: "add",
	}
	fileAction4 := FileAction{
		Title:  "BroomExtra2",
		Name:   "BroomExtra2",
		Url:    "BroomExtra2",
		Kind:   "extra",
		Action: "remove",
	}
	fileAction5 := FileAction{
		Title:  "AroomInstaller1",
		Name:   "AroomInstaller1",
		Url:    "AroomInstaller1",
		Kind:   "installer",
		Action: "add",
	}
	fileAction6 := FileAction{
		Title:  "AroomExtra1",
		Name:   "AroomExtra1",
		Url:    "AroomExtra1",
		Kind:   "extra",
		Action: "remove",
	}
	expectedActions := []Action{
		Action{
			2,
			false,
			nil,
			"add",
		},
		Action{
			2,
			true,
			&fileAction3,
			"",
		},
		Action{
			2,
			true,
			&fileAction4,
			"",
		},
		Action{
			3,
			true,
			&fileAction5,
			"",
		},
		Action{
			3,
			true,
			&fileAction6,
			"",
		},
		Action{
			3,
			false,
			nil,
			"remove",
		},
	}
	expect(*iterator, expectedActions, t)
}
