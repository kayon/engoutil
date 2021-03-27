package engoutil

import (
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
)

func NewCurve(x, y, width, height, strokeWidth float32, points Points, clr uint32) *Shape {
	s := newShape(SHAPE_KIND_CURVE)
	s.attr[0] = x
	s.attr[1] = y
	s.attr[2] = width
	s.attr[3] = height
	s.Render.Drawable = common.Curve{LineWidth: strokeWidth, Points: points.Points()}
	s.Render.Color = NewColor(clr)
	s.Render.SetShader(common.LegacyHUDShader)
	s.Space.Position = engo.Point{X: x, Y: y}
	s.Space.Width = width
	s.Space.Height = height
	return s
}
