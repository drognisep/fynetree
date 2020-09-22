package fynetree

import (
	"errors"
	"fmt"
	"fyne.io/fyne"
	"fyne.io/fyne/widget"
	"strings"
	"sync"
)

// NodeEventHandler is a handler function for node events triggered by the view.
type NodeEventHandler func()

// TreeNode holds a TreeNodeModel's position within the view.
type TreeNode struct {
	widget.BaseWidget
	model             TreeNodeModel
	expanded          bool
	leaf              bool
	OnBeforeExpand    NodeEventHandler
	OnAfterCondense   NodeEventHandler
	OnTappedSecondary func(pe *fyne.PointEvent)

	mux      sync.Mutex
	parent   *TreeNode
	children []fyne.CanvasObject
}

// NewTreeNode constructs a tree node with the given model.
func NewTreeNode(model TreeNodeModel) *TreeNode {
	newNode := &TreeNode{}
	InitTreeNode(newNode, model)
	return newNode
}

// InitTreeNode initializes the given tree node with the given model. If newNode is nil, then a new one will be created.
func InitTreeNode(newNode *TreeNode, model TreeNodeModel) {
	if newNode == nil {
		newNode = &TreeNode{}
	}
	newNode.model = model
	model.SetTreeNode(newNode)
	newNode.OnBeforeExpand = func() {}
	newNode.OnAfterCondense = func() {}
	newNode.OnTappedSecondary = func(pe *fyne.PointEvent) {}
	newNode.leaf = false
	newNode.ExtendBaseWidget(newNode)
}

func (n *TreeNode) TappedSecondary(pe *fyne.PointEvent) {
	if n.OnTappedSecondary != nil {
		n.OnTappedSecondary(pe)
	}
}

func (n *TreeNode) CreateRenderer() fyne.WidgetRenderer {
	return newTreeEntryRenderer(n)
}

// GetParent gets the parent node, or nil if this is a root node.
func (n *TreeNode) GetParent() *TreeNode {
	return n.parent
}

// NumChildren returns how many child nodes this node has.
func (n *TreeNode) NumChildren() int {
	n.mux.Lock()
	defer n.mux.Unlock()
	return len(n.children)
}

// GetModelIconResource gets the icon for this node.
func (n *TreeNode) GetModelIconResource() fyne.Resource {
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

// IsBranch returns whether this is a branch node.
func (n *TreeNode) IsBranch() bool {
	return !n.leaf
}

// SetLeaf sets this node to a leaf node.
func (n *TreeNode) SetLeaf() {
	n.Condense()
	n.leaf = true
	n.Refresh()
}

// SetBranch sets this node to a branch node.
func (n *TreeNode) SetBranch() {
	n.leaf = false
	n.Refresh()
}

// IsExpanded returns whether this node is expanded.
func (n *TreeNode) IsExpanded() bool {
	return n.expanded
}

// IsCondensed returns whether this node is condensed down and child nodes are not shown.
func (n *TreeNode) IsCondensed() bool {
	return !n.expanded
}

// Expand expands the node and triggers the OnBeforeExpand hook in the model if it's a branch and not already expanded.
func (n *TreeNode) Expand() {
	if n.IsBranch() && n.IsCondensed() {
		if n.OnBeforeExpand != nil {
			n.OnBeforeExpand()
		}
		n.showChildren()
		n.expanded = true
		n.Refresh()
	}
}

func (n *TreeNode) showChildren() {
	for _, c := range n.children {
		c.Show()
	}
}

// Condense condenses the node and triggers the AfterCondense hook in the model if it's a branch and not already condensed.
func (n *TreeNode) Condense() {
	if n.IsBranch() && n.IsExpanded() {
		n.expanded = false
		n.hideChildren()
		if n.OnAfterCondense != nil {
			n.OnAfterCondense()
		}
		n.Refresh()
	}
}

func (n *TreeNode) hideChildren() {
	for _, c := range n.children {
		c.Hide()
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

// InsertAt a new TreeNode at the given position as a child of this node.
func (n *TreeNode) InsertAt(position int, node *TreeNode) error {
	n.mux.Lock()
	if node != nil {
		childrenLen := len(n.children)
		if position == childrenLen {
			n.mux.Unlock()
			err := n.Append(node)
			return err
		} else if position == 0 {
			node.Show()
			n.children = append([]fyne.CanvasObject{node}, n.children...)
		} else if position > 0 && position < childrenLen {
			node.Show()
			n.children = append(n.children, nil)
			copy(n.children[(position+1):], n.children[position:])
			n.children[position] = node
		} else {
			n.mux.Unlock()
			return fmt.Errorf("position %d is out of bounds for %d length children", position, childrenLen)
		}
		node.parent = n
		n.mux.Unlock()
		n.Refresh()
		return nil
	}
	n.mux.Unlock()
	return errors.New("unable to insert nil node")
}

func (n *TreeNode) InsertSorted(node *TreeNode) error {
	n.mux.Lock()
	children := n.children
	for i, c := range children {
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

// Append adds a node to the end of the list.
func (n *TreeNode) Append(node *TreeNode) error {
	if node != nil {
		n.mux.Lock()
		n.children = append(n.children, node)
		node.parent = n
		if n.IsCondensed() {
			node.Hide()
		} else {
			node.Show()
		}
		n.mux.Unlock()
		n.Refresh()
		return nil
	}
	return errors.New("unable to append nil node")
}

// Remove the child node at the given position and return it. An error is returned if the index is invalid or the node is not found.
func (n *TreeNode) RemoveAt(position int) (removedNode fyne.CanvasObject, err error) {
	n.mux.Lock()
	removedNode, err = n.removeAtImpl(position)
	n.mux.Unlock()
	n.Refresh()
	return
}

func (n *TreeNode) removeAtImpl(position int) (removedNode fyne.CanvasObject, err error) {
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
	if treeNode, ok := (removedNode).(*TreeNode); ok {
		treeNode.parent = nil
	}
	return
}

// Remove searches for the given node to remove and return it if it exists, returns nil and an error otherwise.
func (n *TreeNode) Remove(node *TreeNode) (removedNode fyne.CanvasObject, err error) {
	n.mux.Lock()
	if node != nil {
		for i, existing := range n.children {
			if existing == node {
				removedNode, err := n.removeAtImpl(i)
				n.mux.Unlock()
				n.Refresh()
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

// RemoveAll unlinks the node and all child nodes.
func (n *TreeNode) RemoveAll() {
	n.mux.Lock()
	numChildren := len(n.children)
	for i := 0; i < numChildren; i++ {
		node, _ := n.removeAtImpl(0)
		if treeNode, ok := (node).(*TreeNode); ok {
			treeNode.RemoveAll()
		}
	}
	n.mux.Unlock()
	n.Refresh()
}
