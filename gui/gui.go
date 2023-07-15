package gui

import (
	"fmt"
	"goImgCmp/proc"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

type MainGui struct {
	myApp        fyne.App
	myWindow     fyne.Window
	selectFolder *SelectFolder
	collectData  *CollectData
	showResults  *ShowResults
	inspect      *Inspect
}

func NewMainGui() *MainGui {
	r := new(MainGui)
	r.myApp = app.New()
	r.myWindow = r.myApp.NewWindow("Find double images")

	r.SetContentSelectFolder()
	return r
}

func (this *MainGui) ShowAndRun() {
	this.myWindow.ShowAndRun()
}

func (this *MainGui) SetContentSelectFolder() {
	fmt.Println("UI: SelectFolder")
	if this.selectFolder == nil {
		this.selectFolder = NewSelectFolder(this.myWindow)
		this.selectFolder.OnRun = this.SetContentCollectData
	}
	this.myWindow.SetContent(this.selectFolder.Content)
}

func (this *MainGui) SetContentCollectData(files []string) {
	fmt.Println("UI: CollectData")
	if this.collectData == nil {
		this.collectData = NewCollectData()
		this.collectData.OnAbort = this.SetContentSelectFolder
		this.collectData.OnFinish = this.SetContentShowResults
	}
	this.myWindow.SetContent(this.collectData.Content)
	go this.collectData.Start(files)
}

func (this *MainGui) SetContentShowResults(data *proc.CompareResult) {
	fmt.Println("UI: ShowResults")
	if this.showResults == nil {
		this.showResults = NewShowResults()
		this.showResults.OnBack = this.SetContentSelectFolder
		this.showResults.OnInspect = this.SetContentInspect
	}
	this.showResults.SetData(data)
	this.myWindow.SetContent(this.showResults.Content)
}

func (this *MainGui) SetContentInspect(
	id1 int, id2 int, data *proc.CompareResult) {

	fmt.Println("UI: Inspect")
	if this.inspect == nil {
		this.inspect = NewInspect(&this.myWindow)
		this.inspect.OnReturn = this.SetContentShowResults
	}
	this.inspect.SetData(id1, id2, data)
	this.myWindow.SetContent(this.inspect.Content)
}
