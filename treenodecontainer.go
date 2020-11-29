package fynetree

import (
	"image/color"
	"sync"

	"fyne.io/fyne"
	"fyne.io/fyne/container"
)

// TreeContainer widget simplifies display of several root tree nodes.
type TreeContainer struct {
	*container.Scroll
	*nodeList
	Background color.Color

	mux           sync.Mutex
	vboxContainer *fyne.Container
}

func NewTreeContainer() *TreeContainer {
	var baseRoots []fyne.CanvasObject
	vboxContainer := container.NewVBox(baseRoots...)
	c := &TreeContainer{
		Scroll:        nil,
		Background:    color.Transparent,
		vboxContainer: vboxContainer,
	}
	c.Scroll = container.NewScroll(vboxContainer)
	c.ExtendBaseWidget(c)
	c.nodeList = &nodeList{
		OnAfterAddition: func(item fyne.CanvasObject) {
			if item == nil {
				panic("Added nil root node")
			}
			if i, ok := item.(*TreeNode); ok {
				i.parent = nil
				c.Refresh()
			}
		},
		OnAfterRemoval: func(item fyne.CanvasObject) {
			if item != nil {
				if i, ok := item.(*TreeNode); ok {
					i.parent = nil
					c.Refresh()
				}
			}
		},
	}

	return c
}

func (t *TreeContainer) NumRoots() int {
	t.mux.Lock()
	defer t.mux.Unlock()
	return t.Len()
}

func (t *TreeContainer) Refresh() {
	t.vboxContainer.Objects = t.nodeList.Objects
	t.vboxContainer.Refresh()
	t.Scroll.Refresh()
}
