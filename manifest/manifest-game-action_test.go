package manifest

import (
	"testing"
)

func TestUpdate(t *testing.T) {

	fileAdditions := []FileAction{
		FileAction{
			Title:  "one",
			Name:   "one",
			Url:    "donotcare",
			Kind:   "installer",
			Action: "add",
		},
		FileAction{
			Title:  "two",
			Name:   "two",
			Url:    "donotcare",
			Kind:   "installer",
			Action: "add",
		},
		FileAction{
			Title:  "three",
			Name:   "three",
			Url:    "donotcare",
			Kind:   "installer",
			Action: "add",
		},
	}

	fileRemovals := []FileAction{
		FileAction{
			Title:  "one",
			Name:   "one",
			Url:    "donotcare",
			Kind:   "installer",
			Action: "remove",
		},
		FileAction{
			Title:  "two",
			Name:   "two",
			Url:    "donotcare",
			Kind:   "installer",
			Action: "remove",
		},
		FileAction{
			Title:  "three",
			Name:   "three",
			Url:    "donotcare",
			Kind:   "installer",
			Action: "remove",
		},
	}

	addAction := GameAction{
		Title:  "test",
		Id:     1,
		Action: "add",
		InstallerActions: map[string]FileAction{
			"one": fileAdditions[0],
			"two": fileAdditions[1],
		},
		ExtraActions: map[string]FileAction{
			"one": fileAdditions[0],
			"two": fileAdditions[1],
		},
	}
	addToRemoveAction := GameAction{
		Title:  "test",
		Id:     1,
		Action: "remove",
		InstallerActions: map[string]FileAction{
			"one": fileRemovals[0],
			"two": fileRemovals[1],
		},
		ExtraActions: map[string]FileAction{
			"one": fileRemovals[0],
			"two": fileRemovals[1],
		},
	}

	addAction.Update(&addToRemoveAction)
	if addAction.Action != "remove" {
		t.Errorf("Update with a removal action should have set the action to remove")
	}

	installerActionOne := addAction.InstallerActions["one"]
	installerActionTwo := addAction.InstallerActions["two"]
	if installerActionOne.Action != "remove" || installerActionTwo.Action != "remove" {
		t.Errorf("Installer actions didn't get properly updated in add to removal conversion")
	}

	extraActionOne := addAction.ExtraActions["one"]
	extraActionTwo := addAction.ExtraActions["two"]
	if extraActionOne.Action != "remove" || extraActionTwo.Action != "remove" {
		t.Errorf("Extra actions didn't get properly updated in add to removal conversion")
	}

	addActionTakeTwo := GameAction{
		Title:  "test",
		Id:     1,
		Action: "add",
		InstallerActions: map[string]FileAction{
			"one": fileAdditions[0],
			"two": fileAdditions[1],
		},
		ExtraActions: map[string]FileAction{
			"one": fileAdditions[0],
			"two": fileAdditions[1],
		},
	}
	addToUpdateAction := GameAction{
		Title:  "test",
		Id:     1,
		Action: "update",
		InstallerActions: map[string]FileAction{
			"two":   fileRemovals[1],
			"three": fileAdditions[2],
		},
		ExtraActions: map[string]FileAction{
			"two":   fileRemovals[1],
			"three": fileAdditions[2],
		},
	}

	addActionTakeTwo.Update(&addToUpdateAction)
	if addActionTakeTwo.Action != "add" {
		t.Errorf("Update with an update action should not have changed the original action's type")
	}

	installerActionOne = addActionTakeTwo.InstallerActions["one"]
	installerActionTwo = addActionTakeTwo.InstallerActions["two"]
	installerActionThree := addActionTakeTwo.InstallerActions["three"]
	if installerActionOne.Action != "add" || installerActionTwo.Action != "remove" || installerActionThree.Action != "add" {
		t.Errorf("Installer actions didn't get properly updated in add to update conversion")
	}

	extraActionOne = addActionTakeTwo.ExtraActions["one"]
	extraActionTwo = addActionTakeTwo.ExtraActions["two"]
	extraActionThree := addActionTakeTwo.ExtraActions["three"]
	if extraActionOne.Action != "add" || extraActionTwo.Action != "remove" || extraActionThree.Action != "add" {
		t.Errorf("Extra actions didn't get properly updated in add to update conversion")
	}

}

func TestGameActionIsNoOp(t *testing.T) {
	addOnly := GameAction{
		Title:            "test",
		Id:               1,
		Action:           "add",
		InstallerActions: map[string]FileAction{},
		ExtraActions:     map[string]FileAction{},
	}

	if addOnly.IsNoOp() {
		t.Errorf("Add game action is not noop")
	}

	removeOnly := GameAction{
		Title:            "test",
		Id:               1,
		Action:           "remove",
		InstallerActions: map[string]FileAction{},
		ExtraActions:     map[string]FileAction{},
	}

	if removeOnly.IsNoOp() {
		t.Errorf("Remove game action is not noop")
	}

	updateOnly := GameAction{
		Title:            "test",
		Id:               1,
		Action:           "update",
		InstallerActions: map[string]FileAction{},
		ExtraActions:     map[string]FileAction{},
	}

	if !updateOnly.IsNoOp() {
		t.Errorf("Update game action without file actions should be noop")
	}

	installerAction := FileAction{
		Title:  "installer",
		Name:   "installer",
		Url:    "donotcare",
		Kind:   "installer",
		Action: "add",
	}

	extraAction := FileAction{
		Title:  "extra",
		Name:   "extra",
		Url:    "donotcare",
		Kind:   "extra",
		Action: "add",
	}

	addAndFileActions := GameAction{
		Title:  "test",
		Id:     1,
		Action: "add",
		InstallerActions: map[string]FileAction{
			"installer": installerAction,
		},
		ExtraActions: map[string]FileAction{
			"extra": extraAction,
		},
	}

	if addAndFileActions.IsNoOp() {
		t.Errorf("Add game action with file actions is not noop")
	}

	removeAndFileActions := GameAction{
		Title:  "test",
		Id:     1,
		Action: "remove",
		InstallerActions: map[string]FileAction{
			"installer": installerAction,
		},
		ExtraActions: map[string]FileAction{
			"extra": extraAction,
		},
	}

	if removeAndFileActions.IsNoOp() {
		t.Errorf("Remove game action with file actions is not noop")
	}

	updateAndFileActions := GameAction{
		Title:  "test",
		Id:     1,
		Action: "update",
		InstallerActions: map[string]FileAction{
			"installer": installerAction,
		},
		ExtraActions: map[string]FileAction{
			"extra": extraAction,
		},
	}

	if updateAndFileActions.IsNoOp() {
		t.Errorf("Update game action with file actions is not noop")
	}
}

func TestGameActionGetInstallerNames(t *testing.T) {
	extraAction := FileAction{
		Title:  "extra",
		Name:   "extra",
		Url:    "donotcare",
		Kind:   "extra",
		Action: "add",
	}

	noInstallers := GameAction{
		Title:            "test",
		Id:               1,
		Action:           "add",
		InstallerActions: map[string]FileAction{},
		ExtraActions: map[string]FileAction{
			"extra": extraAction,
		},
	}

	if len(noInstallers.GetInstallerNames()) != 0 {
		t.Errorf("Game action with no installers should not return any installer names")
	}

	installerOne := FileAction{
		Title:  "oneTitle",
		Name:   "oneName",
		Url:    "donotcare",
		Kind:   "installer",
		Action: "add",
	}

	installerTwo := FileAction{
		Title:  "twoTitle",
		Name:   "twoName",
		Url:    "donotcare",
		Kind:   "installer",
		Action: "add",
	}

	twoInstallers := GameAction{
		Title:  "test",
		Id:     1,
		Action: "add",
		InstallerActions: map[string]FileAction{
			"oneName": installerOne,
			"twoName": installerTwo,
		},
		ExtraActions: map[string]FileAction{
			"extra": extraAction,
		},
	}

	twoInstallersNames := twoInstallers.GetInstallerNames()
	if !(len(twoInstallersNames) == 2 && stringInSlice("oneName", twoInstallersNames) && stringInSlice("twoName", twoInstallersNames)) {
		t.Errorf("Installer names for game with two installers is not as expected")
	}
}

func TestGameActionGetExtraNames(t *testing.T) {
	installerAction := FileAction{
		Title:  "installer",
		Name:   "installer",
		Url:    "donotcare",
		Kind:   "installer",
		Action: "add",
	}

	noExtras := GameAction{
		Title:  "test",
		Id:     1,
		Action: "add",
		InstallerActions: map[string]FileAction{
			"installer": installerAction,
		},
		ExtraActions: map[string]FileAction{},
	}

	if len(noExtras.GetExtraNames()) != 0 {
		t.Errorf("Game action with no extras should not return any extra names")
	}

	extraOne := FileAction{
		Title:  "oneTitle",
		Name:   "oneName",
		Url:    "donotcare",
		Kind:   "extra",
		Action: "add",
	}

	extraTwo := FileAction{
		Title:  "twoTitle",
		Name:   "twoName",
		Url:    "donotcare",
		Kind:   "extra",
		Action: "add",
	}

	twoExtras := GameAction{
		Title:  "test",
		Id:     1,
		Action: "add",
		InstallerActions: map[string]FileAction{
			"installer": installerAction,
		},
		ExtraActions: map[string]FileAction{
			"oneName": extraOne,
			"twoName": extraTwo,
		},
	}

	twoExtrasNames := twoExtras.GetExtraNames()
	if !(len(twoExtrasNames) == 2 && stringInSlice("oneName", twoExtrasNames) && stringInSlice("twoName", twoExtrasNames)) {
		t.Errorf("Extra names for game with two extras is not as expected")
	}
}

func TestGameActionCountFileActions(t *testing.T) {
	noFileActions := GameAction{
		Title:            "test",
		Id:               1,
		Action:           "add",
		InstallerActions: map[string]FileAction{},
		ExtraActions:     map[string]FileAction{},
	}

	if noFileActions.CountFileActions() != 0 {
		t.Errorf("Game action with no files should have file actions count of 0")
	}

	installerOneAction := FileAction{
		Title:  "installer",
		Name:   "installerOne",
		Url:    "donotcare",
		Kind:   "installer",
		Action: "add",
	}

	installerTwoAction := FileAction{
		Title:  "installer",
		Name:   "installerTwo",
		Url:    "donotcare",
		Kind:   "installer",
		Action: "add",
	}

	twoInstallerActions := GameAction{
		Title:  "test",
		Id:     1,
		Action: "add",
		InstallerActions: map[string]FileAction{
			"installerOne": installerOneAction,
			"installerTwo": installerTwoAction,
		},
		ExtraActions: map[string]FileAction{},
	}

	if twoInstallerActions.CountFileActions() != 2 {
		t.Errorf("Game action with two installers should have file actions count of 0")
	}

	extraAction := FileAction{
		Title:  "extra",
		Name:   "extra",
		Url:    "donotcare",
		Kind:   "extra",
		Action: "add",
	}

	oneExtraActions := GameAction{
		Title:            "test",
		Id:               1,
		Action:           "add",
		InstallerActions: map[string]FileAction{},
		ExtraActions: map[string]FileAction{
			"extra": extraAction,
		},
	}

	if oneExtraActions.CountFileActions() != 1 {
		t.Errorf("Game action with one extra should have file actions count of 1")
	}

	mixActions := GameAction{
		Title:  "test",
		Id:     1,
		Action: "add",
		InstallerActions: map[string]FileAction{
			"installerOne": installerOneAction,
			"installerTwo": installerTwoAction,
		},
		ExtraActions: map[string]FileAction{
			"extra": extraAction,
		},
	}

	if mixActions.CountFileActions() != 3 {
		t.Errorf("Game action with two installers and one extra should have file actions count of 3")
	}
}

func TestGameActionActionsLeft(t *testing.T) {
	installerOneAction := FileAction{
		Title:  "installer",
		Name:   "installerOne",
		Url:    "donotcare",
		Kind:   "installer",
		Action: "add",
	}

	installerTwoAction := FileAction{
		Title:  "installer",
		Name:   "installerTwo",
		Url:    "donotcare",
		Kind:   "installer",
		Action: "add",
	}

	extraAction := FileAction{
		Title:  "extra",
		Name:   "extra",
		Url:    "donotcare",
		Kind:   "extra",
		Action: "add",
	}

	mixActionsAdd := GameAction{
		Title:  "test",
		Id:     1,
		Action: "add",
		InstallerActions: map[string]FileAction{
			"installerOne": installerOneAction,
			"installerTwo": installerTwoAction,
		},
		ExtraActions: map[string]FileAction{
			"extra": extraAction,
		},
	}

	if mixActionsAdd.ActionsLeft() != 4 {
		t.Errorf("Add game action with two installers and one extra should have 4 actions left")
	}

	mixActionsRemove := GameAction{
		Title:  "test",
		Id:     1,
		Action: "remove",
		InstallerActions: map[string]FileAction{
			"installerOne": installerOneAction,
			"installerTwo": installerTwoAction,
		},
		ExtraActions: map[string]FileAction{
			"extra": extraAction,
		},
	}

	if mixActionsRemove.ActionsLeft() != 4 {
		t.Errorf("Remove game action with two installers and one extra should have 4 actions left")
	}

	mixActionsUpdate := GameAction{
		Title:  "test",
		Id:     1,
		Action: "update",
		InstallerActions: map[string]FileAction{
			"installerOne": installerOneAction,
			"installerTwo": installerTwoAction,
		},
		ExtraActions: map[string]FileAction{
			"extra": extraAction,
		},
	}

	if mixActionsUpdate.ActionsLeft() != 3 {
		t.Errorf("Update game action with two installers and one extra should have 3 actions left")
	}

}
