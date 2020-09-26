package main

import (
	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
	"github.com/drognisep/fynetree"
	"github.com/drognisep/fynetree/example"
)

func main() {
	myApp := app.New()
	win := myApp.NewWindow("Testing widget")
	win.Resize(fyne.NewSize(640, 480))

	treeContainer := fynetree.NewTreeContainer()
	rootModel := fynetree.NewStaticBoundModel(theme.FolderOpenIcon(), "Tasks")
	notesNode := fynetree.NewTreeNode(fynetree.NewStaticModel(theme.FolderIcon(), "Notes"))
	createPopupFunc := func(msg string) func() {
		return func() {
			dialog.ShowInformation("Hello", msg, win)
		}
	}
	exampleTask := &example.Task{
		Summary:     "Hello!",
		Description: "This is an example Task",
		Menu: fyne.NewMenu("", fyne.NewMenuItem("Say Hello", createPopupFunc("Hello from a popup menu!"))),
	}
	exampleNode := fynetree.NewTreeNode(exampleTask)
	exampleNode.SetLeaf()
	exampleNode.OnTappedSecondary = func(pe *fyne.PointEvent) {
		canvas := fyne.CurrentApp().Driver().CanvasForObject(exampleNode)
		widget.ShowPopUpMenuAtPosition(exampleTask.Menu, canvas, pe.AbsolutePosition)
	}
	exampleNode.OnIconTapped = func(pe *fyne.PointEvent) {createPopupFunc("Hello from icon tapped!")()}
	exampleNode.OnLabelTapped = func(pe *fyne.PointEvent) {createPopupFunc("Hello from label tapped!")()}
	_ = rootModel.Node.Append(exampleNode)
	_ = treeContainer.Append(rootModel.Node)
	_ = treeContainer.Append(notesNode)

	addBtn := widget.NewButton("Add Task", addBtnClicked(rootModel.Node, win))
	btnBox := widget.NewVBox(addBtn)

	split := widget.NewHSplitContainer(treeContainer, fyne.NewContainerWithLayout(
		layout.NewBorderLayout(nil, btnBox, nil, nil),
		btnBox,
		example.NewDetailView(exampleTask),
	))
	split.SetOffset(0.3)

	win.SetContent(split)
	win.ShowAndRun()
}

func addBtnClicked(rootNode *fynetree.TreeNode, window fyne.Window) func() {
	addBtnClicked := func() {
		var summary string
		var desc string

		callback := func(accepted bool) {
			if accepted {
				taskNode := example.NewTaskNode(summary, desc)
				_ = rootNode.InsertSorted(taskNode)
				taskNode.Expand()
				subTask := fynetree.NewTreeNode(fynetree.NewStaticModel(theme.CheckButtonIcon(), "Do this"))
				subTask.SetLeaf()
				_ = taskNode.Append(subTask)
			}
		}

		summaryEntry := widget.NewEntry()
		summaryEntry.OnChanged = func(newSummary string) {
			summary = newSummary
		}
		summaryEntry.PlaceHolder = "Issue summary"
		descEntry := widget.NewMultiLineEntry()
		descEntry.OnChanged = func(newDesc string) {
			desc = newDesc
		}
		descEntry.PlaceHolder = "Optional task description"
		dialog.NewCustomConfirm("Add Task", "Add", "Cancel", fyne.NewContainerWithLayout(
			layout.NewFormLayout(),
			widget.NewLabel("Summary"),
			summaryEntry,
			widget.NewLabel("Description"),
			descEntry,
		), callback, window)
	}
	return addBtnClicked
}
