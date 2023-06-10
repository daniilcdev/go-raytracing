package main

type Material struct {
	Albedo        Vec3
	Roughness     float64
	EmissionPower float64
	EmissionColor Vec3
}

func (m *Material) GetEmission() Vec3 {
	return Mul(m.EmissionColor, m.EmissionPower)
}

type Sphere struct {
	Center     Vec3
	Radius     float64
	MaterialId int
}

type Scene struct {
	Spheres   []Sphere
	Materials []Material
}
