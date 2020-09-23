package fynetree

import (
	"fmt"
	"fyne.io/fyne"
	"testing"
)

var list *NodeList
var listModelA *StaticNodeModel
var listModelB *StaticNodeModel
var listModelC *StaticNodeModel
var listModelD *StaticNodeModel
var listNodeA *TreeNode
var listNodeB *TreeNode
var listNodeC *TreeNode
var listNodeD *TreeNode

func listSetup() {
	list = &NodeList{}
	listModelA = NewStaticBoundModel(nil, "A")
	listModelB = NewStaticBoundModel(nil, "B")
	listModelC = NewStaticBoundModel(nil, "C")
	listModelD = NewStaticBoundModel(nil, "D")
	listNodeA = listModelA.Node
	listNodeB = listModelB.Node
	listNodeC = listModelC.Node
	listNodeD = listModelD.Node
}

func TestNodeList_Append(t *testing.T) {
	listSetup()
	if want, got := 0, list.Len(); got != want {
		t.Errorf("Root node should have no children")
	}
	nodes := []*TreeNode{listNodeA, listNodeB, listNodeC}
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
	listSetup()
	nodes := []*TreeNode{listNodeA, listNodeB, listNodeC, listNodeD}
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
		expectedNode *TreeNode
	}{
		{name: "Remove listNodeA", removeAtPos: 0, nodeName: "listNodeA", expectedNode: listNodeA},
		{name: "Remove listNodeC", removeAtPos: 1, nodeName: "listNodeC", expectedNode: listNodeC},
		{name: "Remove listNodeD", removeAtPos: 1, nodeName: "listNodeD", expectedNode: listNodeD},
		{name: "Remove listNodeB", removeAtPos: 0, nodeName: "listNodeB", expectedNode: listNodeB},
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
	listSetup()
	nodes := []*TreeNode{listNodeA, listNodeB, listNodeC}
	nodeInsertionOrder := []*TreeNode{listNodeB, listNodeC, listNodeA}
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

	err := list.InsertAt(999, listNodeD)
	if err == nil {
		t.Errorf("Should have returned an error for index out of bounds")
	}
	err = list.InsertAt(0, nil)
	if err == nil {
		t.Errorf("Should have returned an error for index out of bounds")
	}
}

func TestNodeList_InsertSorted(t *testing.T) {
	listSetup()
	a := NewTreeNode(NewStaticModel(nil, "A"))
	b := NewTreeNode(NewStaticModel(nil, "B"))
	c := NewTreeNode(NewStaticModel(nil, "c"))
	d := NewTreeNode(NewStaticModel(nil, "D"))
	empty := NewTreeNode(NewStaticModel(nil, ""))

	if childLen := list.Len(); childLen != 0 {
		t.Errorf("List should be empty")
	}

	insertionOrder := []*TreeNode{b, a, d, c, empty}
	for _, c := range insertionOrder {
		err := list.InsertSorted(c)
		if err != nil {
			t.Fatalf("Failed to insert %#v", c)
		}
	}

	items := list.Objects
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
			if items[tc.index] != tc.expectedNode {
				t.Fatalf("%s should be %s", name, tc.nodeName)
			}
		})
	}
}

func TestNodeList_Remove(t *testing.T) {
	listSetup()
	nodes := []*TreeNode{listNodeA, listNodeB, listNodeC}
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

	_, err := list.Remove(listNodeD)
	if err == nil {
		t.Errorf("Should have returned an error when attempting to remove non-existent listNodeD")
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
	listSetup()
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
		insertedNode *TreeNode
		expectedNumInserted int
		insertFunc func() error
	} {
		{testName: "Append", insertedNode: listNodeA, expectedNumInserted: 1, insertFunc: func() error {
			return list.Append(listNodeA)
		}},
		{testName: "InsertSorted", insertedNode: listNodeB, expectedNumInserted: 2, insertFunc: func() error {
			return list.InsertSorted(listNodeB)
		}},
		{testName: "InsertAt", insertedNode: listNodeC, expectedNumInserted: 3, insertFunc: func() error {
			return list.InsertAt(0, listNodeC)
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
		removedNode *TreeNode
		removeFunc func() (fyne.CanvasObject, error)
	} {
		"RemoveAt": {removedNode: listNodeC, removeFunc: func() (fyne.CanvasObject, error) {
			return list.RemoveAt(0)
		}},
		"Remove": {removedNode: listNodeB, removeFunc: func() (fyne.CanvasObject, error) {
			return list.Remove(listNodeB)
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
	listSetup()
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
		insertedNode       *TreeNode
		insertFunc         func() error
	} {
		{testName: "Append", insertedNode: listNodeA, insertFunc: func() error {
			return list.Append(listNodeA)
		}},
		{testName: "InsertSorted", insertedNode: listNodeB, insertFunc: func() error {
			return list.InsertSorted(listNodeB)
		}},
		{testName: "InsertAt", insertedNode: listNodeC, insertFunc: func() error {
			return list.InsertAt(0, listNodeC)
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
		removedNode        *TreeNode
		removeFunc         func() (fyne.CanvasObject, error)
		expectedRemovedNum int
	} {
		{testName: "RemoveAt", removedNode: listNodeC, removeFunc: func() (fyne.CanvasObject, error) {
			return list.RemoveAt(0)
		}, expectedRemovedNum: 1},
		{testName: "Remove", removedNode: listNodeB, removeFunc: func() (fyne.CanvasObject, error) {
			return list.Remove(listNodeB)
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
