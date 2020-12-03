package fynetree

import (
	"testing"

	"fyne.io/fyne"
)

type testNodeLabelState struct {
	nodeLabel *nodeLabel
	treeNode  *TreeNode
}

func (state *testNodeLabelState) setup() *testNodeLabelState {
	state.treeNode = NewStaticBoundModel(nil, "node").Node
	state.nodeLabel = newNodeLabel(state.treeNode, "test")
	return state
}

func (state *testNodeLabelState) teardown() {
	state.nodeLabel.node = nil
	state.nodeLabel = nil
	state.treeNode.model.(*StaticNodeModel).Node = nil
	state.treeNode.model = nil
	state.treeNode = nil
}

func TestNodeLabel_Tapped(t *testing.T) {
	state := &testNodeLabelState{}
	state.setup()
	defer state.teardown()

	var iconTaps int
	tapIncrement := func(pe *fyne.PointEvent) {
		iconTaps += 1
	}
	state.treeNode.OnLabelTapped = tapIncrement
	state.treeNode.OnDoubleTapped = tapIncrement
	state.treeNode.OnTappedSecondary = tapIncrement

	pos := fyne.NewPos(0, 0)
	posEvent := &fyne.PointEvent{
		AbsolutePosition: pos,
		Position:         pos,
	}

	tests := []struct {
		name     string
		test     func()
		expected int
	}{
		{name: "Tapped", test: func() { state.nodeLabel.Tapped(posEvent) }, expected: 1},
		{name: "TappedSecondary", test: func() { state.nodeLabel.TappedSecondary(posEvent) }, expected: 2},
		{name: "DoubleTapped", test: func() { state.nodeLabel.DoubleTapped(posEvent) }, expected: 3},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.test()
			if iconTaps != tc.expected {
				t.Fatalf("Expected %d taps, got %d", tc.expected, iconTaps)
			}
		})
	}
}
