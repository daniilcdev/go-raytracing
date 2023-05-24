package main

import (
	"fmt"
	"image"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	a := app.New()
	w := a.NewWindow("Viewport")
	clock := widget.NewLabel("")
	img := canvas.NewImageFromImage(render())
	img.FillMode = canvas.ImageFillOriginal
	updateTime(clock)

	w.SetContent(container.NewVBox(
		clock, img,
	))
	w.Resize(fyne.NewSize(512+clock.Size().Height, 512))

	go func() {
		for range time.Tick(time.Second) {
			updateTime(clock)
		}
	}()

	w.ShowAndRun()
	tidyUp()
}

func render() image.Image {
	img := image.NewRGBA(image.Rect(0, 0, 512, 512))
	for x := 0; x < 512; x++ {
		for y := 0; y < 512; y++ {
			v := Vec3{float64(x) / (512 - 1), 1 - float64(y)/(512-1), 0.25}
			img.Set(x, y, v.ToRGB())
		}
	}

	return img
}

func updateTime(clock *widget.Label) {
	formatted := time.Now().Format("Time: 03:04:05")
	clock.SetText(formatted)
}

func tidyUp() {
	fmt.Println("Exited")
}
