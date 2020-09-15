package fynetree

import (
	"fyne.io/fyne"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
	"github.com/drognisep/fynetree/model"
)

type expandHandle struct {
	widget.Icon

	node *model.TreeNode
}

func NewExpandHandle(node *model.TreeNode) *expandHandle {
	handle := &expandHandle{
		node: node,
	}
	handle.ExtendBaseWidget(handle)
	handle.Refresh()
	return handle
}

func (e *expandHandle) Refresh() {
	if e.node.IsExpanded() {
		e.SetResource(theme.MenuDropDownIcon())
	} else {
		e.SetResource(theme.MenuExpandIcon())
	}
	e.Icon.Refresh()
}

func (e *expandHandle) Tapped(event *fyne.PointEvent) {
	e.node.ToggleExpand()
}
