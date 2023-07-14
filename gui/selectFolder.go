package gui

import (
	"sort"

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
	rmBtn   *widget.Button
	folder  []string
	list    *widget.List
	selElem int
}

func NewSelectFolder(window fyne.Window) *SelectFolder {
	r := new(SelectFolder)
	r.window = window
	r.folder = make([]string, 0, 8)

	addBtn := widget.NewButton("Add Folder", r.onAddFolder)
	r.rmBtn = widget.NewButton("Del Folder", r.onDeleteFolder)
	r.rmBtn.Disable()
	runBtn := widget.NewButton("Run", r.onRun)
	vbox := container.NewVBox(
		addBtn,
		r.rmBtn,
		runBtn,
	)

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
	r.list.OnSelected = func(i int) {
		r.selElem = i
		r.rmBtn.Enable()
	}
	r.list.OnUnselected = func(_ int) {
		r.rmBtn.Disable()
	}

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
		if uri == nil {
			return
		}
		newFolder := uri.Path()
		// check if our new folder covers old folder or vice versa
		this.folder = mergeIntoFolderList(this.folder, newFolder)
		sort.Strings(this.folder)
		this.list.UnselectAll()
		this.list.Refresh()
	}, this.window).Show()
}

func mergeIntoFolderList(list []string, elem string) []string {
	for i, v := range list {
		if len(v) > len(elem) {
			if v[:len(elem)] == elem {
				// the old folder is a subfolder
				rmElemList := append(list[:i], list[i+1:]...)
				return mergeIntoFolderList(rmElemList, elem)
			}
		} else {
			if v == elem[:len(v)] {
				// the new folder is already covered
				return list
			}
		}
	}
	return append(list, elem)
}

func (this *SelectFolder) onDeleteFolder() {
	i := this.selElem
	this.folder = append(this.folder[:i], this.folder[i+1:]...)
	this.list.UnselectAll()
	this.list.Refresh()
}

func (this *SelectFolder) onRun() {
	if this.OnRun != nil {
		this.OnRun(this.folder)
	} else {
		fmt.Println("Button: Run")
	}
}
