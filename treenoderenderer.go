package fynetree

import (
	"fyne.io/fyne"
	"fyne.io/fyne/widget"
	"image/color"
)

const (
	HierarchyPadding = 24
)

type treeEntryRenderer struct {
	node   *TreeNode
	handle *expandHandle
	icon   *widget.Icon
	label  *widget.Label
}

func newTreeEntryRenderer(node *TreeNode) fyne.WidgetRenderer {
	handle := NewExpandHandle(node)
	icon := widget.NewIcon(node.GetModelIconResource())
	label := widget.NewLabel(node.GetModelText())
	return &treeEntryRenderer{
		node:   node,
		handle: handle,
		icon:   icon,
		label:  label,
	}
}

func (renderer treeEntryRenderer) Layout(container fyne.Size) {
	node := renderer.node
	itemsHeight := renderer.entryItemsMinSize().Height
	handle := renderer.handle
	handleSize := handle.MinSize()
	handleWidth := handleSize.Width
	handle.Move(fyne.NewPos(0, 0))
	handle.Resize(fyne.NewSize(handleWidth, itemsHeight))
	icon := renderer.icon
	iconSize := icon.MinSize()
	var iconWidth int
	if icon.Resource != nil {
		iconWidth = iconSize.Width
		icon.Move(fyne.NewPos(handleWidth, 0))
		icon.Resize(fyne.NewSize(iconWidth, itemsHeight))
	} else {
		iconWidth = 0
	}
	label := renderer.label
	var labelWidth int
	if label.Text != "" {
		labelWidth = container.Width - handleWidth - iconWidth
		label.Move(fyne.NewPos(handleWidth+iconWidth, 0))
		label.Resize(fyne.NewSize(labelWidth, itemsHeight))
	} else {
		labelWidth = 0
	}
	if node.IsBranch() && node.IsExpanded() {
		var runningY int = itemsHeight
		for _, c := range node.children {
			cSize := c.MinSize()
			c.Move(fyne.NewPos(HierarchyPadding, runningY))
			c.Resize(fyne.NewSize(container.Width-HierarchyPadding, cSize.Height))
			runningY += cSize.Height
			c.Show()
		}
	} else {
		for _, c := range node.children {
			c.Hide()
		}
	}
}

func (renderer treeEntryRenderer) MinSize() fyne.Size {
	entryItemsSize := renderer.entryItemsMinSize()
	var childrenSize fyne.Size
	for _, c := range renderer.node.children {
		if c.Visible() {
			childSize := c.MinSize()
			childrenSize = fyne.Size{
				Width:  intMax(childrenSize.Width, childSize.Width),
				Height: childrenSize.Height + childSize.Height,
			}
		}
	}
	return fyne.NewSize(intMax(entryItemsSize.Width, childrenSize.Width+HierarchyPadding), entryItemsSize.Height+childrenSize.Height)
}

func (renderer treeEntryRenderer) entryItemsMinSize() fyne.Size {
	handleSize := renderer.handle.MinSize()
	iconSize := renderer.icon.MinSize()
	labelSize := renderer.label.MinSize()
	entryItemsSize := fyne.NewSize(handleSize.Width+iconSize.Width+labelSize.Width,
		intMax(handleSize.Height, iconSize.Height, labelSize.Height))
	return entryItemsSize
}

func (renderer treeEntryRenderer) Refresh() {
	renderer.updateItemBoxState()
}

func (renderer *treeEntryRenderer) updateItemBoxState() {
	node := renderer.node

	renderer.handle.Refresh()
	// Update icon and label from view model
	iconResource := node.GetModelIconResource()
	labelText := node.GetModelText()

	renderer.icon.SetResource(iconResource)
	if iconResource == nil {
		renderer.icon.Hide()
	} else {
		renderer.icon.Show()
	}
	renderer.label.SetText(labelText)
	if labelText == "" {
		renderer.label.Hide()
	} else {
		renderer.icon.Show()
	}
}

func (renderer treeEntryRenderer) BackgroundColor() color.Color {
	return color.Transparent
}

func (renderer *treeEntryRenderer) Objects() []fyne.CanvasObject {
	return append([]fyne.CanvasObject{renderer.handle, renderer.icon, renderer.label}, renderer.node.children...)
}

func (renderer *treeEntryRenderer) Destroy() {
	renderer.handle.node = nil
	renderer.handle = nil
	renderer.icon = nil
	renderer.label = nil
	renderer.node = nil
}

var _ fyne.Widget = (*TreeNode)(nil)
var _ fyne.CanvasObject = (*TreeNode)(nil)
