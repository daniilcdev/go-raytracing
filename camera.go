package main

type Camera struct {
	Origin     Vec3
	LLC        Vec3
	Horizontal Vec3
	Vertical   Vec3
}

func NewCamera() Camera {
	const aspect = 16.0 / 9.0
	const viewportH = 2.0
	const viewportW = aspect * viewportH
	const focalLength = 1.0

	origin := Vec3{}
	horizontal := Vec3{X: float64(viewportW)}
	vertical := Vec3{Y: float64(viewportH)}

	llc := Subtract(origin, Mul(horizontal, 0.5))
	llc = Subtract(llc, Mul(vertical, 0.5))
	llc = Subtract(llc, Vec3{Z: focalLength})

	return Camera{
		Origin:     origin,
		LLC:        llc,
		Horizontal: horizontal,
		Vertical:   vertical,
	}
}

func (cam *Camera) getRay(u, v float64) Ray {
	dir := Subtract(Add(Add(cam.LLC, Mul(cam.Horizontal, u)), Mul(cam.Vertical, v)), cam.Origin)
	ray := Ray{Origin: cam.Origin, Dir: dir}
	return ray
}
