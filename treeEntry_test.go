package fynetree

import (
	"fyne.io/fyne"
	"fyne.io/fyne/layout"
	"github.com/drognisep/fynetree/model"
	"testing"
	"fyne.io/fyne/test"
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
	entry := NewTreeEntry(rootNode)
	win.SetContent(fyne.NewContainerWithLayout(
		layout.NewVBoxLayout(),
		entry,
	))

	win.ShowAndRun()

	if entry.Hidden {
		t.Errorf("Root entry should not be hidden")
	}
	if entry.itemBox.Hidden {
		t.Errorf("Item box is hidden")
	}

	test.Tap(entry.handle)
	itemPos := entry.itemBox.Position()
	childPos := entry.childBox.Position()

	if offset := childPos.X - itemPos.X; offset != HierarchyPadding {
		t.Errorf("Offset should be %d, not %d", HierarchyPadding, offset)
	}
}
