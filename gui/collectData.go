package gui

import (
	"errors"
	"io/fs"
	"path/filepath"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"

	"goImgCmp/proc"

	"fmt"
)

type CollectData struct {
	Content   fyne.CanvasObject
	OnAbort   func()
	OnFinish  func(*proc.CompareResult)
	shouldRun bool
	current   binding.String
	scan      binding.String
	process   binding.String
	cmpRes    *proc.CompareResult
}

func NewCollectData() *CollectData {
	r := new(CollectData)

	r.current = binding.NewString()
	r.scan = binding.NewString()
	r.process = binding.NewString()

	formLayout := container.New(
		layout.NewFormLayout(),

		widget.NewLabel("Current: "),
		widget.NewLabelWithData(r.current),

		widget.NewLabel("Scan: "),
		widget.NewLabelWithData(r.scan),

		widget.NewLabel("Process: "),
		widget.NewLabelWithData(r.process),

		widget.NewLabel(""),
		widget.NewButton("Abort", func() { r.onAbort() }),
	)
	r.Content = formLayout
	return r
}

func (this *CollectData) Start(folder []string) {
	this.shouldRun = true
	this.current.Set("None")
	this.scan.Set("0")
	this.process.Set("0")
	cmpRes := proc.NewCompareResult()
	this.cmpRes = cmpRes

	fileList := make([]string, 0)
	for i, v := range folder {
		if !this.shouldRun {
			return
		}

		prefix := "(" + strconv.Itoa(i) + ") "
		newFiles, err := this.scanFolder(v, prefix, len(fileList))
		if err != nil {
			fmt.Println("Failed to scan ", v, " because: ", err)
		} else {
			fileList = append(fileList, newFiles...)
		}
	}
	this.setScanProgress(len(fileList), true)

	this.loadJobs(fileList)

	if this.shouldRun {
		this.onFinish()
	}
}

func (this *CollectData) scanFolder(
	folder string, prefix string, offset int) ([]string, error) {

	fileList := make([]string, 0)
	err := filepath.Walk(folder, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !this.shouldRun {
			return errors.New("User Abort")
		}
		this.setCurrentLabel(prefix + path)
		if !proc.JobCanProcess(path) {
			return nil
		}

		fileList = append(fileList, path)
		this.setScanProgress(offset+len(fileList), false)

		return nil
	})
	if err != nil {
		return nil, err
	} else {
		return fileList, nil
	}
}

func (this *CollectData) loadJobs(fileList []string) {
	for i, path := range fileList {
		if !this.shouldRun {
			return
		}
		this.setCurrentLabel(path)
		job := proc.NewJob(path)
		if !job.IsValid() {
			fmt.Println("Invalid job: ", path)
			continue
		}
		this.cmpRes.Add(job)
		this.setProcessingProgress(i+1, false)
	}
	this.setProcessingProgress(len(fileList), true)

}

func (this *CollectData) setCurrentLabel(name string) {
	this.current.Set(name)
}

func (this *CollectData) setScanProgress(p int, f bool) {
	txt := strconv.Itoa(p)
	if f {
		txt = txt + " (done)"
	}
	this.scan.Set(txt)
}

func (this *CollectData) setProcessingProgress(p int, f bool) {
	txt := strconv.Itoa(p)
	if f {
		txt = txt + " (done)"
	}
	this.process.Set(txt)
}

func (this *CollectData) onAbort() {
	this.shouldRun = false
	if this.OnAbort != nil {
		this.OnAbort()
	} else {
		fmt.Println("Action: OnAbort")
	}
}

func (this *CollectData) onFinish() {
	this.shouldRun = false
	if this.OnFinish != nil {
		this.OnFinish(this.cmpRes)
	} else {
		fmt.Println("Action: OnFinish")
	}
}
