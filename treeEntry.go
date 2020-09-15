package fynetree

import (
	"fmt"
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

func (renderer treeEntryRenderer) Layout(container fyne.Size) {
	entry := renderer.entry
	itemBox := entry.itemBox
	childBox := entry.childBox
	minSize := entry.MinSize()
	if !entry.Hidden && minSize.Height != 0 && minSize.Width != 0 {
		itemBoxSize := itemBox.MinSize()
		itemBox.Move(fyne.NewPos(0, 0))
		itemBox.Resize(fyne.NewSize(container.Width, itemBoxSize.Height))
		childBox.Move(fyne.NewPos(itemBoxSize.Height, HierarchyPadding))
		if minSize.Height > itemBoxSize.Height {
			childBox.Resize(fyne.NewSize(container.Width - HierarchyPadding, minSize.Height - itemBoxSize.Height))
		} else {
			childBox.Hide()
		}
	} else {
		entry.Hide()
		itemBox.Hide()
		childBox.Hide()
	}
}

func (renderer treeEntryRenderer) MinSize() fyne.Size {
	return renderer.entry.MinSize()
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
	fmt.Println("TreeEntry.Refresh() called")
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

func (entry *TreeEntry) MinSize() fyne.Size {
	entry.Refresh()
	if entry.itemBox.Hidden {
		return fyne.NewSize(0, 0)
	} else if entry.Node.IsLeaf() || !entry.Node.IsExpanded() {
		return entry.itemBox.MinSize()
	}
	itemBoxSize := entry.itemBox.MinSize()
	childBoxSize := entry.childBox.MinSize()
	runningSize := fyne.NewSize(intMax(itemBoxSize.Width, childBoxSize.Width), itemBoxSize.Height+childBoxSize.Height)
	for _, c := range entry.Node.GetChildren() {
		if childView, ok := (c.View).(*TreeEntry); childView != nil && ok {
			childSize := childView.MinSize()
			runningSize = fyne.NewSize(intMax(childSize.Width, runningSize.Width), childSize.Height + runningSize.Height)
		}
	}
	return runningSize
}
