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
		panic("Can't pass nil node to nodeIcon")
	}
	icon := &nodeIcon{
		node: node,
	}
	icon.ExtendBaseWidget(icon)
	return icon
}

func (n *nodeIcon) Tapped(pe *fyne.PointEvent) {
	if onTapped := n.node.OnIconTapped; onTapped != nil {
		onTapped(pe)
	}
}

func (n *nodeIcon) TappedSecondary(pe *fyne.PointEvent) {
	if onTapped := n.node.OnTappedSecondary; onTapped != nil {
		onTapped(pe)
	}
}
