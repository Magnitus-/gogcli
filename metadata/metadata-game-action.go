package metadata

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

type FileAction struct {
	Tag    string
	Name   string
	Url    string
	Action string
}

type ImageId struct {
	Tag    string
	Name   string
}

type SerializableGameAction struct {
	Title        string
	Slug         string
	Id           int64
	Action       string
	ImageActions map[string]FileAction
}

type GameAction struct {
	Title        string
	Slug         string
	Id           int64
	Action       string
	ImageActions map[ImageId]FileAction
}

func (g *GameAction) MarshalJSON() ([]byte, error) {
	conv := SerializableGameAction{
		Title: g.Title,
		Slug: g.Slug,
		Id: g.Id,
		Action: g.Action,
		ImageActions: map[string]FileAction{},
	}

	for id, image := range g.ImageActions {
		conv.ImageActions[fmt.Sprintf("%s/%s", id.Tag, id.Name)] = image
	}
	
	return json.Marshal(&conv)
}

func (g *GameAction) UnmarshalJSON(data []byte) error {
	var conv SerializableGameAction
	if err := json.Unmarshal(data, &conv); err != nil {
		return err
	}

	g.Title = conv.Title
	g.Slug = conv.Slug
	g.Id = conv.Id
	g.Action = conv.Action
	g.ImageActions = map[ImageId]FileAction{}

	for id, image := range conv.ImageActions {
		res := strings.Split(id, "/")
		if len(res) != 2 {
			return errors.New(fmt.Sprintf("Cannot deserialize metadata game action. Id %s is not in the expected serialized format", id))
		}
		g.ImageActions[ImageId{res[0], res[1]}] = image
	}

	return nil
}

func (g *GameAction) Update(n *GameAction) error {
	if (*n).Action == "update" && (*g).Action == "remove" {
		return errors.New("Cannot change a game removal to a game update. This is an impossible situation.")
	}

	if (*n).Action == "remove" || (*n).Action == "add" {
		(*g).Action = (*n).Action
	}

	for id, _ := range (*n).ImageActions {
		(*g).ImageActions[id] = (*n).ImageActions[id]
	}

	return nil
}

func (g *GameAction) IsNoOp() bool {
	return (*g).ActionsLeft() == 0
}

func (g *GameAction) GetImageIds() []ImageId {
	imageIds := make([]ImageId, len((*g).ImageActions))

	idx := 0
	for id, _ := range (*g).ImageActions {
		imageIds[idx] = id
		idx++
	}

	return imageIds
}

func (g *GameAction) CountFileActions() int {
	return len((*g).ImageActions)
}

func (g *GameAction) ActionsLeft() int {
	actionsCount := g.CountFileActions()
	if (*g).Action != "update" {
		actionsCount++
	}
	return actionsCount
}