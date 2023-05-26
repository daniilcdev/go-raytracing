package main

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"math/rand"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

const aspect float32 = 16.0 / 9.0
const imageW int = 720
const imageH int = int(float32(imageW) / aspect)
const samplesPerPixel int = 1

func renderTexture() *image.RGBA {
	return image.NewRGBA(image.Rect(0, 0, imageW, imageH))
}

var world Hittable
var cam Camera

func main() {
	cam = NewCamera()

	var objs []Hittable = make([]Hittable, 0)
	objs = append(objs, Sphere{
		Center: Vec3{0, 0, -1},
		Radius: 0.5,
	})
	objs = append(objs, Sphere{
		Center: Vec3{0, -100.5, -1},
		Radius: 100,
	})

	world = HittableList{objects: objs}

	a := app.New()
	w := a.NewWindow("Viewport")
	clock := widget.NewLabel("")
	renderTexture := renderTexture()
	img := canvas.NewImageFromImage(renderTexture)
	img.FillMode = canvas.ImageFillOriginal | canvas.ImageFillStretch

	updateTime(clock)

	size := fyne.NewSize(float32(imageW)+clock.Size().Width, float32(imageH))
	hBox := container.NewHBox(container.NewVBox(img), clock)

	w.SetContent(hBox)
	w.Resize(size)

	go func() {
		for range time.Tick(time.Second) {
			updateTime(clock)
		}
	}()

	go renderScene(renderTexture, img)

	w.ShowAndRun()
	tidyUp()
}

func renderScene(renderTexture *image.RGBA, img *canvas.Image) {
	fmt.Println("Cold boot...")
	coldBoot := time.After(time.Second / 4)
	<-coldBoot
	fmt.Println("Done")
	img.Resize(fyne.NewSize(float32(imageW), float32(imageH)))
	img.Refresh()

	t := time.Second / 60
	duration := time.Duration(t)

	for range time.Tick(duration) {
		for y := 0; y < imageH; y++ {
			for x := 0; x < imageW; x++ {
				go perPixel(x, y, renderTexture)
			}
		}

		img.Refresh()
	}
}

func perPixel(x, y int, renderTexture *image.RGBA) {
	var samples Vec3 = Vec3{}
	for s := 0; s < samplesPerPixel; s++ {
		u := (float64(x) + rand.Float64()) / float64(imageW)
		v := 1 - (float64(y)+rand.Float64())/float64(imageH)

		ray := cam.getRay(u, v)
		samples.Add(rayColor(&ray, &world))
	}

	samples.Scale(1. / float64(samplesPerPixel)).Scale(0.1)

	var c color.RGBA = renderTexture.RGBAAt(x, y)
	pixel := Vec3{
		X: float64(c.R) / 255,
		Y: float64(c.G) / 255,
		Z: float64(c.B) / 255,
	}
	pixel.Scale(0.9).Add(samples)
	renderTexture.Set(x, y, pixel.ToRGB())
}

func rayColor(ray *Ray, world *Hittable) Vec3 {
	rec := HitRecord{}

	if (*world).Hit(ray, 0, math.Inf(1), &rec) {
		c := Mul(Add(rec.Normal, Vec3One()), 0.5)
		return c
	}

	unitDir := Normalized(ray.Dir)
	hitDistance := 0.5 * (unitDir.Y + 1.0)
	c := Add(Mul(Vec3One(), 1-hitDistance), Mul(Vec3{0.5, 0.7, 1}, hitDistance))

	return c
}

func updateTime(clock *widget.Label) {
	formatted := time.Now().Format("Time: 15:04:05")
	clock.SetText(formatted)
}

func tidyUp() {
	fmt.Println("Exited")
}
