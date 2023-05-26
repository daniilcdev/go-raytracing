package main

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

const aspect float32 = 16.0 / 9.0
const imageW int = 480
const imageH int = int(float32(imageW) / aspect)

func renderTexture() *image.RGBA {
	return image.NewRGBA(image.Rect(0, 0, imageW, imageH))
}

var world Hittable

func main() {
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
	coldBoot := time.After(time.Second)
	<-coldBoot
	fmt.Println("Done")
	img.Resize(fyne.NewSize(float32(imageW), float32(imageH)))
	img.Refresh()

	const viewportH = 2.0
	const viewportW = aspect * viewportH
	const focalLength = 1.0

	origin := Vec3{}
	horizontal := Vec3{X: float64(viewportW)}
	vertical := Vec3{Y: float64(viewportH)}

	llc := Subtract(origin, Mul(horizontal, 0.5))
	llc = Subtract(llc, Mul(vertical, 0.5))
	llc = Subtract(llc, Vec3{Z: focalLength})

	t := time.Second / 60
	duration := time.Duration(t)

	for range time.Tick(duration) {
		for y := 0; y < imageH; y++ {
			for x := 0; x < imageW; x++ {
				go perPixel(x, y, renderTexture, origin, llc, horizontal, vertical)
			}
		}

		img.Refresh()
	}
}

func perPixel(x, y int, renderTexture *image.RGBA, origin Vec3, llc Vec3, horizontal Vec3, vertical Vec3) {
	u := float64(x) / float64(imageW-1)
	v := 1 - float64(y)/float64(imageH-1)

	dir := Subtract(Add(Add(llc, Mul(horizontal, u)), Mul(vertical, v)), origin)
	ray := Ray{Origin: origin, Dir: dir}
	pixel := rayColor(&ray, &world)
	renderTexture.Set(x, y, pixel)
}

func rayColor(ray *Ray, world *Hittable) color.Color {
	rec := HitRecord{}

	if (*world).Hit(ray, 0, math.Inf(1), &rec) {
		c := Mul(Add(rec.Normal, Vec3One()), 0.5)
		return c.ToRGB()
	}

	unitDir := Normalized(ray.Dir)
	hitDistance := 0.5 * (unitDir.Y + 1.0)
	c := Add(Mul(Vec3One(), 1-hitDistance), Mul(Vec3{0.5, 0.7, 1}, hitDistance))

	return c.ToRGB()
}

func updateTime(clock *widget.Label) {
	formatted := time.Now().Format("Time: 15:04:05")
	clock.SetText(formatted)
}

func tidyUp() {
	fmt.Println("Exited")
}
