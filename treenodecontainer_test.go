package fynetree

import (
	"testing"
)

var treeContainer *TreeContainer
var rootA *TreeNode
var rootB *TreeNode
var rootC *TreeNode
var rootD *TreeNode

func treeContainerTestSetup() {
	treeContainer = NewTreeContainer()
	rootA = NewTreeNode(NewStaticModel(nil, "A"))
	rootB = NewTreeNode(NewStaticModel(nil, "B"))
	rootC = NewTreeNode(NewStaticModel(nil, "C"))
	rootD = NewTreeNode(NewStaticModel(nil, "D"))
}

func TestTreeNodeContainer_Append(t *testing.T) {
	treeContainerTestSetup()
	if want, got := 0, treeContainer.NumRoots(); got != want {
		t.Errorf("Tree container should have no children")
	}

	tests := map[string]*TreeNode{
		"Appending rootA": rootA,
		"Appending rootB": rootB,
		"Appending rootC": rootC,
	}
	var nilParent *TreeNode = nil
	var i int
	for name, n := range tests {
		t.Run(name, func(t *testing.T) {
			err := treeContainer.Append(n)
			if err != nil {
				t.Fatal(err)
			}
			rootLen := treeContainer.NumRoots()
			if want, got := i+1, rootLen; want != got {
				t.Fatalf("Root node size should be %d after iteration %d, not %d", want, i, got)
			}
			if want, got := nilParent, n.parent; want != got {
				t.Fatalf("Parent for node %d should be %v, not %v", i, want, got)
			}
			if want, got := n, treeContainer.roots[rootLen-1]; want != got {
				t.Fatalf("Node was not inserted at the end of the root list")
			}
			i += 1
		})
	}

	err := treeContainer.Append(nil)
	if err == nil {
		t.Errorf("Append should have guarded against a nil value")
	}
}

func TestTreeNodeContainer_RemoveAt(t *testing.T) {
	treeContainerTestSetup()
	nodes := []*TreeNode{rootA, rootB, rootC, rootD}
	for i, n := range nodes {
		err := treeContainer.Append(n)
		if err != nil {
			t.Errorf("Failed to append node %d: %v", i, err)
		}
	}

	if want, got := 4, len(treeContainer.roots); want != got {
		t.Errorf("All nodes should have been appended to the tree container")
	}

	tests := []struct {
		name         string
		removeAtPos  int
		nodeName     string
		expectedNode *TreeNode
	}{
		{name: "Remove rootA", removeAtPos: 0, nodeName: "rootA", expectedNode: rootA},
		{name: "Remove rootC", removeAtPos: 1, nodeName: "rootC", expectedNode: rootC},
		{name: "Remove rootD", removeAtPos: 1, nodeName: "rootD", expectedNode: rootD},
		{name: "Remove rootB", removeAtPos: 0, nodeName: "rootB", expectedNode: rootB},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			removed, err := treeContainer.RemoveAt(tc.removeAtPos)
			if err != nil {
				t.Fatalf("Failed to remove %s: %v", tc.nodeName, err)
			}
			if want, got := tc.expectedNode, removed; want != got {
				t.Fatalf("Returned node does not equal expected %s", tc.nodeName)
			}
		})
	}

	if want, got := 0, treeContainer.NumRoots(); want != got {
		t.Errorf("All nodes should have been removed from the root node")
	}

	var nilNode *TreeNode
	for _, tc := range tests {
		if want, got := nilNode, tc.expectedNode.parent; want != got {
			t.Errorf("%s did not have its parent pointer reset", tc.nodeName)
		}
	}

	_, err := treeContainer.RemoveAt(100)
	if err == nil {
		t.Errorf("Error should have been thrown for index out of bounds")
	}
}

func TestTreeNodeContainer_InsertAt(t *testing.T) {
	treeContainerTestSetup()
	nodes := []*TreeNode{rootA, rootB, rootC}
	nodeInsertionOrder := []*TreeNode{rootB, rootC, rootA}
	insertionIndices := []int{0, 0, 1}

	for i, position := range insertionIndices {
		node := nodes[i]
		err := treeContainer.InsertAt(position, node)
		if err != nil {
			t.Errorf("Failed to insert %v at position %d", node, position)
		}
	}

	for i, want := range nodeInsertionOrder {
		if got := treeContainer.roots[i]; got != want {
			t.Errorf("Node %v at incorrect position %d, wanted %v", got, i, want)
		}
	}

	err := treeContainer.InsertAt(999, rootD)
	if err == nil {
		t.Errorf("Should have returned an error for index out of bounds")
	}
	err = treeContainer.InsertAt(0, nil)
	if err == nil {
		t.Errorf("Should have returned an error for index out of bounds")
	}
}

func TestTreeNodeContainer_InsertSorted(t *testing.T) {
	treeContainerTestSetup()
	a := NewTreeNode(NewStaticModel(nil, "A"))
	b := NewTreeNode(NewStaticModel(nil, "B"))
	c := NewTreeNode(NewStaticModel(nil, "c"))
	d := NewTreeNode(NewStaticModel(nil, "D"))
	empty := NewTreeNode(NewStaticModel(nil, ""))
	insertOrder := []*TreeNode{b, a, d, c, empty}

	if childLen := treeContainer.NumRoots(); childLen != 0 {
		t.Errorf("Root node should be empty")
	}

	for _, n := range insertOrder {
		err := treeContainer.InsertSorted(n)
		if err != nil {
			t.Fatalf("Error occurred inserting node sorted: %v", err)
		}
	}

	roots := treeContainer.roots
	tests := map[string]struct {
		index        int
		expectedNode *TreeNode
		nodeName     string
	}{
		"First node":  {index: 0, expectedNode: empty, nodeName: "empty"},
		"Second node": {index: 1, expectedNode: a, nodeName: "A"},
		"Third node":  {index: 2, expectedNode: b, nodeName: "B"},
		"Fourth node": {index: 3, expectedNode: c, nodeName: "C"},
		"Fifth node":  {index: 4, expectedNode: d, nodeName: "D"},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			if roots[tc.index] != tc.expectedNode {
				t.Fatalf("%s should be %s", name, tc.nodeName)
			}
		})
	}
}

func TestTreeNodeContainer_Remove(t *testing.T) {
	treeContainerTestSetup()
	nodes := []*TreeNode{rootA, rootB, rootC}
	for i, n := range nodes {
		err := treeContainer.Append(n)
		if err != nil {
			t.Errorf("Error occurred appending node %d", i)
		}
	}

	for i, n := range nodes {
		removedNode, err := treeContainer.Remove(n)
		if err != nil {
			t.Errorf("Error occurred removing node %d", i)
		}
		if want, got := n, removedNode; want != got {
			t.Errorf("Wrong node removed from root. Wanted %v, got %v", want, got)
		}
	}

	_, err := treeContainer.Remove(rootD)
	if err == nil {
		t.Errorf("Should have returned an error when attempting to remove non-existent nodeD")
	}

	_, err = treeContainer.Remove(nil)
	if err == nil {
		t.Errorf("Should have returned an error when attempting to remove nil node")
	}

	if want, got := 0, len(treeContainer.roots); want != got {
		t.Errorf("Not all nodes were removed")
	}
}

func TestTreeNodeContainer_RemoveAll(t *testing.T) {
	treeContainerTestSetup()
	nodes := []*TreeNode{rootA, rootB, rootC}
	for _, n := range nodes {
		err := treeContainer.Append(n)
		if err != nil {
			t.Errorf("Expected Append error to be nil: %v", err)
		}
	}
	err := rootC.Append(rootD)
	if err != nil {
		t.Errorf("Expected Append error to be nil: %v", err)
	}

	treeContainer.RemoveAll()
	if got, want := len(treeContainer.roots), 0; got != want {
		t.Errorf("Root node should be empty, has %d children", got)
	}

	if got, want := len(rootC.children), 0; got != want {
		t.Errorf("Node C should be empty, has %d children", got)
	}
}
