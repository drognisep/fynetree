package model

import (
	"fmt"
	"fyne.io/fyne"
	"testing"
)

type mockModel struct {
	changeListeners []ModelChangeListener
	text            string
	leaf            bool
}

func (m mockModel) GetIconResource() *fyne.Resource {
	return nil
}

func (m mockModel) GetText() string {
	return m.text
}

var rootModel TreeNodeModel
var rootNode *TreeNode
var modelA mockModel
var modelB mockModel
var modelC mockModel
var modelD mockModel
var nodeA *TreeNode
var nodeB *TreeNode
var nodeC *TreeNode
var nodeD *TreeNode

func setup() {
	rootModel = mockModel{text: "root"}
	rootNode = NewTreeNode(rootModel)

	modelA = mockModel{text: "A"}
	modelB = mockModel{text: "B"}
	modelC = mockModel{text: "C"}
	modelD = mockModel{text: "D"}
	nodeA = NewTreeNode(modelA)
	nodeB = NewTreeNode(modelB)
	nodeC = NewTreeNode(modelC)
	nodeD = NewTreeNode(modelD)

	for _, m := range []mockModel{modelA, modelB, modelC, modelD} {
		m.leaf = true
	}
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
