package fynetree

import (
	"errors"
	"fmt"
	"strings"
	"sync"

	"fyne.io/fyne"
)

type nodeList struct {
	OnAfterAddition func(item fyne.CanvasObject)
	OnAfterRemoval  func(item fyne.CanvasObject)

	mux     sync.Mutex
	Objects []fyne.CanvasObject
}

func (n *nodeList) Len() int {
	return len(n.Objects)
}

// InsertAt a new TreeNode at the given position as a child of this Objects.
func (n *nodeList) InsertAt(position int, node *TreeNode) error {
	n.mux.Lock()
	if node != nil {
		childrenLen := n.Len()
		if position == childrenLen {
			n.mux.Unlock()
			return n.Append(node)
		} else if position == 0 {
			node.Show()
			n.Objects = append([]fyne.CanvasObject{node}, n.Objects...)
		} else if position > 0 && position < childrenLen {
			node.Show()
			n.Objects = append(n.Objects, nil)
			copy(n.Objects[(position+1):], n.Objects[position:])
			n.Objects[position] = node
		} else {
			n.mux.Unlock()
			return fmt.Errorf("position %d is out of bounds for %d length children", position, childrenLen)
		}
		n.mux.Unlock()
		if n.OnAfterAddition != nil {
			n.OnAfterAddition(node)
		}
		return nil
	}
	n.mux.Unlock()
	return errors.New("unable to insert nil node")
}

func (n *nodeList) InsertSorted(node *TreeNode) error {
	n.mux.Lock()
	nodes := n.Objects
	for i, c := range nodes {
		if treeNode, ok := c.(*TreeNode); ok {
			if strings.ToUpper(node.GetModelText()) <= strings.ToUpper(treeNode.GetModelText()) {
				n.mux.Unlock()
				return n.InsertAt(i, node)
			}
		}
	}
	n.mux.Unlock()
	return n.Append(node)
}

// Append adds a node to the end of the Objects.
func (n *nodeList) Append(node *TreeNode) error {
	if node != nil {
		n.mux.Lock()
		n.Objects = append(n.Objects, node)
		n.mux.Unlock()
		if n.OnAfterAddition != nil {
			n.OnAfterAddition(node)
		}
		return nil
	}
	return errors.New("unable to append nil node")
}

// Remove the child node at the given position and return it. An error is returned if the index is invalid or the node is not found.
func (n *nodeList) RemoveAt(position int) (removedNode fyne.CanvasObject, err error) {
	n.mux.Lock()
	removedNode, err = n.removeAtImpl(position)
	n.mux.Unlock()
	return
}

func (n *nodeList) removeAtImpl(position int) (removedNode fyne.CanvasObject, err error) {
	childrenLen := len(n.Objects)
	if position == 0 {
		removedNode = n.Objects[position]
		n.Objects = n.Objects[1:]
	} else if position > 0 && position < (childrenLen-1) {
		removedNode = n.Objects[position]
		n.Objects = append(n.Objects[0:position], n.Objects[(position+1):]...)
	} else if position < childrenLen {
		removedNode = n.Objects[position]
		n.Objects = n.Objects[:position]
	} else {
		err = fmt.Errorf("position %d is out of bounds for %d length children", position, childrenLen)
		return
	}
	if n.OnAfterRemoval != nil {
		n.OnAfterRemoval(removedNode)
	}
	return
}

// Remove searches for the given node to remove and return it if it exists, returns nil and an error otherwise.
func (n *nodeList) Remove(node *TreeNode) (removedNode fyne.CanvasObject, err error) {
	n.mux.Lock()
	if node != nil {
		for i, existing := range n.Objects {
			if existing == node {
				removedNode, err := n.removeAtImpl(i)
				n.mux.Unlock()
				return removedNode, err
			}
		}
	} else {
		n.mux.Unlock()
		return nil, errors.New("unable to reference nil node")
	}
	n.mux.Unlock()
	return nil, errors.New("unable to locate node")
}

// IndexOf returns the index of the given node in the list if it's present, -1 otherwise.
func (n *nodeList) IndexOf(node *TreeNode) int {
	for i, obj := range n.Objects {
		if nodeFound, ok := obj.(*TreeNode); ok {
			if nodeFound == node {
				return i
			}
		}
	}
	return -1
}
