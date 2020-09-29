package fynetree

import (
	"fyne.io/fyne"
	"fyne.io/fyne/driver/desktop"
	"fyne.io/fyne/widget"
)

type nodeIcon struct {
	widget.Icon

	node *TreeNode
}

func newNodeIcon(node *TreeNode, resource fyne.Resource) *nodeIcon {
	if node == nil {
		panic("Can't pass nil node to nodeIcon")
	}
	icon := &nodeIcon{
		node: node,
	}
	icon.SetResource(resource)
	icon.ExtendBaseWidget(icon)
	return icon
}

func (icon *nodeIcon) Tapped(pe *fyne.PointEvent) {
	if onTapped := icon.node.OnIconTapped; onTapped != nil {
		onTapped(pe)
	}
}

func (icon *nodeIcon) TappedSecondary(pe *fyne.PointEvent) {
	if onTapped := icon.node.OnTappedSecondary; onTapped != nil {
		onTapped(pe)
	}
}

func (icon *nodeIcon) MouseDown(me *desktop.MouseEvent) {
	switch me.Button {
	case desktop.LeftMouseButton:
		if onIconTapped := icon.node.OnIconTapped; onIconTapped != nil {
			onIconTapped(&me.PointEvent)
		}
		break
	case desktop.RightMouseButton:
		if OnTappedSecondary := icon.node.OnTappedSecondary; OnTappedSecondary != nil {
			OnTappedSecondary(&me.PointEvent)
		}
	}
}

func (icon *nodeIcon) MouseUp(_ *desktop.MouseEvent) {
}
