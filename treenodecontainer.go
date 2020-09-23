package fynetree

import (
	"fyne.io/fyne"
	"fyne.io/fyne/widget"
	"image/color"
	"sync"
)

// TreeContainer widget simplifies display of several root tree nodes.
type TreeContainer struct {
	widget.ScrollContainer
	*NodeList
	Background color.Color

	mux           sync.Mutex
	vboxContainer *widget.Box
}

func NewTreeContainer() *TreeContainer {
	var baseRoots []fyne.CanvasObject
	vboxContainer := widget.NewVBox(baseRoots...)
	container := &TreeContainer{
		Background:    color.Transparent,
		vboxContainer: vboxContainer,
	}
	container.ExtendBaseWidget(container)
	container.ScrollContainer.Content = vboxContainer
	container.NodeList = &NodeList{
		OnAfterAddition: func(item fyne.CanvasObject) {
			if item == nil {
				panic("Added nil root node")
			}
			if i, ok := item.(*TreeNode); ok {
				i.parent = nil
				container.Refresh()
			}
		},
		OnAfterRemoval:  func(item fyne.CanvasObject) {
			if item != nil {
				if i, ok := item.(*TreeNode); ok {
					i.parent = nil
					container.Refresh()
				}
			}
		},
	}

	return container
}

func (t *TreeContainer) NumRoots() int {
	t.mux.Lock()
	defer t.mux.Unlock()
	return t.Len()
}

func (t *TreeContainer) Refresh() {
	t.vboxContainer.Children = t.NodeList.Objects
	t.vboxContainer.Refresh()
	t.ScrollContainer.Refresh()
}
