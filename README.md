# fynetree

This library provides a tree widget implementation that can be used in lieu of one being provided
by the [fyne framework](https://fyne.io) itself.

I really enjoy using the framework and its accompanying tooling, but not being able to show
data in a hierarchical view was a real blocker for me, so I decided to create a widget for this.
Which really helped me understand what's going on under the covers more too. :)

## Task list
- [x] ~~Create base tree node widget with custom layout and dynamic event handling~~
- [ ] Create tree node container
- [ ] Provide icon and text tap event hooks
- [ ] Try out some selection model ideas
- [ ] Handle custom secondary tap menu and logic
- [ ] Possibly create factory methods to create leaf/branch nodes instead of setting leaf
explicitly after creation

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
func NewStaticModel(resource fyne.Resource, text string) TreeNodeModel
```

### ViewModel and View
The fyne framework seems to want the `Widget` implementation to be the view model, while the
`WidgetRenderer` defines the actual view logic. I follow this pattern by exposing a `TreeNode`
type which keeps track of expanded/condensed and leaf/branch state, interfaces with the model,
exposes event-based integration points, and manages the set of child nodes.

```golang
type TreeNode struct {
	widget.BaseWidget
	model         model.TreeNodeModel
	expanded      bool
	beforeExpand  model.NodeEventHandler
	afterCondense model.NodeEventHandler
	leaf          bool

	mux      sync.Mutex
	parent   *TreeNode
	children []fyne.CanvasObject
}

// Defined in the model package

// NodeEventHandler is a handler function for node events triggered by the view.
type NodeEventHandler func()
```

The view is completely defined by the renderer. This keeps view behavior and complexity neatly
hidden behind the widget state facade.

## How to use it
Root tree nodes can be simply wrapped in a `NewVBox` to present the view one would expect. I
plan on creating a tree node container to make this more straight forward.

```golang
// This is an excerpt from the example app. Clone the project to try it out.

/* ... */
rootNode := fynetree.NewTreeNode(model.NewStaticModel(theme.FolderOpenIcon(), "Tasks"))

// This implements model.TreeNodeModel
exampleTask := &Task{
    Summary:     "Hello!",
    Description: "This is an example Task",
}

exampleNode := fynetree.NewTreeNode(exampleTask)
exampleNode.SetLeaf() // Branches by default
_ = rootNode.Append(exampleNode)
treeContainer := widget.NewVBox(rootNode)
scrollContainer := widget.NewScrollContainer(treeContainer)
/* ... */
```

Nodes can be dynamically added/inserted/removed from any other node, nodes can be
programmatically expanded or condensed, event receivers can be registered, and the icon/text
data can be changed by the model with a node refresh.

And that's about it! More features are planned, so check back often. If you see an opportunity
for improvement, please create an issue detailing what you want to see.
