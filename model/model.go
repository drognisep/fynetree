package model

import (
	"fyne.io/fyne"
)

// TreeNodeModel is the interface to user defined data.
type TreeNodeModel interface {
	// GetIconResource should return the user defined icon Resource to show in the view, or nil if no icon is needed.
	GetIconResource() fyne.Resource

	// GetText should return the user defined Text to display for this node in the view, or "" if no Text is needed.
	GetText() string
}

var _ TreeNodeModel = (*StaticNodeModel)(nil)

type StaticNodeModel struct {
	Resource fyne.Resource
	Text     string
}

func (s *StaticNodeModel) GetIconResource() fyne.Resource {
	return s.Resource
}

func (s *StaticNodeModel) GetText() string {
	return s.Text
}

// NewStaticModel creates a TreeNodeModel with fixed values that never change.
func NewStaticModel(resource fyne.Resource, text string) TreeNodeModel {
	return &StaticNodeModel{
		Resource: resource,
		Text:     text,
	}
}

// NodeEventHandler is a handler function for node events triggered by the view.
type NodeEventHandler func()
