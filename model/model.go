package model

import (
	"fyne.io/fyne"
)

// TreeNodeModel is the interface to user defined data.
type TreeNodeModel interface {
	// GetIconResource should return the user defined icon resource to show in the view, or nil if no icon is needed.
	GetIconResource() fyne.Resource

	// GetText should return the user defined text to display for this node in the view, or "" if no text is needed.
	GetText() string
}

var _ TreeNodeModel = (*staticNodeModel)(nil)

type staticNodeModel struct {
	resource fyne.Resource
	text     string
}

func (s *staticNodeModel) GetIconResource() fyne.Resource {
	return s.resource
}

func (s *staticNodeModel) GetText() string {
	return s.text
}

// NewStaticModel creates a TreeNodeModel with fixed values that never change.
func NewStaticModel(resource fyne.Resource, text string) TreeNodeModel {
	return &staticNodeModel{
		resource: resource,
		text:     text,
	}
}

// NodeEventHandler is a handler function for node events triggered by the view.
type NodeEventHandler func()
