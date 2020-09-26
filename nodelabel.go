package fynetree

import (
	"fyne.io/fyne"
	"fyne.io/fyne/widget"
)

type nodeLabel struct {
	widget.Label

	node *TreeNode
}

func newNodeLabel(node *TreeNode, text string) *nodeLabel {
	if node == nil {
		panic("Can't pass nil node to nodeLabel")
	}
	label := &nodeLabel{
		node: node,
	}
	label.ExtendBaseWidget(label)
	return label
}

func (n *nodeLabel) Tapped(pe *fyne.PointEvent) {
	if onTapped := n.node.OnLabelTapped; onTapped != nil {
		onTapped(pe)
	}
}

func (n *nodeLabel) TappedSecondary(pe *fyne.PointEvent) {
	if onTapped := n.node.OnTappedSecondary; onTapped != nil {
		onTapped(pe)
	}
}
