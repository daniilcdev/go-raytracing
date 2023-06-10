package main

type HitRecord struct {
	ObjectIndex int
	Point       Vec3
	Normal      Vec3
	Distance    float64
	FrontFace   bool
}

type Hittable interface {
	WithRadius
	MaterialIndex() int
	CenterPoint() Vec3
}

type WithRadius interface {
	R() float64
}

type Sphere struct {
	Center     Vec3
	Radius     float64
	MaterialId int
}

func (s Sphere) MaterialIndex() int {
	return s.MaterialId
}

func (s Sphere) CenterPoint() Vec3 {
	return s.Center
}

func (s Sphere) R() float64 {
	return s.Radius
}

type HittableList struct {
	objects []Hittable
}
