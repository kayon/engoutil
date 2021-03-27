package engoutil

import (
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
)

// NewRect
func NewRect(x, y, width, height, strokeWidth float32, strokeColor, fillColor uint32) *Shape {
	s := newShape(SHAPE_KIND_RECT)
	s.attr[0] = x
	s.attr[1] = y
	s.attr[2] = width
	s.attr[3] = height
	s.Render.Drawable = common.Rectangle{BorderWidth: strokeWidth, BorderColor: NewColor(strokeColor)}
	s.Render.Color = NewColor(fillColor)
	s.Render.SetShader(common.LegacyHUDShader)
	s.Space.Position = engo.Point{X: x, Y: y}
	s.Space.Width = width
	s.Space.Height = height
	return s
}
