package main

type Material struct {
	Albedo        Vec3
	Roughness     float64
	EmissionPower float64
	EmissionColor Vec3
}

func (m *Material) GetEmmision() Vec3 {
	return Mul(m.EmissionColor, m.EmissionPower)
}
