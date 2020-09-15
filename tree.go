package fynetree

import (
	"fyne.io/fyne"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
	"github.com/drognisep/fynetree/model"
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

var _ model.TreeView = (*TreeItem)(nil)

// Create Tree widget that wraps a set of tree items in a scroll container with a good background

// Create TreeItemHandle widget that extends icon to be tappable

type TreeItem struct {
	Icon *widget.Icon
	Label *widget.Label
	node *model.TreeNode
	// Add tree item handle that is rendered if node IsLeaf() returns false
}

func NewTreeItem(customModel model.TreeNodeModel) *TreeItem {
	return NewTreeItemFromNode(model.NewTreeNode(customModel))
}
func NewTreeItemFromNode(node *model.TreeNode) *TreeItem {
	if node != nil {
		item := &TreeItem{}
		item.node = node
		node.View = item
		item.Icon = widget.NewIcon(node.GetModelIconResource())
		item.Label = widget.NewLabel(node.GetModelText())
		item.Refresh()
		return item
	} else {
		panic("Nil node pointer")
	}
}

func (item *TreeItem) Refresh() {
	node := item.node
	iconResource := node.GetModelIconResource()
	labelText := node.GetModelText()
	if iconResource == nil {
		item.Icon.Hide()
	} else {
		item.Icon.SetResource(iconResource)
	}
	if labelText == "" {
		item.Label.Hide()
	} else {
		item.Label.SetText(labelText)
	}
}

func (item *TreeItem) RefreshAll() {
	item.Refresh()
	if !item.node.IsLeaf() && item.node.IsExpanded() {
		for _, c := range item.node.GetChildren() {
			nodeView := c.View
			if nodeView != nil {
				if childItem, ok := nodeView.(*TreeItem); ok {
					childItem.RefreshAll()
				}
			}
		}
	}
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

type spreader struct {
	FixHorizontal bool
	FixVertical   bool

	size   fyne.Size
	pos    fyne.Position
	hidden bool
}

var _ fyne.CanvasObject = (*spreader)(nil)
// ExpandVertical returns whether or not this spacer expands on the vertical axis
func (s *spreader) ExpandVertical() bool {
	return !s.FixVertical
}

// ExpandHorizontal returns whether or not this spacer expands on the horizontal axis
func (s *spreader) ExpandHorizontal() bool {
	return !s.FixHorizontal
}

// Size returns the current size of this Spacer
func (s *spreader) Size() fyne.Size {
	return s.size
}

// Resize sets a new size for the Spacer - this will be called by the layout
func (s *spreader) Resize(size fyne.Size) {
	s.size = size
}

// Position returns the current position of this Spacer
func (s *spreader) Position() fyne.Position {
	return s.pos
}

// Move sets a new position for the Spacer - this will be called by the layout
func (s *spreader) Move(pos fyne.Position) {
	s.pos = pos
}

// MinSize returns a 0 size as a Spacer can shrink to no actual size
func (s *spreader) MinSize() fyne.Size {
	return s.size
}

// Visible returns true if this spacer should affect the layout
func (s *spreader) Visible() bool {
	return !s.hidden
}

// Show sets the Spacer to be part of the layout calculations
func (s *spreader) Show() {
	s.hidden = false
}

// Hide removes this Spacer from layout calculations
func (s *spreader) Hide() {
	s.hidden = true
}

// Refresh does nothing for a spacer but is part of the CanvasObject definition
func (s *spreader) Refresh() {
}

// NewSpacer returns a spacer object which can fill vertical and horizontal
// space. This is primarily used with a box layout.
func NewSpacer() fyne.CanvasObject {
	return &spreader{}
}