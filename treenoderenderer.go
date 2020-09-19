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
	node       *TreeNode
	handle     *expandHandle
	icon       *widget.Icon
	label      *widget.Label
	childBox   *widget.Box
	childNodes []fyne.CanvasObject
}

func newTreeEntryRenderer(node *TreeNode) fyne.WidgetRenderer {
	handle := NewExpandHandle(node)
	icon := widget.NewIcon(node.GetModelIconResource())
	label := widget.NewLabel(node.GetModelText())
	var childNodes []fyne.CanvasObject
	for _, c := range node.GetChildren() {
		childNodes = append(childNodes, c)
	}
	return &treeEntryRenderer{
		node:     node,
		handle:   handle,
		icon:     icon,
		label:    label,
		childBox: widget.NewVBox(childNodes...),
		childNodes: childNodes,
	}
}

func (renderer treeEntryRenderer) Layout(container fyne.Size) {
	node := renderer.node
	itemsHeight := renderer.entryItemsMinSize().Height
	handle := renderer.handle
	handleSize := handle.MinSize()
	var handleWidth int
	if handle.Visible() {
		handleWidth = handleSize.Width
		handle.Move(fyne.NewPos(0, 0))
		handle.Resize(fyne.NewSize(handleWidth, itemsHeight))
	} else {
		handleWidth = 0
	}
	icon := renderer.icon
	iconSize := icon.MinSize()
	var iconWidth int
	if icon.Visible() {
		iconWidth = iconSize.Width
		icon.Move(fyne.NewPos(handleWidth, 0))
		icon.Resize(fyne.NewSize(iconWidth, itemsHeight))
	} else {
		iconWidth = 0
	}
	label := renderer.label
	var labelWidth int
	if label.Visible() {
		labelWidth = container.Width - handleWidth - iconWidth
		label.Move(fyne.NewPos(handleWidth+iconWidth, 0))
		label.Resize(fyne.NewSize(labelWidth, itemsHeight))
	} else {
		labelWidth = 0
	}
	childBox := renderer.childBox
	if node.IsBranch() && node.IsExpanded() {
		childBox.Show()
		childBoxSize := childBox.MinSize()
		childBox.Move(fyne.NewPos(HierarchyPadding, itemsHeight))
		childBox.Resize(fyne.NewSize(container.Width-HierarchyPadding, childBoxSize.Height))
	} else {
		childBox.Hide()
	}
}

func (renderer treeEntryRenderer) MinSize() fyne.Size {
	entryItemsSize := renderer.entryItemsMinSize()
	childBoxSize := renderer.childBox.MinSize()
	return fyne.NewSize(intMax(entryItemsSize.Width, childBoxSize.Width+HierarchyPadding), entryItemsSize.Height+childBoxSize.Height)
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
	renderer.recreateChildBox()
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

func (renderer treeEntryRenderer) recreateChildBox() {
	node := renderer.node
	var childObjects []fyne.CanvasObject
	for _, c := range node.GetChildren() {
		childObjects = append(childObjects, c)
	}
	renderer.childBox.Children = childObjects
	renderer.childBox.Refresh()
}

func (renderer treeEntryRenderer) BackgroundColor() color.Color {
	return color.Transparent
}

func (renderer *treeEntryRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{renderer.handle, renderer.icon, renderer.label, renderer.childBox}
}

func (renderer *treeEntryRenderer) Destroy() {
	renderer.handle.node = nil
	renderer.handle = nil
	renderer.icon = nil
	renderer.label = nil
	renderer.childBox = nil
	renderer.node = nil
}

var _ fyne.Widget = (*TreeNode)(nil)
var _ fyne.CanvasObject = (*TreeNode)(nil)
