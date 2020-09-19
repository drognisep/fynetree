package main

import (
	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/theme"
	"github.com/drognisep/fynetree"
	"github.com/drognisep/fynetree/model"
)

func main() {
	myApp := app.New()
	win := myApp.NewWindow("Testing widget")
	win.Resize(fyne.NewSize(300, 200))

	rootNode := fynetree.NewTreeNode(model.NewStaticModel(theme.FolderOpenIcon(), "Projects"))
	task1 := fynetree.NewTreeNode(model.NewStaticModel(nil, "Task 1"))
	_ = rootNode.Append(task1)
	_ = rootNode.Append(fynetree.NewTreeNode(model.NewStaticModel(nil, "Task 2")))
	_ = task1.Append(fynetree.NewTreeNode(model.NewStaticModel(nil, "Sub-Task 1-1")))

	win.SetContent(rootNode)
	win.ShowAndRun()
}
