package fynetree

import "testing"

var treeContainer *TreeContainer

func containerSetup() {
	treeNodeSetup()
	treeContainer = NewTreeContainer()
}

func TestTreeContainer_AddRemove(t *testing.T) {
	containerSetup()
	if nodeA.parent != nil {
		t.Fatalf("Node A is in an invalid initial state")
	}

	if err := treeContainer.Append(nodeA); err != nil {
		t.Fatalf("Failed to append node: %v", err)
	} else if nodeA.parent != nil {
		t.Fatalf("Node A's parent was set to an actual node")
	}

	removedObject, err := treeContainer.Remove(nodeA)
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
