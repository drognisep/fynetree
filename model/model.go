package model

import (
	"errors"
	"fmt"
	"fyne.io/fyne"
)

// TreeNodeModel is the interface to user defined data.
type TreeNodeModel interface {
	// GetIconResource should return the user defined icon resource to show in the view, or nil if no icon is needed.
	GetIconResource() *fyne.Resource

	// GetText should return the user defined text to display for this node in the view, or "" if no text is needed.
	GetText() string
}

// NodeEventHandler is a handler function for node events triggered by the view.
type NodeEventHandler func()

// ModelChangeListener is called to alert the view that state has changed.
type ModelChangeListener func()

// TreeNode holds a TreeNodeModel's position within the view.
type TreeNode struct {
	parent        *TreeNode
	children      []*TreeNode
	model         TreeNodeModel
	expanded      bool
	beforeExpand  NodeEventHandler
	afterCondense NodeEventHandler
	modelChanged  ModelChangeListener
	leaf          bool
}

// NewTreeNode constructs a tree node with the given model.
func NewTreeNode(model TreeNodeModel) *TreeNode {
	newNode := TreeNode{}
	newNode.model = model
	newNode.beforeExpand = func() {}
	newNode.afterCondense = func() {}
	newNode.modelChanged = func() {}
	newNode.leaf = false
	return &newNode
}

// GetParent gets the parent node, or nil if this is a root node.
func (n *TreeNode) GetParent() *TreeNode {
	return n.parent
}

// GetModelIconResource gets the icon for this node.
func (n *TreeNode) GetModelIconResource() *fyne.Resource {
	return n.model.GetIconResource()
}

// GetModelText gets the text for this node.
func (n *TreeNode) GetModelText() string {
	return n.model.GetText()
}

// IsLeaf returns whether this is a leaf node.
func (n *TreeNode) IsLeaf() bool {
	return n.leaf
}

// SetLeaf sets this node to a leaf node.
func (n *TreeNode) SetLeaf() {
	n.leaf = true
}

// SetBranch sets this node to a branch node.
func (n *TreeNode) SetBranch() {
	n.leaf = false
}

// OnBeforeExpand sets the model handler for before the node has been expanded.
func (n *TreeNode) OnBeforeExpand(handler NodeEventHandler) {
	n.beforeExpand = handler
}

// OnAfterCondense sets the model handler for after the node has been condensed.
func (n *TreeNode) OnAfterCondense(handler NodeEventHandler) {
	n.afterCondense = handler
}

// IsExpanded returns whether this node is expanded.
func (n *TreeNode) IsExpanded() bool {
	return n.expanded
}

// Expand expands the node and triggers the BeforeExpand hook in the model if it's not already expanded.
func (n *TreeNode) Expand() {
	if !n.expanded {
		n.beforeExpand()
		n.expanded = true
	}
}

// Condense condenses the node and triggers the AfterCondense hook in the model if it's not already condensed.
func (n *TreeNode) Condense() {
	if n.expanded {
		n.expanded = false
		n.afterCondense()
	}
}

// ToggleExpand toggles the expand state of the node.
func (n *TreeNode) ToggleExpand() {
	if n.expanded {
		n.Condense()
	} else {
		n.Expand()
	}
}

// OnModelChanged sets the view handler for when the model has changed.
func (n *TreeNode) OnModelChanged(handler ModelChangeListener) {
	n.modelChanged = handler
}

// ModelChanged triggered by the model to alert the view that the model has changed.
func (n *TreeNode) ModelChanged() {
	n.modelChanged()
}

// InsertAt a new TreeNode at the given position as a child of this node.
func (n *TreeNode) InsertAt(position int, node *TreeNode) error {
	if node != nil {
		childrenLen := len(n.children)
		if position == childrenLen {
			err := n.Append(node)
			return err
		} else if position == 0 {
			n.children = append([]*TreeNode{node}, n.children...)
		} else if position > 0 && position < childrenLen {
			n.children = append(n.children, nil)
			copy(n.children[(position+1):], n.children[position:])
			n.children[position] = node
		} else {
			return fmt.Errorf("position %d is out of bounds for %d length children", position, childrenLen)
		}
		node.parent = n
		return nil
	}
	return errors.New("unable to insert nil node")
}

// Append adds a node to the end of the list.
func (n *TreeNode) Append(node *TreeNode) error {
	if node != nil {
		n.children = append(n.children, node)
		node.parent = n
		return nil
	}
	return errors.New("unable to append nil node")
}

// Remove the child node at the given position and return it. An error is returned if the index is invalid or the node is not found.
func (n *TreeNode) RemoveAt(position int) (removedNode *TreeNode, err error) {
	childrenLen := len(n.children)
	if position == 0 {
		removedNode = n.children[position]
		n.children = n.children[1:]
	} else if position > 0 && position < (childrenLen-1) {
		removedNode = n.children[position]
		n.children = append(n.children[0:position], n.children[(position+1):]...)
	} else if position < childrenLen {
		removedNode = n.children[position]
		n.children = n.children[:position]
	} else {
		err = fmt.Errorf("position %d is out of bounds for %d length children", position, childrenLen)
		return
	}
	removedNode.parent = nil
	return
}

// Remove searches for the given node to remove and return it if it exists, returns nil and an error otherwise.
func (n *TreeNode) Remove(node *TreeNode) (removedNode *TreeNode, err error) {
	if node != nil {
		for i, existing := range n.children {
			if existing == node {
				return n.RemoveAt(i)
			}
		}
	} else {
		return nil, errors.New("unable to reference nil node")
	}
	return nil, errors.New("unable to locate node")
}
