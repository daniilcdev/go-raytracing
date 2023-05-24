package main

type Ray struct {
	Origin Vec3
	Dir    Vec3
}

func (r *Ray) At(t float64) Vec3 {
	origin := r.Origin
	dir := r.Dir
	return *dir.Scale(t).Add(origin)
}
