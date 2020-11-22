# fynetree

This library provides a tree widget implementation that can be used in lieu of one being provided
by the [fyne framework](https://fyne.io) itself.

I really enjoy using the framework and its accompanying tooling, but not being able to show
data in a hierarchical view was a real blocker for me, so I decided to create a widget for this.
Which really helped me understand what's going on under the covers more too. :)

## Task list
- [x] ~~Create base tree node widget with custom layout and dynamic event handling~~
- [x] ~~Create tree node container~~
- [x] ~~Provide InsertSorted method~~
- [x] ~~Handle custom secondary tap menu and logic~~
- [x] ~~Provide icon and text tap event hooks~~
- [ ] Won't do for now ~~Try out some selection model ideas~~
- [x] ~~Possibly create factory methods to create leaf/branch nodes instead of setting leaf
explicitly after creation~~

## How to get it
It's a pure Go library, so using plain `go get github.com/drognisep/fynetree` will get you started.

## How it's organized
The library is meant to follow (more or less) an MVVM structure, borrowing a lot from the base
framework.

### Model
The model is provided by the library consumer, and currently has very few requirements. More
will be added as more interesting functionality comes online.

```golang
// TreeNodeModel is the interface to user defined data.
type TreeNodeModel interface {
	// GetIconResource should return the user defined icon resource to show in the view, or nil if no icon is needed.
	GetIconResource() fyne.Resource

	// GetText should return the user defined text to display for this node in the view, or "" if no text is needed.
	GetText() string
}
```

This optionally returns a `fyne.Resource` which is populated in the icon slot. If `nil` is passed
then no icon will be rendered. The `GetText` method just returns the string that should be shown
for the tree entry.

To make simple things simple, a factory function is provided which creates a static model for
situations where the icon/text will not be changing. `nil` values are accepted for either
argument.

```golang
// NewStaticModel creates a TreeNodeModel with fixed values that never change.
func NewStaticModel(resource fyne.Resource, text string) TreeNodeModel {
	return &StaticNodeModel{
		Resource: resource,
		Text:     text,
	}
}
```

Of course, the static model type is also exposed in case there's a need to create or extend it
directly.

```golang
type StaticNodeModel struct {
	Resource fyne.Resource
	Text     string
}
```

### ViewModel and View
The fyne framework seems to want the `Widget` implementation to be the view model, while the
`WidgetRenderer` defines the actual view logic. I follow this pattern by exposing a `TreeNode`
type which keeps track of expanded/condensed and leaf/branch state, interfaces with the model,
exposes event-based integration points, and manages the set of child nodes.

Nodes can be dynamically added/inserted/removed from any other node, nodes can be
programmatically expanded or condensed, event receivers can be registered, and the icon/text
data can be changed by the model with a node refresh.

```golang
// NodeEventHandler is a handler function for node events triggered by the view.
type NodeEventHandler func()

// TapEventHandler is a handler function for tap events triggered by the view.
type TapEventHandler func(pe *fyne.PointEvent)

// TreeNode holds a TreeNodeModel's position within the view.
type TreeNode struct {
	widget.BaseWidget
	*nodeList
	model             TreeNodeModel
	expanded          bool
	leaf              bool
	OnBeforeExpand    NodeEventHandler
	OnAfterCondense   NodeEventHandler
	OnTappedSecondary TapEventHandler
	OnIconTapped      TapEventHandler
	OnLabelTapped     TapEventHandler

	mux      sync.Mutex
	parent   *TreeNode
}
```

The view is completely defined by the renderer. This keeps view behavior and complexity neatly
hidden behind the widget state facade.

## How to use it
Root tree nodes can be added to a `TreeContainer` to keep them together to be added to a view
once. This type extends the base fyne scroll container to ensure that the tree is flexible
enough to add to most any consumer view without compromising layout expectations.

```golang
// This is an excerpt from the example app. Clone the project to try it out.

/* ... */
// Create a ready-made container
treeContainer := fynetree.NewTreeContainer()
// Used to make a node and model at the same time
rootModel := fynetree.NewStaticBoundModel(theme.FolderOpenIcon(), "Tasks")
// Or created separately with a provided model
notesNode := fynetree.NewTreeNode(fynetree.NewStaticModel(theme.FolderIcon(), "Notes"))
createPopupFunc := func(msg string) func() {
    return func() {
        dialog.ShowInformation("Hello", msg, win)
    }
}
// Task defined elsewhere
exampleTask := &example.Task{
    Summary:     "Hello!",
    Description: "This is an example Task",
    Menu:        fyne.NewMenu("", fyne.NewMenuItem("Say Hello", createPopupFunc("Hello from a popup menu!"))),
}
// Factory methods for creating leaf/branch nodes, can be easily changed later
exampleNode := fynetree.NewLeafTreeNode(exampleTask)
// General event handler
exampleNode.OnTappedSecondary = func(pe *fyne.PointEvent) {
    canvas := fyne.CurrentApp().Driver().CanvasForObject(exampleNode)
    widget.ShowPopUpMenuAtPosition(exampleTask.Menu, canvas, pe.AbsolutePosition)
}
exampleNode.OnDoubleTapped = func(pe *fyne.PointEvent) { createPopupFunc("Hello from node double-tapped")() }
// Icon tap event handler
exampleNode.OnIconTapped = func(pe *fyne.PointEvent) { createPopupFunc("Hello from icon tapped!")() }
// This would block the node double-tapped event
// exampleNode.OnLabelTapped = func(pe *fyne.PointEvent) { createPopupFunc("Hello from label tapped!")() }
_ = rootModel.Node.Append(exampleNode)
_ = treeContainer.Append(rootModel.Node)
_ = treeContainer.Append(notesNode)
/* ... */
```

And that's about it! More features are planned, so check back often. If you see an opportunity
for improvement, please create an issue detailing what you want to see.
