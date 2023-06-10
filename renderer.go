package main

import (
	"fmt"
	"image"
	"math"
)

var frameIndex int = 1
var finalImage *image.RGBA

var accumulationData []Vec3

func Render() {
	if frameIndex == 1 {
		resetAccumulatedData()
		fmt.Println("Reset accumulated data")
	}

	rectSize := finalImage.Rect.Size()

	for y := 0; y < rectSize.Y; y++ {
		for x := 0; x < rectSize.X; x++ {
			color := perPixel(x, y)

			accumulatedPixel := Add(accumulationData[x+y*rectSize.X], color)
			accumulationData[x+y*rectSize.X] = accumulatedPixel

			accumulatedPixel = Divide(accumulatedPixel, float64(frameIndex))
			accumulatedPixel = Clamp01(accumulatedPixel)

			finalImage.Set(x, y, accumulatedPixel.ToRGB())
		}
	}

	frameIndex++
}

func Resize(w, h int) {
	if finalImage != nil &&
		(finalImage.Rect.Size().X == w && finalImage.Rect.Size().Y == h) {
		return
	}

	finalImage = image.NewRGBA(image.Rect(0, 0, w, h))
	accumulationData = make([]Vec3, w*h)
}

func GetFinalImage() *image.RGBA {
	return finalImage
}

func perPixel(x, y int) Vec3 {
	ray := Ray{}
	ray.Origin = cam.position
	ray.Dir = cam.rayDirections[x+y*cam.viewportWidth]

	light := Vec3{}
	contribution := Vec3One()

	for i := 0; i < 5; i++ {
		payload := traceRay(ray)
		if payload.Distance < 0 {
			break
		}

		sphere := scene.Spheres[payload.ObjectIndex]
		material := scene.Materials[sphere.MaterialId]

		light = Add(light, material.GetEmission())
		contribution = MWiseMul(contribution, material.Albedo)
		ray.Origin = Add(payload.Point, Mul(payload.Normal, 0.0001))
		ray.Dir = Normalized(Add(payload.Normal, RandomInUnitSphere(rng)))
	}

	return light
}

func traceRay(r Ray) HitRecord {
	closestSphere := -1
	hitDistance := math.MaxFloat64

	for i := range scene.Spheres {
		sphere := scene.Spheres[i]
		origin := Subtract(r.Origin, sphere.Center)
		a := Dot(r.Dir, r.Dir)
		b := Dot(origin, r.Dir) * 2.0
		c := Dot(origin, origin) - sphere.Radius*sphere.Radius
		discriminant := b*b - 4.0*a*c
		if discriminant < 0 {
			continue
		}

		closest := (-b - math.Sqrt(discriminant)) / (2.0 * a)
		// closest1 := (-b + math.Sqrt(discriminant)) / (2.0 * a)

		if closest > 0.0 && closest < hitDistance {
			hitDistance = closest
			closestSphere = i
		}
	}

	if closestSphere < 0 {
		return HitRecord{Distance: -1}
	}

	return closestHit(r, hitDistance, closestSphere)
}

func closestHit(r Ray, distance float64, objectIndex int) HitRecord {
	payload := HitRecord{}
	payload.Distance = distance
	payload.ObjectIndex = objectIndex

	closestSphere := scene.Spheres[objectIndex]
	origin := Subtract(r.Origin, closestSphere.Center)
	payload.Point = Add(origin, Mul(r.Dir, distance))
	payload.Normal = Normalized(payload.Point)
	payload.Point = Add(payload.Point, closestSphere.Center)

	return payload
}

func resetAccumulatedData() {
	if len(accumulationData) == 0 {
		return
	}

	accumulationData[0] = Vec3{}
	for bp := 1; bp < len(accumulationData); bp *= 2 {
		copy(accumulationData[bp:], accumulationData[:bp])
	}
}
