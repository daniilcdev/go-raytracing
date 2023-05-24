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
const imageW int = 480
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

func getHitDistance(center *Vec3, radius float64, r *Ray) float64 {
	oc := Subtract(r.Origin, *center)
	a := r.Dir.SqrMag()
	half_b := Dot(oc, r.Dir)
	c := oc.SqrMag() - radius*radius

	discriminant := half_b*half_b - a*c
	if discriminant > 0 {
		return (-half_b - math.Sqrt(discriminant)) / a
	} else {
		return -1
	}
}

func rayColor(ray *Ray) color.Color {
	hitDistance := getHitDistance(&Vec3{Z: -1}, 0.5, ray)
	if hitDistance > 0 {
		n := Normalized(Subtract(ray.At(hitDistance), Vec3{Z: -1}))
		return n.Add(Vec3One()).Scale(0.5).ToRGB()
	}

	unitDir := Normalized(ray.Dir)
	hitDistance = 0.5 * (unitDir.Y + 1.0)
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
