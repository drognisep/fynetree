package fynetree

import (
	"fyne.io/fyne"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
	"github.com/drognisep/fynetree/model"
)

const (
	HierarchyPadding = 32
)

type Tree struct {
	roots []*model.TreeNode
}

func (t *Tree) GetContent() fyne.CanvasObject {
	var rootContent []fyne.CanvasObject
	for _, r := range t.roots {
		rootContent = append(rootContent, nodeContent(r))
	}
	return fyne.NewContainerWithLayout(
		layout.NewVBoxLayout(),
		rootContent...
	)
}

func (t *Tree) SetRoots(roots ...*model.TreeNode) {
	t.roots = roots
}

func nodeContent(node *model.TreeNode) fyne.CanvasObject {
	var nodeObjects []fyne.CanvasObject
	iconResource := node.GetModelIconResource()
	if iconResource != nil {
		nodeObjects = append(nodeObjects, widget.NewIcon(iconResource))
	}
	modelText := node.GetModelText()
	if modelText != "" {
		nodeObjects = append(nodeObjects, widget.NewLabel(modelText))
	}
	var container fyne.CanvasObject
	childNodes := node.GetChildren()
	if len(childNodes) > 0 {
		var childCanvasObjects []fyne.CanvasObject
		if !node.IsLeaf() {
			if node.IsExpanded() {
				nodeObjects = prependExpandedIcon(nodeObjects)
				for _, c := range childNodes {
					childCanvasObjects = append(childCanvasObjects, nodeContent(c))
				}
			} else {
				nodeObjects = prependCondensedIcon(nodeObjects)
			}
		}
		container = fyne.NewContainerWithLayout(
			layout.NewVBoxLayout(),
			fyne.NewContainerWithLayout(layout.NewHBoxLayout(), nodeObjects...),
			widget.NewHBox(hierarchySpacer(), widget.NewVBox(childCanvasObjects...)),
		)
	} else {
		container = fyne.NewContainerWithLayout(
			layout.NewHBoxLayout(),
			nodeObjects...,
		)
	}
	return container
}

func prependCondensedIcon(nodeObjects []fyne.CanvasObject) []fyne.CanvasObject {
	condensedIcon := widget.NewIcon(theme.MenuExpandIcon())
	nodeObjects = append([]fyne.CanvasObject{condensedIcon}, nodeObjects...)
	return nodeObjects
}

func prependExpandedIcon(nodeObjects []fyne.CanvasObject) []fyne.CanvasObject {
	expandedIcon := widget.NewIcon(theme.MenuDropDownIcon())
	nodeObjects = append([]fyne.CanvasObject{expandedIcon}, nodeObjects...)
	return nodeObjects
}

func hierarchySpacer() *layout.Spacer {
	var spacer *layout.Spacer = &layout.Spacer{
		FixHorizontal: true,
	}
	spacer.Resize(fyne.NewSize(HierarchyPadding, 0))
	return spacer
}
