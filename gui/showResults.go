package gui

import (
	"fmt"
	"goImgCmp/proc"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type ShowResults struct {
	Content   fyne.CanvasObject
	OnBack    func()
	OnInspect func(int, int, *proc.CompareResult)
	dataRaw   *proc.CompareResult
	data      []entry
	list      *widget.List
	selItem   int
	filter    binding.Float
}

type entry struct {
	id1  int
	id2  int
	text string
}

func NewShowResults() *ShowResults {
	r := new(ShowResults)

	// Initialize Data
	r.data = make([]entry, 0)
	r.filter = binding.NewFloat()

	// Create UI Elements
	backBtn := widget.NewButton("Back", r.onBack)
	inspectBtn := widget.NewButton("Inspect", r.onInspect)
	inspectBtn.Disable()

	slider := widget.NewSliderWithData(0, 1, r.filter)
	slider.Step = 0.01
	r.list = widget.NewList(
		func() int {
			return len(r.data)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("template\njob1\njob2")
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(r.data[i].text)
		},
	)
	r.list.OnSelected = func(i widget.ListItemID) {
		inspectBtn.Enable()
		r.selItem = i
	}
	r.list.OnUnselected = func(_ widget.ListItemID) {
		inspectBtn.Disable()
	}

	bottom := container.New(layout.NewGridLayoutWithColumns(3),
		inspectBtn, widget.NewLabel(""), backBtn)
	r.Content = container.NewBorder(
		slider,
		bottom,
		nil,
		nil,
		r.list,
	)

	// Attach Listener
	r.filter.AddListener(binding.NewDataListener(func() {
		r.SetData(r.dataRaw)
	}))

	return r
}

func (this *ShowResults) onBack() {
	if this.OnBack != nil {
		this.OnBack()
	}
}

func (this *ShowResults) onInspect() {
	d := this.data[this.selItem]
	if this.OnInspect != nil {
		this.OnInspect(d.id1, d.id2, this.dataRaw)
	} else {
		fmt.Println("OnInspect(", d.id1, ", ", d.id2, ",...)")
	}
}

func (this *ShowResults) SetData(data *proc.CompareResult) {
	this.dataRaw = data

	filterVal, err := this.filter.Get()
	if err != nil {
		filterVal = 0
	}

	this.list.UnselectAll()
	this.data = make([]entry, 0)
	if data == nil {
		this.list.Refresh()
		return
	}
	for i1 := 0; i1 < data.Len(); i1++ {
		n1 := data.At(i1).Filename
		for i2 := 0; i2 < i1; i2++ {
			val := data.Diff(i2, i1)
			if val < filterVal {
				continue
			}
			n2 := data.At(i2).Filename
			txt := fmt.Sprintf("Equal: %v\n%s\n%s", val, n1, n2)
			this.data = append(this.data, entry{
				id1:  i2,
				id2:  i1,
				text: txt,
			})
		}
	}
	this.list.Refresh()
}
