package example

import (
	"fyne.io/fyne"
	"fyne.io/fyne/theme"
	"github.com/drognisep/fynetree"
)

var _ fynetree.TreeNodeModel = (*Task)(nil)

type Task struct {
	Summary     string
	Description string
	Node        *fynetree.TreeNode
	Menu        *fyne.Menu
}

func (t *Task) SetTreeNode(node *fynetree.TreeNode) {
	t.Node = node
}

func NewTaskNode(summary, description string) *fynetree.TreeNode {
	task := &Task{
		Summary:     summary,
		Description: description,
	}
	node := fynetree.NewTreeNode(task)
	return node
}

func (t *Task) GetIconResource() fyne.Resource {
	return theme.CheckButtonCheckedIcon()
}

func (t *Task) GetText() string {
	return t.Summary
}
