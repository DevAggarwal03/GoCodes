package main

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main () {
	fmt.Println("Hello world")

	a := app.New()
	newWin := a.NewWindow("Hello World") 
	
	caption := widget.NewLabel("Random Number Generator")


	button := widget.NewButton("Button", func () {
		time := time.Now().Format("TIME: 03:04:05")
		caption.SetText(time)	
	})

	newWin.SetContent(container.NewVBox(caption, button))
	newWin.ShowAndRun()
}
