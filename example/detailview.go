package example

import (
	"image/color"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
	"github.com/drognisep/fynetree/util"
)

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
		Color:    theme.TextColor(),
		Text:     view.Task.Summary,
		TextSize: int(defaultTextSize * 1.5),
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
	d.description.Move(fyne.NewPos(spacer.Width, summarySize.Height+spacer.Height))
	d.description.Resize(fyne.NewSize(container.Width, descSize.Height))
}

func (d *detailViewRenderer) MinSize() fyne.Size {
	return util.ColumnMinSize(d.summary.MinSize(), spacer, d.description.MinSize())
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
