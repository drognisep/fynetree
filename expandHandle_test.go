package fynetree

import (
	"fyne.io/fyne"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/test"
	"testing"
)

type expandHandleModel struct{}

func (e *expandHandleModel) GetIconResource() fyne.Resource {
	return nil
}

func (e *expandHandleModel) GetText() string {
	return ""
}

func TestNewExpandHandle(t *testing.T) {
	testApp := test.NewApp()
	win := testApp.NewWindow("Some Window")

	node := NewTreeNode(&expandHandleModel{})
	node.SetBranch()
	node.Expand()
	handle := NewExpandHandle(node)

	if node.IsLeaf() || !node.IsExpanded() {
		t.Errorf("Expected node to be a branch and expanded. IsLeaf = %v, IsExpanded = %v", node.IsLeaf(), node.IsExpanded())
	}

	win.SetContent(fyne.NewContainerWithLayout(
		layout.NewHBoxLayout(),
		handle,
	))

	win.ShowAndRun()

	test.Tap(handle)
	if node.IsExpanded() == true {
		t.Errorf("Expected node to now be condensed")
	}
}

func TestNewExpandHandleHiddenWhenLeaf(t *testing.T) {
	testApp := test.NewApp()
	win := testApp.NewWindow("Some Window")

	node := NewTreeNode(&expandHandleModel{})
	node.SetLeaf()
	handle := NewExpandHandle(node)

	if !node.IsLeaf() {
		t.Errorf("Expected node to be a leaf")
	}

	win.SetContent(fyne.NewContainerWithLayout(
		layout.NewHBoxLayout(),
		handle,
	))

	win.ShowAndRun()

	if handle.Hidden == false {
		t.Errorf("Expected handle to be hidden if the node is a leaf")
	}
}
