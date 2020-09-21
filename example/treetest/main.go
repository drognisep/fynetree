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
	"github.com/drognisep/fynetree/model"
)

func main() {
	myApp := app.New()
	win := myApp.NewWindow("Testing widget")
	win.Resize(fyne.NewSize(640, 480))

	treeContainer := fynetree.NewTreeContainer()
	rootNode := fynetree.NewTreeNode(model.NewStaticModel(theme.FolderOpenIcon(), "Tasks"))
	exampleTask := &example.Task{
		Summary:     "Hello!",
		Description: "This is an example Task",
	}
	exampleNode := fynetree.NewTreeNode(exampleTask)
	exampleNode.SetLeaf()
	_ = rootNode.Append(exampleNode)
	_ = treeContainer.Append(rootNode)

	addBtn := widget.NewButton("Add Task", addBtnClicked(rootNode, win))
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
				subTask := fynetree.NewTreeNode(model.NewStaticModel(theme.CheckButtonIcon(), "Do this"))
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
