package fynetree

import (
	"fyne.io/fyne"
	"testing"
)

type testNodeIconState struct {
	nodeIcon *nodeIcon
	treeNode *TreeNode
}

func (state *testNodeIconState) setup() *testNodeIconState {
	state.treeNode = NewStaticBoundModel(nil, "node").Node
	state.nodeIcon = newNodeIcon(state.treeNode, nil)
	return state
}

func (state *testNodeIconState) teardown() {
	state.nodeIcon.node = nil
	state.nodeIcon = nil
	state.treeNode.model.(*StaticNodeModel).Node = nil
	state.treeNode.model = nil
	state.treeNode = nil
}

func TestNodeIcon_Tapped(t *testing.T) {
	state := &testNodeIconState{}
	state.setup()
	defer state.teardown()

	var iconTaps int
	tapIncrement := func(pe *fyne.PointEvent) {
		iconTaps += 1
	}
	state.treeNode.OnIconTapped = tapIncrement
	state.treeNode.OnDoubleTapped = tapIncrement
	state.treeNode.OnTappedSecondary = tapIncrement

	pos := fyne.NewPos(0, 0)
	posEvent := &fyne.PointEvent{
		AbsolutePosition: pos,
		Position:         pos,
	}

	tests := []struct {
		name string
		test func()
		expected int
	} {
		{name: "Tapped", test: func() {state.nodeIcon.Tapped(posEvent)}, expected: 1},
		{name: "TappedSecondary", test: func() {state.nodeIcon.TappedSecondary(posEvent)}, expected: 2},
		{name: "DoubleTapped", test: func() {state.nodeIcon.DoubleTapped(posEvent)}, expected: 3},
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
