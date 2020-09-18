package fynetree

import (
	"fyne.io/fyne/test"
	"github.com/drognisep/fynetree/model"
	"testing"
)

var rootNode *model.TreeNode
var modelA model.TreeNodeModel
var modelB model.TreeNodeModel
var modelC model.TreeNodeModel
var modelD model.TreeNodeModel
var nodeA *model.TreeNode
var nodeB *model.TreeNode
var nodeC *model.TreeNode
var nodeD *model.TreeNode

func setup() {
	rootNode = model.NewTreeNode(model.NewStaticModel(nil, "root"))
	rootNode.SetBranch()
	modelA = model.NewStaticModel(nil, "A")
	modelB = model.NewStaticModel(nil, "B")
	modelC = model.NewStaticModel(nil, "C")
	modelD = model.NewStaticModel(nil, "D")
	nodeA = model.NewTreeNode(modelA)
	nodeB = model.NewTreeNode(modelB)
	nodeC = model.NewTreeNode(modelC)
	nodeC.SetBranch()
	nodeD = model.NewTreeNode(modelD)

	_ = rootNode.Append(nodeA)
	_ = rootNode.Append(nodeB)
	_ = rootNode.Append(nodeC)
	nodeC.OnBeforeExpand(func() {
		if len(nodeC.GetChildren()) == 0 {
			_ = nodeC.Append(nodeD)
		}
	})
}

func TestNewTreeEntry(t *testing.T) {
	setup()
	testApp := test.NewApp()
	win := testApp.NewWindow("Testing")
	rootEntry := NewTreeEntry(rootNode)
	win.SetContent(rootEntry)

	win.ShowAndRun()

	if rootEntry.Hidden {
		t.Errorf("Root rootEntry should not be hidden")
	}
	if nodeA.View.Visible() {
		t.Errorf("Node A should not be visible yet")
	}

	rootNode.Expand()
	if !nodeA.View.Visible() {
		t.Errorf("Node A is not visible after expanding the root node")
	}
}
