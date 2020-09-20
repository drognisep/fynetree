package example

import (
	"fyne.io/fyne"
	"fyne.io/fyne/theme"
	"github.com/drognisep/fynetree"
	"github.com/drognisep/fynetree/model"
)

var _ model.TreeNodeModel = (*Task)(nil)

type Task struct {
	Summary     string
	Description string
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
