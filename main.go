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

const aspect float64 = 16.0 / 9.0
const imageW int = 640
const imageH int = int(float64(imageW) / aspect)

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
	w.Resize(fyne.NewSize(float32(imageW)+clock.Size().Height, float32(imageH)))

	go func() {
		for range time.Tick(time.Second) {
			updateTime(clock)
		}
	}()

	w.ShowAndRun()
	tidyUp()
}

func render() image.Image {
	const viewportH = 2.0
	const viewportW = aspect * viewportH
	const focalLength = 1.0

	origin := Vec3{}
	horizontal := Vec3{X: viewportW}
	vertical := Vec3{Y: viewportH}

	llc := Subtract(origin, Mul(horizontal, 0.5))
	llc = Subtract(llc, Mul(vertical, 0.5))
	llc = Subtract(llc, Vec3{Z: focalLength})

	img := image.NewRGBA(image.Rect(0, 0, imageW, imageH))
	for x := 0; x < imageW; x++ {
		for y := 0; y < imageH; y++ {
			u := float64(x) / float64(imageW-1)
			v := 1 - float64(y)/float64(imageH-1)

			dir := Subtract(Add(Add(llc, Mul(horizontal, u)), Mul(vertical, v)), origin)
			ray := Ray{Origin: origin, Dir: dir}
			pixel := rayColor(&ray)
			img.Set(x, y, pixel)
		}
	}

	return img
}

func getCastDistance(center *Vec3, radius float64, r *Ray) float64 {
	oc := Subtract(r.Origin, *center)
	a := Dot(r.Dir, r.Dir)
	b := 2.0 * Dot(oc, r.Dir)
	c := Dot(oc, oc) - radius*radius

	discriminant := b*b - 4*a*c
	if discriminant > 0 {
		return (-b - math.Sqrt(discriminant)) / (2 * a)
	} else {
		return -1
	}
}

func rayColor(ray *Ray) color.Color {
	hitDistance := getCastDistance(&Vec3{Z: -1}, 0.5, ray)
	if hitDistance > 0 {
		n := Normalized(Subtract(ray.At(hitDistance), Vec3{Z: -1}))
		return n.Add(Vec3One()).Scale(0.5).ToRGB()
	}

	unitDir := Normalized(ray.Dir)
	t := 0.5 * (unitDir.Y + 1.0)
	c := Add(Mul(Vec3One(), 1-t), Mul(Vec3{0.5, 0.7, 1}, t))

	return c.ToRGB()
}

func updateTime(clock *widget.Label) {
	formatted := time.Now().Format("Time: 03:04:05")
	clock.SetText(formatted)
}

func tidyUp() {
	fmt.Println("Exited")
}
