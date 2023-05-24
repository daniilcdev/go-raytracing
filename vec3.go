package main

import "math"

type Vec3 struct {
	X, Y, Z float64
}

func (v3 *Vec3) Add(other Vec3) *Vec3 {
	v3.X += other.X
	v3.Y += other.Y
	v3.Z += other.Z
	return v3
}

func (v3 *Vec3) Scale(s float64) *Vec3 {
	v3.X *= s
	v3.Y *= s
	v3.Z *= s

	return v3
}

func (v3 *Vec3) Mag() float64 {
	return math.Sqrt(v3.SqrMag())
}

func (v3 *Vec3) SqrMag() float64 {
	return v3.X*v3.X + v3.Y*v3.Y + v3.Z*v3.Z
}

func Vec3One() Vec3 {
	return Vec3{1, 1, 1}
}

func Cross(a Vec3, b Vec3) Vec3 {
	return Vec3{
		X: a.Y*b.Z - a.Z*b.Y,
		Y: a.Z*b.X - a.X*b.Z,
		Z: a.X*b.Y - a.Y*b.X,
	}
}

func Dot(a Vec3, b Vec3) float64 {
	return a.X*b.X + a.Y*b.Y + a.Z*b.Z
}

func Normalized(v Vec3) Vec3 {
	return *v.Scale(1 / v.Mag())
}

func Mul(v Vec3, t float64) Vec3 {
	return Vec3{v.X * t, v.Y * t, v.Z * t}
}

func Subtract(a Vec3, b Vec3) Vec3 {
	return Vec3{a.X - b.X, a.Y - b.Y, a.Z - b.Z}
}

func Add(a Vec3, b Vec3) Vec3 {
	return Vec3{a.X + b.X, a.Y + b.Y, a.Z + b.Z}
}
