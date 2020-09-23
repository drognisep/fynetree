package fynetree

import (
	"fyne.io/fyne"
)

// TreeNodeModel is the interface to user defined data.
type TreeNodeModel interface {
	// GetIconResource should return the user defined icon Resource to show in the view, or nil if no icon is needed.
	GetIconResource() fyne.Resource

	// GetText should return the user defined Text to display for this node in the view, or "" if no Text is needed.
	GetText() string

	// SetTreeNode is called by node initialization when the model is bound to the node.
	SetTreeNode(node *TreeNode)
}

var _ TreeNodeModel = (*StaticNodeModel)(nil)

type StaticNodeModel struct {
	Resource fyne.Resource
	Text     string
	Node     *TreeNode
}

func (s *StaticNodeModel) SetTreeNode(node *TreeNode) {
	s.Node = node
}

func (s *StaticNodeModel) GetIconResource() fyne.Resource {
	return s.Resource
}

func (s *StaticNodeModel) GetText() string {
	return s.Text
}

// NewStaticModel creates a TreeNodeModel with fixed values that never change.
func NewStaticModel(resource fyne.Resource, text string) *StaticNodeModel {
	return &StaticNodeModel{
		Resource: resource,
		Text:     text,
	}
}

// NewStaticBoundModel creates a TreeNode and StaticNodeModel at the same time, and returns the bound model.
func NewStaticBoundModel(resource fyne.Resource, text string) *StaticNodeModel {
	model := NewStaticModel(resource, text)
	InitTreeNode(nil, model)
	return model
}
