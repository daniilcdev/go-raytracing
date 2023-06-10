package main

import (
	"fmt"
	"math/rand"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

const aspect float32 = 16.0 / 9.0
const imageW int = 400
const imageH int = int(float32(imageW) / aspect)

var rng *rand.Rand = rand.New(rand.NewSource(105))

var world HittableList
var cam Camera

var materials [3]Material

func main() {
	cam = NewCameraAt(Vec3{0.0, 0.1, 10})

	materials = [3]Material{}
	materials[0] = Material{Albedo: Vec3{1, 0, 1}}
	materials[1] = Material{Albedo: Vec3{.2, 0.3, 1}, Roughness: 0.1}
	materials[2] = Material{Albedo: Vec3{.8, 0.5, 0.2}, Roughness: 0.1, EmissionPower: 2., EmissionColor: Vec3{.8, 0.5, 0.2}}

	var objs []Hittable = make([]Hittable, 0)
	objs = append(objs,
		Sphere{
			Center:     Vec3{0, 0.0, 0},
			Radius:     1,
			MaterialId: 0,
		},
		Sphere{
			Center:     Vec3{2, 0, 0},
			Radius:     1,
			MaterialId: 2,
		},
		Sphere{
			Center:     Vec3{.0, -101, 0},
			Radius:     100,
			MaterialId: 1,
		},
	)

	world = HittableList{objects: objs}

	a := app.New()
	w := a.NewWindow("Viewport")

	clock := widget.NewLabel("FPS: 60.00")

	Resize(imageW, imageH)
	renderTexture := GetFinalImage()

	img := canvas.NewImageFromImage(renderTexture)
	img.FillMode = canvas.ImageFillOriginal | canvas.ImageFillStretch

	size := fyne.NewSize(float32(imageW)+clock.Size().Width, float32(imageH))
	hBox := container.NewHBox(container.NewVBox(img), clock)

	w.SetContent(hBox)
	w.Resize(size)

	go renderScene(img, clock)

	w.ShowAndRun()
	tidyUp()
}

func renderScene(img *canvas.Image,
	clock *widget.Label) {
	fmt.Println("Cold boot...")
	coldBoot := time.After(time.Second / 4)
	<-coldBoot
	fmt.Println("Done")
	img.Resize(fyne.NewSize(float32(imageW), float32(imageH)))
	img.Refresh()

	t := time.Second / 60
	loopDuration := time.Duration(t)

	for {
		start := time.Now()
		Render()
		duration := time.Since(start)

		clock.SetText(fmt.Sprintf("FPS: %f", 1.0/duration.Seconds())[:10])
		img.Refresh()

		go time.After(loopDuration)
	}
}

func tidyUp() {
	fmt.Println("Exited")
}
