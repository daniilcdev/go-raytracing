package main

import "image/color"

func (v3 *Vec3) ToRGB() color.Color {
	return color.RGBA{R: uint8(v3.X * 255),
		G: uint8(v3.Y * 255),
		B: uint8(v3.Z * 255),
		A: 255}
}
