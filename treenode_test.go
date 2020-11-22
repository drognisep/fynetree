package fynetree

import (
	"fmt"
	"fyne.io/fyne/test"
	"testing"
)

var rootNode *TreeNode
var rootModel *StaticNodeModel
var modelA TreeNodeModel
var modelB TreeNodeModel
var modelC TreeNodeModel
var modelD TreeNodeModel
var nodeA *TreeNode
var nodeB *TreeNode
var nodeC *TreeNode
var nodeD *TreeNode

func treeNodeSetup() {
	rootModel = NewStaticBoundModel(nil, "root")
	rootNode = rootModel.Node
	modelA = NewStaticModel(nil, "A")
	modelB = NewStaticModel(nil, "B")
	modelC = NewStaticModel(nil, "C")
	modelD = NewStaticModel(nil, "D")
	nodeA = NewTreeNode(modelA)
	nodeB = NewTreeNode(modelB)
	nodeC = NewTreeNode(modelC)
	nodeD = NewTreeNode(modelD)
}

func TestNewTreeNode(t *testing.T) {
	treeNodeSetup()

	_ = rootNode.Append(nodeA)
	_ = rootNode.Append(nodeB)
	_ = rootNode.Append(nodeC)
	nodeC.OnBeforeExpand = func() {
		if nodeC.NumChildren() == 0 {
			_ = nodeC.Append(nodeD)
		}
	}

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

func TestNewBranchTreeNode(t *testing.T) {
	newNode := NewBranchTreeNode(NewStaticModel(nil, "New branch"))
	if newNode == nil {
		t.Fatal("'newNode' is nil")
	}

	if got := newNode.IsBranch(); got != true {
		t.Fatalf("Wanted IsBranch = true, got %v", got)
	}
}

func TestNewLeafTreeNode(t *testing.T) {
	newNode := NewLeafTreeNode(NewStaticModel(nil, "New leaf"))
	if newNode == nil {
		t.Fatal("'newNode' is nil")
	}

	if got := newNode.IsLeaf(); got != true {
		t.Fatalf("Wanted IsLeaf = true, got %v", got)
	}
}

func TestTreeNode_AddRemove(t *testing.T) {
	treeNodeSetup()
	if nodeA.parent != nil {
		t.Fatalf("Node A is in an invalid initial state")
	}

	if err := rootNode.Append(nodeA); err != nil {
		t.Fatalf("Failed to append node: %v", err)
	} else if nodeA.parent != rootNode {
		t.Fatalf("Node A's parent was not set to the root node")
	}

	removedObject, err := rootNode.Remove(nodeA)
	if err != nil {
		t.Fatalf("Error occurred removing node: %v", err)
	} else if removedObject == nil {
		t.Fatalf("Removed node not returned from remove method")
	} else if removedObject != nodeA {
		t.Fatalf("Unexpected node removed from root node")
	}

	if removedNode, ok := removedObject.(*TreeNode); ok {
		if removedNode.parent != nil {
			t.Fatalf("Removed node's parent pointer was not reset")
		}
	} else {
		t.Fatalf("Removed object is not a TreeNode: %T", removedObject)
	}
}
