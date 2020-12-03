package fynetree

import (
	"fyne.io/fyne"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
)

type expandHandle struct {
	widget.Icon

	node *TreeNode
}

func NewExpandHandle(node *TreeNode) *expandHandle {
	handle := &expandHandle{
		node: node,
	}
	handle.ExtendBaseWidget(handle)
	handle.Refresh()
	return handle
}

func (e *expandHandle) Refresh() {
	if e.node != nil {
		if e.node.IsLeaf() {
			e.Hide()
		} else {
			e.Show()
			if e.node.IsExpanded() {
				e.SetResource(theme.MenuDropDownIcon())
			} else {
				e.SetResource(theme.MenuExpandIcon())
			}
		}
	}
}

func (e *expandHandle) Tapped(_ *fyne.PointEvent) {
	e.node.ToggleExpand()
}
