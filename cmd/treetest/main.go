package main

import (
	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
	"github.com/drognisep/fynetree"
	"github.com/drognisep/fynetree/model"
	"github.com/drognisep/fynetree/util"
	"image/color"
)

func main() {
	myApp := app.New()
	win := myApp.NewWindow("Testing widget")
	win.Resize(fyne.NewSize(640, 480))

	rootNode := fynetree.NewTreeNode(model.NewStaticModel(theme.FolderOpenIcon(), "Tasks"))
	exampleTask := &Task{
		Summary:     "Hello!",
		Description: "This is an example Task",
	}
	exampleNode := fynetree.NewTreeNode(exampleTask)
	exampleNode.SetLeaf()
	_ = rootNode.Append(exampleNode)
	_ = rootNode.Append(fynetree.NewTreeNode(model.NewStaticModel(nil, "Some other task")))
	treeContainer := widget.NewVBox(rootNode)
	scrollContainer := widget.NewScrollContainer(treeContainer)
	split := widget.NewHSplitContainer(scrollContainer, NewDetailView(exampleTask))
	split.SetOffset(0.3)
	win.SetContent(split)
	win.ShowAndRun()
}

var _ model.TreeNodeModel = (*Task)(nil)

type Task struct {
	Summary     string
	Description string
}

func (t *Task) GetIconResource() fyne.Resource {
	return theme.ConfirmIcon()
}

func (t *Task) GetText() string {
	return t.Summary
}

var _ fyne.Widget = (*DetailView)(nil)

type DetailView struct {
	widget.BaseWidget
	Task *Task
}

func NewDetailView(task *Task) *DetailView {
	view := &DetailView{
		Task: task,
	}
	view.ExtendBaseWidget(view)
	return view
}

func (d *DetailView) CreateRenderer() fyne.WidgetRenderer {
	return newDetailViewRenderer(d)
}

var _ fyne.WidgetRenderer = (*detailViewRenderer)(nil)

type detailViewRenderer struct {
	view        *DetailView
	summary     *canvas.Text
	description *canvas.Text
}

func newDetailViewRenderer(view *DetailView) *detailViewRenderer {
	defaultTextSize := float64(fyne.CurrentApp().Settings().Theme().TextSize())
	summary := &canvas.Text{
		Color:     theme.TextColor(),
		Text:      view.Task.Summary,
		TextSize:  int(defaultTextSize * 1.5),
		TextStyle: fyne.TextStyle{
			Bold: true,
		},
	}
	description := canvas.NewText(view.Task.Description, theme.TextColor())
	return &detailViewRenderer{
		view:        view,
		summary:     summary,
		description: description,
	}
}

var spacer = fyne.Size{
	Width:  5,
	Height: 15,
}

func (d *detailViewRenderer) Layout(container fyne.Size) {
	summarySize := d.summary.MinSize()
	d.summary.Move(fyne.NewPos(spacer.Width, 0))
	d.summary.Resize(fyne.NewSize(container.Width, summarySize.Height))

	descSize := d.description.MinSize()
	d.description.Move(fyne.NewPos(spacer.Width, summarySize.Height + spacer.Height))
	d.description.Resize(fyne.NewSize(container.Width, descSize.Height))
}

func (d *detailViewRenderer) MinSize() fyne.Size {
	return util.ColumnMaxSize(d.summary.MinSize(), spacer, d.description.MinSize())
}

func (d *detailViewRenderer) Refresh() {
	d.summary.Text = d.view.Task.Summary
	d.description.Text = d.view.Task.Description
}

func (d *detailViewRenderer) BackgroundColor() color.Color {
	return theme.BackgroundColor()
}

func (d *detailViewRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{d.summary, d.description}
}

func (d *detailViewRenderer) Destroy() {
	d.summary = nil
	d.description = nil
}
