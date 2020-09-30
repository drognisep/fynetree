package fynetree

import (
	"fyne.io/fyne"
	"fyne.io/fyne/widget"
)

type nodeIcon struct {
	widget.Icon

	node *TreeNode
}

func newNodeIcon(node *TreeNode, resource fyne.Resource) *nodeIcon {
	if node == nil {
		panic("Can't pass nil node to nodeLabel")
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
	icon.node.TappedSecondary(pe)
}

func (icon *nodeIcon) DoubleTapped(pe *fyne.PointEvent) {
	icon.node.DoubleTapped(pe)
}
