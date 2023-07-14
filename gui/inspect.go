package gui

import (
	"fmt"
	"goImgCmp/proc"
	"image"
	"image/color"
	"math"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type Inspect struct {
	Content  fyne.CanvasObject
	OnReturn func(*proc.CompareResult)
	id1      int
	id2      int
	dataRaw  *proc.CompareResult
}

func NewInspect() *Inspect {
	r := new(Inspect)

	r.Content = widget.NewLabel("WORKING")
	return r
}

func (this *Inspect) SetData(id1 int, id2 int, dataRaw *proc.CompareResult) {
	this.id1 = id1
	this.id2 = id2
	this.dataRaw = dataRaw

	img1 := newImage(dataRaw.At(id1))
	img2 := newImage(dataRaw.At(id2))
	imgDiff := newDiffImage(img1.Image, img2.Image)

	this.Content = container.New(
		layout.NewGridLayout(2),
		img1, img2,
		imgDiff, widget.NewButton("Back", this.onReturn),
	)
}

func newImage(job *proc.Job) *canvas.Image {
	img := canvas.NewImageFromImage(job.GetImage())
	img.FillMode = canvas.ImageFillContain
	img.ScaleMode = canvas.ImageScaleFastest
	return img
}

func (this *Inspect) onReturn() {
	if this.OnReturn != nil {
		this.OnReturn(this.dataRaw)
	} else {
		fmt.Println("OnReturn from Inspect")
	}
}

func newDiffImage(img1, img2 image.Image) *canvas.Image {
	if img1 == nil || img2 == nil {
		return nil
	}
	b1 := img1.Bounds()
	b2 := img2.Bounds()
	minX := minInt(b1.Dx(), b2.Dx())
	minY := minInt(b1.Dy(), b2.Dy())

	diffColor := func(c1, c2 color.Color) color.RGBA {
		r1, g1, b1, _ := c1.RGBA()
		r2, g2, b2, _ := c2.RGBA()

		r3 := diffColor(r1, r2)
		g3 := diffColor(g1, g2)
		b3 := diffColor(b1, b2)

		return color.RGBA{
			R: r3,
			G: g3,
			B: b3,
			A: 255,
		}
	}

	r := image.NewRGBA(image.Rect(0, 0, minX, minY))
	for y := 0; y < minY; y++ {
		for x := 0; x < minX; x++ {
			diff := diffColor(img1.At(x, y), img2.At(x, y))
			r.SetRGBA(x, y, diff)
		}
	}

	return canvas.NewImageFromImage(r)
}

func minInt(i1, i2 int) int {
	if i1 < i2 {
		return i1
	} else {
		return i2
	}
}

func maxInt(i1, i2 int) int {
	if i1 > i2 {
		return i1
	} else {
		return i2
	}
}

func diffColor(i1, i2 uint32) uint8 {
	f1 := (float64)(i1) / (float64)(65536)
	f2 := (float64)(i2) / (float64)(65536)

	fd := math.Abs(f1 - f2)

	// scale up these lows
	fd = 1 - fd
	fd = fd * fd * fd * fd
	fd = 1 - fd

	if fd < 0 {
		fd = 0
	} else if fd > 1 {
		fd = 1
	}

	return (uint8)(fd * 255)
}
