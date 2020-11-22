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
	label.SetText(text)
	label.ExtendBaseWidget(label)
	return label
}

func (label *nodeLabel) Tapped(pe *fyne.PointEvent) {
	if onTapped := label.node.OnLabelTapped; onTapped != nil {
		onTapped(pe)
	}
}

func (label *nodeLabel) TappedSecondary(pe *fyne.PointEvent) {
	label.node.TappedSecondary(pe)
}

func (label *nodeLabel) DoubleTapped(pe *fyne.PointEvent) {
	label.node.DoubleTapped(pe)
}
