package main

import (
	"math"

	"github.com/go-gl/mathgl/mgl64"
)

type Camera struct {
	Fov, nearClip, farClip float64

	forwardDirection Vec3
	position         Vec3

	// cached ray directions
	rayDirections []Vec3

	viewportWidth, viewportHeight int

	projection, view, inverseProjection, inverseView mgl64.Mat4
}

func (c *Camera) recalculateView() {
	eye := mgl64.Vec3{c.position.X, c.position.Y, c.position.Z}
	target := Add(c.position, c.forwardDirection)
	eyeCenter := mgl64.Vec3{target.X, target.Y, target.Z}
	c.view = mgl64.LookAtV(
		eye,
		eyeCenter,
		mgl64.Vec3{0, -1, 0},
	)

	c.inverseView = c.view.Inv()
}

func (c *Camera) recalculateProjection() {
	c.projection = mgl64.Perspective(
		(c.Fov*math.Pi)/180,
		float64(c.viewportWidth)/float64(c.viewportHeight),
		c.nearClip,
		c.farClip,
	)

	c.inverseProjection = c.projection.Inv()
}

func (c *Camera) recalculateRayDirections() {
	c.rayDirections = make([]Vec3, c.viewportWidth*c.viewportHeight)

	for y := 0; y < c.viewportHeight; y++ {
		for x := 0; x < c.viewportWidth; x++ {

			coord := mgl64.Vec2{float64(x)/float64(c.viewportWidth)*2 - 1,
				float64(y)/float64(c.viewportHeight)*2 - 1,
			}

			target := c.inverseProjection.Mul4x1(coord.Vec4(1, 1))
			n := target.Vec3()
			for i := range n {
				n[i] /= target.W()
			}

			p := c.inverseView.Mul4x1(n.Normalize().Vec4(0)).Vec3()

			rayDirection := Vec3{p[0], p[1], p[2]}
			c.rayDirections[x+y*c.viewportWidth] = rayDirection
		}
	}

}

func (c *Camera) Resize(w, h int) {
	if w == c.viewportWidth && h == c.viewportHeight {
		return
	}

	c.viewportHeight = h
	c.viewportWidth = w
	c.recalculateProjection()
	c.recalculateRayDirections()
}
