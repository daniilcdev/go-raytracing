package main

import "math"

type HitRecord struct {
	Point     Vec3
	Normal    Vec3
	Distance  float64
	FrontFace bool
}

func (rec *HitRecord) setFaceNormal(ray *Ray, outwardNormal *Vec3) {
	rec.FrontFace = Dot(ray.Dir, *outwardNormal) < 0
	if rec.FrontFace {
		rec.Normal = *outwardNormal
	} else {
		rec.Normal = Mul(*outwardNormal, -1)
	}
}

type Hittable interface {
	Hit(ray *Ray, tMin float64, tMax float64, rec *HitRecord) bool
}

type Sphere struct {
	Center Vec3
	Radius float64
}

func (s Sphere) Hit(ray *Ray, tMin float64, tMax float64, rec *HitRecord) bool {
	oc := Subtract(ray.Origin, s.Center)
	a := ray.Dir.SqrMag()
	half_b := Dot(oc, ray.Dir)
	c := oc.SqrMag() - s.Radius*s.Radius

	discriminant := half_b*half_b - a*c
	if discriminant < 0 {
		return false
	}
	sqrtd := math.Sqrt(discriminant)
	root := (-half_b - sqrtd) / a

	if root < tMin || tMax < root {
		root = (-half_b + sqrtd) / a
		if root < tMin || tMax < root {
			return false
		}
	}

	rec.Distance = root
	rec.Point = ray.At(root)
	outwardNormal := Divide(Subtract(rec.Point, s.Center), s.Radius)
	rec.setFaceNormal(ray, &outwardNormal)

	return true
}

type HittableList struct {
	objects []Hittable
}

func (w HittableList) Hit(ray *Ray, tMin float64, tMax float64, rec *HitRecord) bool {
	var tempRec HitRecord = *rec
	var didHitAnything bool = false
	closest := tMax

	for i := 0; i < len(w.objects); i++ {
		o := w.objects[i]
		if o.Hit(ray, tMin, closest, &tempRec) {
			didHitAnything = true
			closest = tempRec.Distance

			rec.Point = tempRec.Point
			rec.Distance = tempRec.Distance
			rec.Normal = tempRec.Normal
		}
	}

	return didHitAnything
}
