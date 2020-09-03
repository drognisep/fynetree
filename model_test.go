package fynetree

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

func (m mockModel) AddChangeListener(listener ModelChangeListener) {
	m.changeListeners = append(m.changeListeners, listener)
}

func (m mockModel) BeforeExpand() {}

func (m mockModel) AfterCondense() {}

func (m mockModel) IsLeaf() bool {
	return m.leaf
}

var rootModel TreeNodeModel
var rootNode *TreeNode
var modelA mockModel
var modelB mockModel
var modelC mockModel
var nodeA *TreeNode
var nodeB *TreeNode
var nodeC *TreeNode

func setup() {
	rootModel = mockModel{text: "root"}
	rootNode = NewTreeNode(rootModel)

	modelA = mockModel{text: "A"}
	modelB = mockModel{text: "B"}
	modelC = mockModel{text: "C"}
	nodeA = NewTreeNode(modelA)
	nodeB = NewTreeNode(modelB)
	nodeC = NewTreeNode(modelC)

	for _, m := range []mockModel{modelA, modelB, modelC} {
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
}

func TestTreeNode_RemoveAt(t *testing.T) {
	setup()
	nodes := []*TreeNode{nodeA, nodeB, nodeC}
	for i, n := range nodes {
		err := rootNode.Append(n)
		if err != nil {
			t.Errorf("Failed to append node %d: %v", i, err)
		}
	}

	if want, got := 3, len(rootNode.children); want != got {
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
}
