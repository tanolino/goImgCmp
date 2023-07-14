package main

import (
	"fmt"
	"goImgCmp/gui"
)

func main() {
	fmt.Println("Starting")
	mainGui := gui.NewMainGui()
	mainGui.ShowAndRun()
}
