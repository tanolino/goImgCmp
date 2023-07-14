package proc

import (
	"bytes"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"math"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/image/webp"
)

const divider int = 16
const dividerF float64 = 16

type Job struct {
	Filename string
	width    int
	height   int
	vecs     [divider][divider][4]float64
}

func NewJob(f string) *Job {
	r := new(Job)
	r.Filename = f
	r.process()
	return r
}

func (this *Job) IsValid() bool {
	return this.width > 0 && this.height > 0
}

var readWithImage = []string{".png", ".gif", ".jpeg", ".jpg"}
var readWithWebp = ".webp"

func JobCanProcess(file string) bool {
	ext := strings.ToLower(filepath.Ext(file))
	for _, v := range readWithImage {
		if ext == v {
			return true
		}
	}
	if ext == readWithWebp {
		return true
	}
	return false
}

func (this *Job) process() {
	img := this.GetImage()
	if img != nil {
		this.calculateVectors(img)
	}
}

func (this *Job) GetImage() image.Image {
	buffer := readToBuffer(this.Filename)
	if buffer == nil {
		return nil
	}

	img := decode(this.Filename, buffer)
	if img == nil {
		return nil
	} else {
		return img
	}
}

func readToBuffer(file string) []byte {
	var buffer, err = os.ReadFile(file)
	if err != nil {
		fmt.Println("Failed to read ", file, " because ", err)
		return nil
	} else {
		return buffer
	}
}

func decode(file string, buffer []byte) image.Image {
	ext := strings.ToLower(filepath.Ext(file))
	reader := bytes.NewReader(buffer)

	for _, v := range readWithImage {
		if ext == v {
			img, _, err := image.Decode(reader)
			if err != nil {
				fmt.Println("Failed to decode ", file, " because ", err)
			} else {
				return img
			}
		}
	}
	if ext == readWithWebp {
		img, err := webp.Decode(reader)
		if err != nil {
			fmt.Println("Failed to decode ", file, " because ", err)
		} else {
			return img
		}
	}
	fmt.Println("Failed to find decoder for file ", file)
	return nil
}

func (this *Job) calculateVectors(img image.Image) {
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	if height < divider || width < divider {
		fmt.Println("Image ", this.Filename, " too small")
		return
	}

	for row := 0; row < divider; row++ {
		for col := 0; col < divider; col++ {
			res := calcVector(img, row, col)
			this.vecs[row][col] = res
		}
	}

	this.width = width
	this.height = height
}

func calcVector(img image.Image, row int, col int) [4]float64 {
	var r, g, b, a int64

	bounds := img.Bounds()
	minX := calcEffectiveStart(bounds.Min.X, bounds.Max.X, col)
	maxX := calcEffectiveEnd(bounds.Min.X, bounds.Max.X, col)
	minY := calcEffectiveStart(bounds.Min.Y, bounds.Max.Y, row)
	maxY := calcEffectiveEnd(bounds.Min.Y, bounds.Max.Y, row)

	for x := minX; x < maxX; x++ {
		for y := minY; y < maxY; y++ {
			cR, cG, cB, cA := img.At(x, y).RGBA()
			r += (int64)(cR >> 8)
			g += (int64)(cG >> 8)
			b += (int64)(cB >> 8)
			a += (int64)(cA >> 8)
		}
	}

	pixelF := (float64)((maxX - minX) * (maxY - minY) * 256)
	return [4]float64{
		(float64)(r) / pixelF,
		(float64)(g) / pixelF,
		(float64)(b) / pixelF,
		(float64)(a) / pixelF,
	}
}

func calcEffectiveStart(min, max, pos int) int {
	if pos == 0 {
		return min
	} else {
		diff := max - min
		return min + pos*(diff/divider)
	}
}

func calcEffectiveEnd(min, max, pos int) int {
	if pos+1 == divider {
		return max
	} else {
		diff := max - min
		return min + (pos+1)*(diff/divider)
	}
}

func (this *Job) CompareTo(that *Job) float64 {
	var eqal float64 = 0
	if !this.IsValid() {
		fmt.Println("Invalid Job(1)")
		return eqal
	}
	if !that.IsValid() {
		fmt.Println("Invalid Job(2)")
		return eqal
	}
	// row
	for r := 0; r < divider; r++ {
		// column
		for c := 0; c < divider; c++ {
			// colour
			var colorDiff float64 = 0
			for i := 0; i < 3; i++ {
				colorDiff += math.Abs(this.vecs[r][c][i] - that.vecs[r][c][i])
				// fmt.Println("ColorDiff ", colorDiff)
			}
			colorDiff /= 3
			colorDiff = 1 - colorDiff
			eqal += colorDiff
		}
	}
	return eqal / (dividerF * dividerF)
}
