package fynetree

import (
	"errors"
	"fmt"
	"fyne.io/fyne"
	"fyne.io/fyne/widget"
	"image/color"
	"strings"
	"sync"
)

// TreeContainer widget simplifies display of several root tree nodes.
type TreeContainer struct {
	widget.ScrollContainer
	Background color.Color

	mux           sync.Mutex
	roots         []fyne.CanvasObject
	vboxContainer *widget.Box
}

func NewTreeContainer() *TreeContainer {
	var baseRoots []fyne.CanvasObject
	vboxContainer := widget.NewVBox(baseRoots...)
	container := &TreeContainer{
		Background:    color.Transparent,
		roots:         baseRoots,
		vboxContainer: vboxContainer,
	}
	container.ExtendBaseWidget(container)
	container.ScrollContainer.Content = vboxContainer
	return container
}

func (t *TreeContainer) NumRoots() int {
	return len(t.roots)
}

// InsertAt a new TreeNode at the given position as a child of this node.
func (t *TreeContainer) InsertAt(position int, node *TreeNode) error {
	t.mux.Lock()
	if node != nil {
		childrenLen := len(t.roots)
		if position == childrenLen {
			t.mux.Unlock()
			err := t.Append(node)
			return err
		} else if position == 0 {
			node.Show()
			t.roots = append([]fyne.CanvasObject{node}, t.roots...)
		} else if position > 0 && position < childrenLen {
			node.Show()
			t.roots = append(t.roots, nil)
			copy(t.roots[(position+1):], t.roots[position:])
			t.roots[position] = node
		} else {
			t.mux.Unlock()
			return fmt.Errorf("position %d is out of bounds for %d length children", position, childrenLen)
		}
		node.parent = nil
		t.mux.Unlock()
		t.Refresh()
		return nil
	}
	t.mux.Unlock()
	return errors.New("unable to insert nil node")
}

// InsertSorted inserts a root node sorted by label.
func (t *TreeContainer) InsertSorted(node *TreeNode) error {
	t.mux.Lock()
	roots := t.roots
	for i, c := range roots {
		if treeNode, ok := c.(*TreeNode); ok {
			if strings.ToUpper(node.GetModelText()) <= strings.ToUpper(treeNode.GetModelText()) {
				t.mux.Unlock()
				return t.InsertAt(i, node)
			}
		}
	}
	t.mux.Unlock()
	return t.Append(node)
}

// Append adds a node to the end of the list.
func (t *TreeContainer) Append(node *TreeNode) error {
	if node != nil {
		t.mux.Lock()
		t.roots = append(t.roots, node)
		node.parent = nil
		node.Show()
		t.mux.Unlock()
		t.Refresh()
		return nil
	}
	return errors.New("unable to append nil node")
}

// Remove the child node at the given position and return it. An error is returned if the index is invalid or the node is not found.
func (t *TreeContainer) RemoveAt(position int) (removedNode fyne.CanvasObject, err error) {
	t.mux.Lock()
	removedNode, err = t.removeAtImpl(position)
	t.mux.Unlock()
	t.Refresh()
	return
}

func (t *TreeContainer) removeAtImpl(position int) (removedNode fyne.CanvasObject, err error) {
	childrenLen := len(t.roots)
	if position == 0 {
		removedNode = t.roots[position]
		t.roots = t.roots[1:]
	} else if position > 0 && position < (childrenLen-1) {
		removedNode = t.roots[position]
		t.roots = append(t.roots[0:position], t.roots[(position+1):]...)
	} else if position < childrenLen {
		removedNode = t.roots[position]
		t.roots = t.roots[:position]
	} else {
		err = fmt.Errorf("position %d is out of bounds for %d length children", position, childrenLen)
		return
	}
	if treeNode, ok := (removedNode).(*TreeNode); ok {
		treeNode.parent = nil
	}
	return
}

// Remove searches for the given node to remove and return it if it exists, returns nil and an error otherwise.
func (t *TreeContainer) Remove(node *TreeNode) (removedNode fyne.CanvasObject, err error) {
	t.mux.Lock()
	if node != nil {
		for i, existing := range t.roots {
			if existing == node {
				removedNode, err := t.removeAtImpl(i)
				t.mux.Unlock()
				t.Refresh()
				return removedNode, err
			}
		}
	} else {
		t.mux.Unlock()
		return nil, errors.New("unable to reference nil node")
	}
	t.mux.Unlock()
	return nil, errors.New("unable to locate node")
}

// RemoveAll unlinks the node and all child nodes.
func (t *TreeContainer) RemoveAll() {
	t.mux.Lock()
	numChildren := len(t.roots)
	for i := 0; i < numChildren; i++ {
		node, _ := t.removeAtImpl(0)
		if treeNode, ok := (node).(*TreeNode); ok {
			treeNode.RemoveAll()
		}
	}
	t.mux.Unlock()
	t.Refresh()
}

func (t *TreeContainer) Refresh() {
	t.vboxContainer.Children = t.roots
	t.vboxContainer.Refresh()
	t.ScrollContainer.Refresh()
}
