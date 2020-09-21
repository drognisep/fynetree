package fynetree

import (
	"fmt"
	"fyne.io/fyne/test"
	"github.com/drognisep/fynetree/model"
	"testing"
)

var rootNode *TreeNode
var modelA model.TreeNodeModel
var modelB model.TreeNodeModel
var modelC model.TreeNodeModel
var modelD model.TreeNodeModel
var nodeA *TreeNode
var nodeB *TreeNode
var nodeC *TreeNode
var nodeD *TreeNode

func setup() {
	rootNode = NewTreeNode(model.NewStaticModel(nil, "root"))
	rootNode.SetBranch()
	modelA = model.NewStaticModel(nil, "A")
	modelB = model.NewStaticModel(nil, "B")
	modelC = model.NewStaticModel(nil, "C")
	modelD = model.NewStaticModel(nil, "D")
	nodeA = NewTreeNode(modelA)
	nodeB = NewTreeNode(modelB)
	nodeC = NewTreeNode(modelC)
	nodeC.SetBranch()
	nodeD = NewTreeNode(modelD)
}

func TestNewTreeEntry(t *testing.T) {
	setup()

	_ = rootNode.Append(nodeA)
	_ = rootNode.Append(nodeB)
	_ = rootNode.Append(nodeC)
	nodeC.OnBeforeExpand(func() {
		if nodeC.NumChildren() == 0 {
			_ = nodeC.Append(nodeD)
		}
	})

	testApp := test.NewApp()
	win := testApp.NewWindow("Testing")
	win.SetContent(rootNode)

	win.ShowAndRun()

	if rootNode.Hidden {
		t.Errorf("Root rootEntry should not be hidden")
	}
	if nodeA.Visible() {
		t.Errorf("Node A should not be visible yet")
	}

	rootNode.Expand()
	if !nodeA.Visible() {
		t.Errorf("Node A is not visible after expanding the root node")
	}

	fmt.Println("Appended Node D to Node C")
	nodeC.Expand()
	if !nodeD.Visible() {
		t.Errorf("Node D should be visible after expanding")
	}

	win.Close()
}

func TestTreeNode_Append(t *testing.T) {
	setup()
	if want, got := 0, len(rootNode.children); got != want {
		t.Errorf("Root node should have no children")
	}
	nodes := []*TreeNode{nodeA, nodeB, nodeC}
	for i, n := range nodes {
		fmt.Println("TestTreeNode_Append iteration", i)
		err := rootNode.Append(n)
		if err != nil {
			t.Error(err)
		}
		childLen := len(rootNode.children)
		if want, got := i+1, childLen; want != got {
			t.Errorf("Root node size should be %d after iteration %d, not %d", want, i, got)
		}
		if want, got := rootNode, n.parent; want != got {
			t.Errorf("Parent for node %d should be %v, not %v", i, want, got)
		}
		if want, got := n, rootNode.children[childLen-1]; want != got {
			t.Errorf("Node was not inserted at the end of the child list")
		}
	}

	err := rootNode.Append(nil)
	if err == nil {
		t.Errorf("Append should have guarded against a nil value")
	}
}

func TestTreeNode_RemoveAt(t *testing.T) {
	setup()
	nodes := []*TreeNode{nodeA, nodeB, nodeC, nodeD}
	for i, n := range nodes {
		err := rootNode.Append(n)
		if err != nil {
			t.Errorf("Failed to append node %d: %v", i, err)
		}
	}

	if want, got := 4, len(rootNode.children); want != got {
		t.Errorf("All nodes should have been appended to the root node")
	}

	removed, err := rootNode.RemoveAt(0)
	if err != nil {
		t.Errorf("Failed to remove nodeA: %v", err)
	}
	if want, got := nodeA, removed; want != got {
		t.Errorf("Returned node does not equal expected nodeA")
	}

	removed, err = rootNode.RemoveAt(1)
	if err != nil {
		t.Errorf("Failed to remove nodeC: %v", err)
	}
	if want, got := nodeC, removed; want != got {
		t.Errorf("Returned node does not equal expected nodeC")
	}

	removed, err = rootNode.RemoveAt(1)
	if err != nil {
		t.Errorf("Failed to remove nodeD: %v", err)
	}
	if want, got := nodeD, removed; want != got {
		t.Errorf("Returned node does not equal expected nodeD")
	}

	removed, err = rootNode.RemoveAt(0)
	if err != nil {
		t.Errorf("Failed to remove nodeB: %v", err)
	}
	if want, got := nodeB, removed; want != got {
		t.Errorf("Returned node does not equal expected nodeB")
	}

	if want, got := 0, len(rootNode.children); want != got {
		t.Errorf("All nodes should have been removed from the root node")
	}

	var nilNode *TreeNode
	for i, n := range nodes {
		if want, got := nilNode, n.parent; want != got {
			t.Errorf("Node %d did not have its parent pointer reset", i)
		}
	}

	_, err = rootNode.RemoveAt(100)
	if err == nil {
		t.Errorf("Error should have been thrown for index out of bounds")
	}
}

func TestTreeNode_InsertAt(t *testing.T) {
	setup()
	nodes := []*TreeNode{nodeA, nodeB, nodeC}
	nodeInsertionOrder := []*TreeNode{nodeB, nodeC, nodeA}
	insertionIndices := []int{0, 0, 1}

	for i, position := range insertionIndices {
		node := nodes[i]
		err := rootNode.InsertAt(position, node)
		if err != nil {
			t.Errorf("Failed to insert %v at position %d", node, position)
		}
	}

	for i, want := range nodeInsertionOrder {
		if got := rootNode.children[i]; got != want {
			t.Errorf("Node %v at incorrect position %d, wanted %v", got, i, want)
		}
	}

	err := rootNode.InsertAt(999, nodeD)
	if err == nil {
		t.Errorf("Should have returned an error for index out of bounds")
	}
	err = rootNode.InsertAt(0, nil)
	if err == nil {
		t.Errorf("Should have returned an error for index out of bounds")
	}
}

func TestTreeNode_InsertSorted(t *testing.T) {
	root := NewTreeNode(model.NewStaticModel(nil, "root"))
	a := NewTreeNode(model.NewStaticModel(nil, "A"))
	b := NewTreeNode(model.NewStaticModel(nil, "B"))
	c := NewTreeNode(model.NewStaticModel(nil, "c"))
	d := NewTreeNode(model.NewStaticModel(nil, "D"))
	empty := NewTreeNode(model.NewStaticModel(nil, ""))

	if childLen := root.NumChildren(); childLen != 0 {
		t.Errorf("Root node should be empty")
	}

	_ = root.InsertSorted(b)
	_ = root.InsertSorted(a)
	_ = root.InsertSorted(d)
	_ = root.InsertSorted(c)
	_ = root.InsertSorted(empty)

	children := root.children
	if children[0] != empty {
		t.Errorf("First node should be empty")
	}
	if children[1] != a {
		t.Errorf("Second node should be A")
	}
	if children[2] != b {
		t.Errorf("Third node should be B")
	}
	if children[3] != c {
		t.Errorf("Fourth node should be C")
	}
	if children[4] != d {
		t.Errorf("Fifth node should be D")
	}
}

func TestTreeNode_Remove(t *testing.T) {
	setup()
	nodes := []*TreeNode{nodeA, nodeB, nodeC}
	for i, n := range nodes {
		err := rootNode.Append(n)
		if err != nil {
			t.Errorf("Error occurred appending node %d", i)
		}
	}

	for i, n := range nodes {
		removedNode, err := rootNode.Remove(n)
		if err != nil {
			t.Errorf("Error occurred removing node %d", i)
		}
		if want, got := n, removedNode; want != got {
			t.Errorf("Wrong node removed from root. Wanted %v, got %v", want, got)
		}
	}

	_, err := rootNode.Remove(nodeD)
	if err == nil {
		t.Errorf("Should have returned an error when attempting to remove non-existent nodeD")
	}

	_, err = rootNode.Remove(nil)
	if err == nil {
		t.Errorf("Should have returned an error when attempting to remove nil node")
	}

	if want, got := 0, len(rootNode.children); want != got {
		t.Errorf("Not all nodes were removed")
	}
}

func TestTreeNode_RemoveAll(t *testing.T) {
	setup()
	nodes := []*TreeNode{nodeA, nodeB, nodeC}
	for _, n := range nodes {
		err := rootNode.Append(n)
		if err != nil {
			t.Errorf("Expected Append error to be nil: %v", err)
		}
	}
	err := nodeC.Append(nodeD)
	if err != nil {
		t.Errorf("Expected Append error to be nil: %v", err)
	}

	rootNode.RemoveAll()
	if got, want := len(rootNode.children), 0; got != want {
		t.Errorf("Root node should be empty, has %d children", got)
	}

	if got, want := len(nodeC.children), 0; got != want {
		t.Errorf("Node C should be empty, has %d children", got)
	}
}
