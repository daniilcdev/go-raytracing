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
const imageW int = 400
const imageH int = int(float32(imageW) / aspect)
const maxDepth int = 10

func renderTexture() *image.RGBA {
	return image.NewRGBA(image.Rect(0, 0, imageW, imageH))
}

var world Hittable
var cam Camera

func main() {
	cam = NewCamera()

	var objs []Hittable = make([]Hittable, 0)
	objs = append(objs,
		Sphere{
			Center: Vec3{0, -100.5, -1},
			Radius: 100,
		},
		Sphere{
			Center: Vec3{0, 0, -1},
			Radius: 0.5,
		},
		Sphere{
			Center: Vec3{-1, 0, -1},
			Radius: 0.5,
		},
		Sphere{
			Center: Vec3{1, 0, -1},
			Radius: 0.5,
		},
	)

	world = HittableList{objects: objs}

	a := app.New()
	w := a.NewWindow("Viewport")
	clock := widget.NewLabel("FPS: 60.00")
	renderTexture := renderTexture()
	img := canvas.NewImageFromImage(renderTexture)
	img.FillMode = canvas.ImageFillOriginal | canvas.ImageFillStretch

	size := fyne.NewSize(float32(imageW)+clock.Size().Width, float32(imageH))
	hBox := container.NewHBox(container.NewVBox(img), clock)

	w.SetContent(hBox)
	w.Resize(size)

	go renderScene(renderTexture, img, clock)

	w.ShowAndRun()
	tidyUp()
}

func renderScene(renderTexture *image.RGBA, img *canvas.Image, clock *widget.Label) {
	fmt.Println("Cold boot...")
	coldBoot := time.After(time.Second / 4)
	<-coldBoot
	fmt.Println("Done")
	img.Resize(fyne.NewSize(float32(imageW), float32(imageH)))
	img.Refresh()

	t := time.Second / 60
	loopDuration := time.Duration(t)

	for {
		pixels := make(chan Pixel, imageW*imageH)

		start := time.Now()
		go writePixels(renderTexture, pixels)

		for pixel := range pixels {
			renderTexture.Set(pixel.X, pixel.Y, pixel.Color.ToRGB())
		}

		duration := time.Since(start)

		clock.SetText(fmt.Sprintf("FPS: %f", 1.0/duration.Seconds())[:10])

		img.Refresh()

		go time.After(loopDuration)
	}
}

func writePixels(renderTexture *image.RGBA, buffer chan Pixel) {
	for y := 0; y < imageH; y++ {
		for x := 0; x < imageW; x++ {
			sample := Vec3{}
			u := (float64(x) + rand.Float64()*2 - 1) / float64(imageW)
			v := 1 - (float64(y)+rand.Float64()*2-1)/float64(imageH)

			ray := cam.getRay(u, v)
			sample.Add(rayColor(&ray, &world, maxDepth))

			const lerp float64 = 0.5
			sample.Sqrt().Scale(1 - lerp)

			var c color.RGBA = renderTexture.RGBAAt(x, y)
			pixelColor := Vec3{
				X: float64(c.R) / 255,
				Y: float64(c.G) / 255,
				Z: float64(c.B) / 255,
			}

			pixelColor = Add(Mul(pixelColor, lerp), sample)

			buffer <- Pixel{X: x, Y: y, Color: pixelColor}
		}
	}

	close(buffer)
}

func rayColor(ray *Ray, world *Hittable, depth int) Vec3 {
	rec := HitRecord{}
	if depth <= 0 {
		return Vec3{}
	}

	if (*world).Hit(ray, 0, math.Inf(1), &rec) {
		target := Add(Add(rec.Point, rec.Normal), RandomInUnitSphere())
		nextRay := Ray{rec.Point, Subtract(target, rec.Point)}
		nextColor := rayColor(&nextRay, world, depth-1)
		return Mul(nextColor, 0.5)
	}

	unitDir := Normalized(ray.Dir)
	hitDistance := 0.5 * (unitDir.Y + 1.0)
	c := Add(Mul(Vec3One(), 1-hitDistance), Mul(Vec3{0.5, 0.7, 1}, hitDistance))

	return c
}

func tidyUp() {
	fmt.Println("Exited")
}
