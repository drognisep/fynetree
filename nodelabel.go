package fynetree

import (
	"fyne.io/fyne"
	"fyne.io/fyne/driver/desktop"
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
	if onTapped := label.node.OnIconTapped; onTapped != nil {
		onTapped(pe)
	}
}

func (label *nodeLabel) TappedSecondary(pe *fyne.PointEvent) {
	if onTapped := label.node.OnTappedSecondary; onTapped != nil {
		onTapped(pe)
	}
}

func (label *nodeLabel) MouseDown(me *desktop.MouseEvent) {
	switch me.Button {
	case desktop.LeftMouseButton:
		if onIconTapped := label.node.OnLabelTapped; onIconTapped != nil {
			onIconTapped(&me.PointEvent)
		}
		break
	case desktop.RightMouseButton:
		if OnTappedSecondary := label.node.OnTappedSecondary; OnTappedSecondary != nil {
			OnTappedSecondary(&me.PointEvent)
		}
	}
}

func (label *nodeLabel) MouseUp(_ *desktop.MouseEvent) {
}
