package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"fmt"
)

type SelectFolder struct {
	Content fyne.CanvasObject
	OnRun   func([]string)
	window  fyne.Window
	folder  []string
	list    *widget.List
}

func NewSelectFolder(window fyne.Window) *SelectFolder {
	r := new(SelectFolder)
	r.window = window
	r.folder = make([]string, 0, 8)
	r.list = widget.NewList(
		func() int {
			return len(r.folder)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("template")
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(r.folder[i])
		},
	)

	addBtn := widget.NewButton("Add Folder", r.onAddFolder)
	runBtn := widget.NewButton("Run", r.onRun)
	vbox := container.NewVBox(
		addBtn,
		runBtn,
	)

	border := container.NewBorder(
		nil,
		vbox,
		nil,
		nil,
		r.list,
	)

	r.Content = border
	return r
}

func (this *SelectFolder) onAddFolder() {
	dialog.NewFolderOpen(func(uri fyne.ListableURI, _ error) {
		if uri != nil {
			this.folder = append(this.folder, uri.Path())
			this.list.Refresh()
		}
	}, this.window).Show()
}

func (this *SelectFolder) onRun() {
	if this.OnRun != nil {
		this.OnRun(this.folder)
	} else {
		fmt.Println("Button: Run")
	}
}
