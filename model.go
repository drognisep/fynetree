package fynetree

import (
	"errors"
	"fmt"
	"fyne.io/fyne"
)

// TreeNode holds a TreeNodeModel's position within the view.
type TreeNode struct {
	parent   *TreeNode
	children []*TreeNode
	model    TreeNodeModel
	expanded bool
}

func NewTreeNode(model TreeNodeModel) *TreeNode {
	newNode := TreeNode{}
	newNode.model = model
	return &newNode
}

func (n *TreeNode) GetParent() *TreeNode {
	return n.parent
}

func (n *TreeNode) IsExpanded() bool {
	return n.expanded
}

func (n *TreeNode) Expand() {
	n.model.BeforeExpand()
	n.expanded = true
}

func (n *TreeNode) Condense() {
	n.expanded = false
	n.model.AfterCondense()
}

func (n *TreeNode) ToggleExpand() {
	if n.expanded {
		n.Condense()
	} else {
		n.Expand()
	}
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
func (n *TreeNode) Append(node *TreeNode) error {
	if node != nil {
		n.children = append(n.children, node)
		node.parent = n
		return nil
	}
	return errors.New("unable to append nil node")
}

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

// ModelChangeListener is called to alert the view that state has changed.
type ModelChangeListener func()

// TreeNodeModel is the interface to user defined data.
type TreeNodeModel interface {
	// GetIconResource should return the user defined icon resource to show in the view, or nil if no icon is needed.
	GetIconResource() *fyne.Resource

	// GetText should return the user defined text to display for this node in the view, or "" if no text is needed.
	GetText() string

	// AddChangeListener is a hook to allow the view to listen for change events and re-render the node.
	// The provided listener is expected to be called by the user when the model view should be refreshed.
	// It is also expected that multiple listeners may be added.
	AddChangeListener(listener ModelChangeListener)

	// BeforeExpand is an event hook called by the view before expansion. The model can take this opportunity to load
	// children in the event that lazy loading is desired.
	BeforeExpand()

	// AfterCondense is an event hook called by the view after condensing. The model can use this to optionally unload
	// resources in child nodes.
	AfterCondense()

	// IsLeaf is checked to determine whether this model supports node expansion.
	IsLeaf() bool
}
