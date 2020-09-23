package data

import (
	"fmt"
	"fyne.io/fyne"
	"github.com/drognisep/fynetree"
	"testing"
)

var list *nodeList
var modelA *fynetree.StaticNodeModel
var modelB *fynetree.StaticNodeModel
var modelC *fynetree.StaticNodeModel
var modelD *fynetree.StaticNodeModel
var nodeA *fynetree.TreeNode
var nodeB *fynetree.TreeNode
var nodeC *fynetree.TreeNode
var nodeD *fynetree.TreeNode

func setup() {
	list = &nodeList{}
	modelA = fynetree.NewStaticBoundModel(nil, "A")
	modelB = fynetree.NewStaticBoundModel(nil, "B")
	modelC = fynetree.NewStaticBoundModel(nil, "C")
	modelD = fynetree.NewStaticBoundModel(nil, "D")
	nodeA = modelA.Node
	nodeB = modelB.Node
	nodeC = modelC.Node
	nodeD = modelD.Node
}

func TestNodeList_Append(t *testing.T) {
	setup()
	if want, got := 0, list.Len(); got != want {
		t.Errorf("Root node should have no children")
	}
	nodes := []*fynetree.TreeNode{nodeA, nodeB, nodeC}
	for i, n := range nodes {
		fmt.Println("TestNodeList_Append iteration", i)
		err := list.Append(n)
		if err != nil {
			t.Error(err)
		}
		childLen := list.Len()
		if want, got := i+1, childLen; want != got {
			t.Errorf("Root node size should be %d after iteration %d, not %d", want, i, got)
		}
		if want, got := n, list.Objects[childLen-1]; want != got {
			t.Errorf("Node was not inserted at the end of the child Objects")
		}
	}

	err := list.Append(nil)
	if err == nil {
		t.Errorf("Append should have guarded against a nil value")
	}
}

func TestNodeList_RemoveAt(t *testing.T) {
	setup()
	nodes := []*fynetree.TreeNode{nodeA, nodeB, nodeC, nodeD}
	for i, n := range nodes {
		err := list.Append(n)
		if err != nil {
			t.Errorf("Failed to append node %d: %v", i, err)
		}
	}

	if want, got := 4, list.Len(); want != got {
		t.Errorf("All nodes should have been appended to the root node")
	}



	tests := []struct {
		name         string
		removeAtPos  int
		nodeName     string
		expectedNode *fynetree.TreeNode
	}{
		{name: "Remove nodeA", removeAtPos: 0, nodeName: "nodeA", expectedNode: nodeA},
		{name: "Remove nodeC", removeAtPos: 1, nodeName: "nodeC", expectedNode: nodeC},
		{name: "Remove nodeD", removeAtPos: 1, nodeName: "nodeD", expectedNode: nodeD},
		{name: "Remove nodeB", removeAtPos: 0, nodeName: "nodeB", expectedNode: nodeB},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			removed, err := list.RemoveAt(tc.removeAtPos)
			if err != nil {
				t.Fatalf("Failed to remove %s: %v", tc.nodeName, err)
			}
			if want, got := tc.expectedNode, removed; want != got {
				t.Fatalf("Returned node does not equal expected %s", tc.nodeName)
			}
		})
	}

	if want, got := 0, list.Len(); want != got {
		t.Errorf("All nodes should have been removed from the root node")
	}

	_, err := list.RemoveAt(100)
	if err == nil {
		t.Errorf("Error should have been thrown for index out of bounds")
	}
}

func TestNodeList_InsertAt(t *testing.T) {
	setup()
	nodes := []*fynetree.TreeNode{nodeA, nodeB, nodeC}
	nodeInsertionOrder := []*fynetree.TreeNode{nodeB, nodeC, nodeA}
	insertionIndices := []int{0, 0, 1}

	for i, position := range insertionIndices {
		node := nodes[i]
		err := list.InsertAt(position, node)
		if err != nil {
			t.Errorf("Failed to insert %v at position %d", node, position)
		}
	}

	for i, want := range nodeInsertionOrder {
		if got := list.Objects[i]; got != want {
			t.Errorf("Node %v at incorrect position %d, wanted %v", got, i, want)
		}
	}

	err := list.InsertAt(999, nodeD)
	if err == nil {
		t.Errorf("Should have returned an error for index out of bounds")
	}
	err = list.InsertAt(0, nil)
	if err == nil {
		t.Errorf("Should have returned an error for index out of bounds")
	}
}

func TestNodeList_InsertSorted(t *testing.T) {
	setup()
	a := fynetree.NewTreeNode(fynetree.NewStaticModel(nil, "A"))
	b := fynetree.NewTreeNode(fynetree.NewStaticModel(nil, "B"))
	c := fynetree.NewTreeNode(fynetree.NewStaticModel(nil, "c"))
	d := fynetree.NewTreeNode(fynetree.NewStaticModel(nil, "D"))
	empty := fynetree.NewTreeNode(fynetree.NewStaticModel(nil, ""))

	if childLen := list.Len(); childLen != 0 {
		t.Errorf("List should be empty")
	}

	insertionOrder := []*fynetree.TreeNode{b, a, d, c, empty}
	for _, c := range insertionOrder {
		err := list.InsertSorted(c)
		if err != nil {
			t.Fatalf("Failed to insert %#v", c)
		}
	}

	items := list.Objects
	tests := map[string]struct {
		index        int
		expectedNode *fynetree.TreeNode
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
			if items[tc.index] != tc.expectedNode {
				t.Fatalf("%s should be %s", name, tc.nodeName)
			}
		})
	}
}

func TestNodeList_Remove(t *testing.T) {
	setup()
	nodes := []*fynetree.TreeNode{nodeA, nodeB, nodeC}
	for i, n := range nodes {
		err := list.Append(n)
		if err != nil {
			t.Errorf("Error occurred appending node %d", i)
		}
	}

	for i, n := range nodes {
		removedNode, err := list.Remove(n)
		if err != nil {
			t.Errorf("Error occurred removing node %d", i)
		}
		if want, got := n, removedNode; want != got {
			t.Errorf("Wrong node removed from root. Wanted %v, got %v", want, got)
		}
	}

	_, err := list.Remove(nodeD)
	if err == nil {
		t.Errorf("Should have returned an error when attempting to remove non-existent nodeD")
	}

	_, err = list.Remove(nil)
	if err == nil {
		t.Errorf("Should have returned an error when attempting to remove nil node")
	}

	if want, got := 0, list.Len(); want != got {
		t.Errorf("Not all nodes were removed")
	}
}

func TestNodeList_OnAfterAddition(t *testing.T) {
	setup()
	if list.Len() != 0 {
		t.Fatalf("List should be empty")
	}
	var lastNode fyne.CanvasObject
	var numInserted int
	list.OnAfterAddition = func(node fyne.CanvasObject) {
		lastNode = node
		numInserted += 1
	}

	insertTests := []struct {
		testName string
		insertedNode *fynetree.TreeNode
		expectedNumInserted int
		insertFunc func() error
	} {
		{testName: "Append", insertedNode: nodeA, expectedNumInserted: 1, insertFunc: func() error {
			return list.Append(nodeA)
		}},
		{testName: "InsertSorted", insertedNode: nodeB, expectedNumInserted: 2, insertFunc: func() error {
			return list.InsertSorted(nodeB)
		}},
		{testName: "InsertAt", insertedNode: nodeC, expectedNumInserted: 3, insertFunc: func() error {
			return list.InsertAt(0, nodeC)
		}},
	}

	for _, tc := range insertTests {
		t.Run(tc.testName, func(t *testing.T) {
			if err := tc.insertFunc(); err != nil {
				t.Fatalf("Error occurred during %s to list %v", tc.testName, err)
			} else if got, want := lastNode, tc.insertedNode; got != want {
				t.Fatalf("lastNode was not set to expected node")
			} else if got, want := numInserted, tc.expectedNumInserted; got != want {
				t.Fatalf("numInserted was not incremented")
			}
		})
	}

	removalTests := map[string]struct {
		removedNode *fynetree.TreeNode
		removeFunc func() (fyne.CanvasObject, error)
	} {
		"RemoveAt": {removedNode: nodeC, removeFunc: func() (fyne.CanvasObject, error) {
			return list.RemoveAt(0)
		}},
		"Remove": {removedNode: nodeB, removeFunc: func() (fyne.CanvasObject, error) {
			return list.Remove(nodeB)
		}},
	}

	expectedLastNode := lastNode
	expectedNumInserted := numInserted
	for name, tc := range removalTests {
		t.Run(name, func(t *testing.T) {
			removedNode, err := tc.removeFunc()
			if err != nil {
				t.Fatalf("Error occurred removing node: %v", err)
			} else if got, want := removedNode, tc.removedNode; got != want {
				t.Fatalf("Unexpected node removed: %#v", removedNode)
			} else if lastNode != expectedLastNode {
				t.Fatalf("Last node was changed after removal operation")
			} else if numInserted != expectedNumInserted {
				t.Fatalf("Num inserted %d does not match expected %d", numInserted, expectedNumInserted)
			}
		})
	}
}

func TestNodeList_OnAfterRemoval(t *testing.T) {
	setup()
	if list.Len() != 0 {
		t.Fatalf("List should be empty")
	}
	var lastNode fyne.CanvasObject
	var numRemoved int
	list.OnAfterRemoval = func(node fyne.CanvasObject) {
		lastNode = node
		numRemoved += 1
	}

	insertTests := []struct {
		testName           string
		insertedNode       *fynetree.TreeNode
		insertFunc         func() error
	} {
		{testName: "Append", insertedNode: nodeA, insertFunc: func() error {
			return list.Append(nodeA)
		}},
		{testName: "InsertSorted", insertedNode: nodeB, insertFunc: func() error {
			return list.InsertSorted(nodeB)
		}},
		{testName: "InsertAt", insertedNode: nodeC, insertFunc: func() error {
			return list.InsertAt(0, nodeC)
		}},
	}

	expectedLastNode := lastNode
	expectedNumRemoved := numRemoved
	for _, tc := range insertTests {
		t.Run(tc.testName, func(t *testing.T) {
			if err := tc.insertFunc(); err != nil {
				t.Fatalf("Error occurred during %s to list %v", tc.testName, err)
			} else if got, want := lastNode, expectedLastNode; got != want {
				t.Fatalf("lastNode was set unexpectedly")
			} else if got, want := numRemoved, expectedNumRemoved; got != want {
				t.Fatalf("numRemoved was changed")
			}
		})
	}

	removalTests := []struct {
		testName           string
		removedNode        *fynetree.TreeNode
		removeFunc         func() (fyne.CanvasObject, error)
		expectedRemovedNum int
	} {
		{testName: "RemoveAt", removedNode: nodeC, removeFunc: func() (fyne.CanvasObject, error) {
			return list.RemoveAt(0)
		}, expectedRemovedNum: 1},
		{testName: "Remove", removedNode: nodeB, removeFunc: func() (fyne.CanvasObject, error) {
			return list.Remove(nodeB)
		}, expectedRemovedNum: 2},
	}

	for _, tc := range removalTests {
		t.Run(tc.testName, func(t *testing.T) {
			removedNode, err := tc.removeFunc()
			if err != nil {
				t.Fatalf("Error occurred removing node: %v", err)
			} else if got, want := removedNode, tc.removedNode; got != want {
				t.Fatalf("Unexpected node removed: %#v", removedNode)
			} else if lastNode != tc.removedNode {
				t.Fatalf("Last node was changed after removal operation")
			} else if numRemoved != tc.expectedRemovedNum {
				t.Fatalf("Num inserted %d does not match expected %d", numRemoved, expectedNumRemoved)
			}
		})
	}
}
