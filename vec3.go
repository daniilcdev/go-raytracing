package main

import (
	"math"
	"math/rand"
)

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

func (v3 *Vec3) Sqrt() *Vec3 {
	v3.X = math.Sqrt(v3.X)
	v3.Y = math.Sqrt(v3.Y)
	v3.Z = math.Sqrt(v3.Z)

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
	return Divide(v, v.Mag())
}

func Mul(v Vec3, t float64) Vec3 {
	return Vec3{v.X * t, v.Y * t, v.Z * t}
}

func MWiseMul(v, s Vec3) Vec3 {
	return Vec3{v.X * s.X, v.Y * s.Y, v.Z * s.Z}
}

func Divide(v Vec3, t float64) Vec3 {
	return Vec3{v.X / t, v.Y / t, v.Z / t}
}

func Clamp01(v Vec3) Vec3 {
	return Vec3{
		math.Max(math.Min(1, v.X), 0),
		math.Max(math.Min(1, v.Y), 0),
		math.Max(math.Min(1, v.Z), 0),
	}
}

func Subtract(a Vec3, b Vec3) Vec3 {
	return Vec3{a.X - b.X, a.Y - b.Y, a.Z - b.Z}
}

func Add(a Vec3, b Vec3) Vec3 {
	return Vec3{a.X + b.X, a.Y + b.Y, a.Z + b.Z}
}

func Random(rng *rand.Rand) Vec3 {
	return Vec3{rng.Float64(), rng.Float64(), rng.Float64()}
}

func RandomRange(min, max float64, rng *rand.Rand) Vec3 {
	return Vec3{
		rng.Float64()*(max-min) + min,
		rng.Float64()*(max-min) + min,
		rng.Float64()*(max-min) + min}
}

func RandomInUnitSphere(rng *rand.Rand) Vec3 {
	return Normalized(RandomRange(-1, 1, rng))
}
