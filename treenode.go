package fynetree

import (
	"fyne.io/fyne"
	"fyne.io/fyne/widget"
	"sync"
)

// NodeEventHandler is a handler function for node events triggered by the view.
type NodeEventHandler func()

// TapEventHandler is a handler function for tap events triggered by the view.
type TapEventHandler func(pe *fyne.PointEvent)

// TreeNode holds a TreeNodeModel's position within the view.
type TreeNode struct {
	widget.BaseWidget
	*nodeList
	model             TreeNodeModel
	expanded          bool
	leaf              bool
	OnBeforeExpand    NodeEventHandler
	OnAfterCondense   NodeEventHandler
	OnTappedSecondary TapEventHandler
	OnIconTapped      TapEventHandler
	OnLabelTapped     TapEventHandler
	OnTapped          TapEventHandler
	OnDoubleTapped    TapEventHandler

	mux    sync.Mutex
	parent *TreeNode
}

// NewTreeNode constructs a tree node with the given model.
func NewTreeNode(model TreeNodeModel) *TreeNode {
	newNode := &TreeNode{}
	InitTreeNode(newNode, model)
	return newNode
}

func NewBranchTreeNode(model TreeNodeModel) *TreeNode {
	return NewTreeNode(model)
}

func NewLeafTreeNode(model TreeNodeModel) *TreeNode {
	leaf := NewTreeNode(model)
	leaf.SetLeaf()
	return leaf
}

// InitTreeNode initializes the given tree node with the given model. If newNode is nil, then a new one will be created.
func InitTreeNode(newNode *TreeNode, model TreeNodeModel) {
	if newNode == nil {
		newNode = &TreeNode{}
	}
	newNode.model = model
	model.SetTreeNode(newNode)
	newNode.initNodeListEvents()
	newNode.OnBeforeExpand = func() {}
	newNode.OnAfterCondense = func() {}
	newNode.OnTappedSecondary = func(pe *fyne.PointEvent) {}
	newNode.leaf = false
	newNode.ExtendBaseWidget(newNode)
}

func (n *TreeNode) initNodeListEvents() {
	n.nodeList = &nodeList{
		OnAfterAddition: func(item fyne.CanvasObject) {
			if item == nil {
				panic("Inserted nil object")
			}
			if i, ok := item.(*TreeNode); ok {
				i.parent = n
				n.Refresh()
			}
		},
		OnAfterRemoval: func(item fyne.CanvasObject) {
			if item != nil {
				if i, ok := item.(*TreeNode); ok {
					i.parent = nil
					n.Refresh()
				}
			}
		},
	}
}

func (n *TreeNode) TappedSecondary(pe *fyne.PointEvent) {
	if n.OnTappedSecondary != nil {
		n.OnTappedSecondary(pe)
	}
}

func (n *TreeNode) DoubleTapped(pe *fyne.PointEvent) {
	if n.OnDoubleTapped != nil {
		n.OnDoubleTapped(pe)
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
	return n.Len()
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
	for _, c := range n.nodeList.Objects {
		c.Show()
	}
}

// Condense condenses the node and triggers the AfterCondense hook in the model if it's a branch and not already condensed.
func (n *TreeNode) Condense() {
	if n.IsBranch() && n.IsExpanded() {
		n.expanded = false
		n.hideChildren()
		n.Refresh()
		if n.OnAfterCondense != nil {
			n.OnAfterCondense()
		}
	}
}

func (n *TreeNode) hideChildren() {
	for _, c := range n.nodeList.Objects {
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
