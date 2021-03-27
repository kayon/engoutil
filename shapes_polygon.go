package engoutil

import (
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
)

// NewPolygon this is ComplexTriangles
func NewPolygon(x, y, width, height float32, strokeWidth float32, points Points, strokeColor, fillColor uint32) *Shape {
	s := newShape(SHAPE_KIND_POLYGON)
	s.attr[0] = x
	s.attr[1] = y
	s.attr[2] = width
	s.attr[3] = height
	s.Render.Drawable = common.ComplexTriangles{Points: points.Points(), BorderWidth: strokeWidth, BorderColor: NewColor(strokeColor)}
	s.Render.Color = NewColor(fillColor)
	s.Render.SetShader(common.LegacyHUDShader)
	s.Space.Position = engo.Point{X: x, Y: y}
	s.Space.Width = width
	s.Space.Height = height
	return s
}
