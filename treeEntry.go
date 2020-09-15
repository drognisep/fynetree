package fynetree

import (
	"fyne.io/fyne"
	"fyne.io/fyne/widget"
	"github.com/drognisep/fynetree/model"
	"image/color"
)

const (
	HierarchyPadding = 32
)

var _ model.TreeView = (*TreeEntry)(nil)

type treeEntryRenderer struct {
	entry          *TreeEntry
}

func (renderer treeEntryRenderer) Layout(size fyne.Size) {
	itemBoxSize := renderer.entry.itemBox.MinSize()
	childBoxSize := renderer.entry.childBox.MinSize()

	renderer.entry.itemBox.Move(fyne.NewPos(0, 0))
	renderer.entry.itemBox.Resize(fyne.NewSize(size.Width, itemBoxSize.Height))

	renderer.entry.childBox.Move(fyne.NewPos(HierarchyPadding, itemBoxSize.Height))
	renderer.entry.childBox.Resize(fyne.NewSize(size.Width - HierarchyPadding, childBoxSize.Height))
}

func (renderer treeEntryRenderer) MinSize() fyne.Size {
	renderer.Refresh()
	itemBoxSize := renderer.entry.itemBox.MinSize()
	childBoxSize := renderer.entry.childBox.MinSize()
	return fyne.NewSize(intMax(itemBoxSize.Width, childBoxSize.Width), itemBoxSize.Height + childBoxSize.Height)
}

func (renderer treeEntryRenderer) Refresh() {
	renderer.entry.updateItemBoxState()
	renderer.entry.updateChildState()
}

func (renderer treeEntryRenderer) BackgroundColor() color.Color {
	return color.Transparent
}

func (renderer *treeEntryRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{renderer.entry.itemBox, renderer.entry.childBox}
}

func (renderer *treeEntryRenderer) Destroy() {
	// Not sure what to do here since each node's view is a static representation of the view model.
	return
}

var _ fyne.Widget = (*TreeEntry)(nil)
var _ fyne.CanvasObject = (*TreeEntry)(nil)
type TreeEntry struct {
	widget.BaseWidget

	Node     *model.TreeNode
	handle   *expandHandle
	icon     *widget.Icon
	label    *widget.Label
	itemBox  *widget.Box
	childBox *widget.Box
}

func NewTreeEntry(node *model.TreeNode) *TreeEntry {
	if node == nil {
		return nil
	}
	handle := NewExpandHandle(node)
	icon := widget.NewIcon(nil)
	label := widget.NewLabel("")
	view := &TreeEntry{
		Node:    node,
		handle:  handle,
		icon:    icon,
		label:   label,
		itemBox: widget.NewHBox(handle, icon, label),
		childBox: widget.NewVBox(),
	}
	node.View = view
	node.OnModelChanged(func() {
		view.Refresh()
	})
	view.Refresh()

	return view
}

func (entry *TreeEntry) CreateRenderer() fyne.WidgetRenderer {
	return &treeEntryRenderer{
		entry: entry,
	}
}

func (entry *TreeEntry) Refresh() {
	entry.updateItemBoxState()
	entry.updateChildState()
}

func (entry *TreeEntry) updateItemBoxState() {
	node := entry.Node

	entry.handle.Refresh()
	// Update icon and label from view model
	iconResource := node.GetModelIconResource()
	labelText := node.GetModelText()

	entry.icon.SetResource(iconResource)
	if iconResource == nil {
		entry.icon.Hide()
	} else {
		entry.icon.Show()
	}
	entry.label.SetText(labelText)
	if labelText == "" {
		entry.label.Hide()
	} else {
		entry.icon.Show()
	}

	parent := node.GetParent()
	if parent != nil && !parent.IsExpanded() {
		entry.itemBox.Hide()
	}
}

func (entry *TreeEntry) updateChildState() {
	node := entry.Node

	if node.IsLeaf() || len(node.GetChildren()) == 0 {
		entry.childBox.Hide()
	} else {
		if node.IsExpanded() {
			entry.childBox = widget.NewVBox()
			entry.childBox.Show()
			childEntries := entry.updateChildren()
			entry.childBox.Children = childEntries
		} else {
			entry.childBox.Hide()
		}
	}
}

func (entry *TreeEntry) updateChildren() []fyne.CanvasObject {
	var childEntries []fyne.CanvasObject
	for _, c := range entry.Node.GetChildren() {
		view := c.View
		childView, ok := view.(*TreeEntry)
		if ok && childView != nil {
			childView.Refresh()
		} else {
			childView = NewTreeEntry(c)
			c.View = childView
		}
		childEntries = append(childEntries, childView)
	}
	return childEntries
}
