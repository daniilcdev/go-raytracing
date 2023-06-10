package main

type HitRecord struct {
	ObjectIndex int
	Point       Vec3
	Normal      Vec3
	Distance    float64
	FrontFace   bool
}
