package fynetree

import (
	"image/color"
	"sync"

	"fyne.io/fyne"
	"fyne.io/fyne/container"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
)

var _ fyne.Widget = (*TreeContainer)(nil)

// TreeContainer widget simplifies display of several root tree nodes.
type TreeContainer struct {
	widget.BaseWidget
	*nodeList
	Background color.Color

	mux           sync.Mutex
	vboxContainer *fyne.Container
}

func NewTreeContainer() *TreeContainer {
	var baseRoots []fyne.CanvasObject
	vboxContainer := container.NewVBox(baseRoots...)
	c := &TreeContainer{
		Background:    color.Transparent,
		vboxContainer: vboxContainer,
	}
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
}

func (t *TreeContainer) CreateRenderer() fyne.WidgetRenderer {
	return newTreeContainerRenderer(t)
}

var _ fyne.WidgetRenderer = (*treeContainerRenderer)(nil)

type treeContainerRenderer struct {
	scrollContainer *container.Scroll
	treeContainer   *TreeContainer
}

func newTreeContainerRenderer(treeContainer *TreeContainer) *treeContainerRenderer {
	t := &treeContainerRenderer{
		treeContainer: treeContainer,
	}
	t.scrollContainer = container.NewScroll(container.NewVBox(t.treeContainer.Objects...))
	return t
}

func (t *treeContainerRenderer) BackgroundColor() color.Color {
	return color.Transparent
}

func (t *treeContainerRenderer) Destroy() {
	t.scrollContainer = nil
	t.treeContainer = nil
}

func (t *treeContainerRenderer) Layout(_ fyne.Size) {
	y := theme.Padding()
	for _, i := range t.Objects() {
		iSize := i.MinSize()
		i.Resize(iSize)
		i.Move(fyne.NewPos(theme.Padding(), y))
		y = iSize.Height + y
	}
}

func (t *treeContainerRenderer) MinSize() fyne.Size {
	runningSize := fyne.NewSize(0, 0)
	for _, i := range t.Objects() {
		iSize := i.MinSize()
		runningSize = fyne.NewSize(runningSize.Max(iSize).Width, runningSize.Height+iSize.Height)
	}
	return runningSize
}

func (t *treeContainerRenderer) Objects() []fyne.CanvasObject {
	return t.treeContainer.Objects
}

func (t *treeContainerRenderer) Refresh() {
}
