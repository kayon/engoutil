package engoutil

import (
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
	"github.com/EngoEngine/math"
)

// NewCircle
// cx, cy are the center of the circle
func NewCircle(cx, cy, radius, arc, strokeWidth float32, strokeColor, fillColor uint32) *Shape {
	s := newShape(SHAPE_KIND_CIRCLE)
	arc = math.Mod(arc, 360)
	if arc == 0 {
		arc = 360
	}
	s.attr[0] = cx
	s.attr[1] = cy
	s.attr[2] = radius
	s.attr[3] = arc
	size := radius * 2
	s.Render.Drawable = common.Circle{Arc: arc, BorderWidth: strokeWidth, BorderColor: NewColor(strokeColor)}
	s.Render.Color = NewColor(fillColor)
	s.Render.SetShader(common.LegacyHUDShader)
	s.Space.Position = engo.Point{X: cx - radius, Y: cy - radius}
	s.Space.Width = size
	s.Space.Height = size
	return s
}

// (*Shape) SetRadius
func (s *Shape) SetRadius(radius float32) {
	if !s.requireKind(SHAPE_KIND_CIRCLE, "SetRadius") {
		return
	}
	if s.attr[2] == radius {
		return
	}
	s.Transform(s.attr[0], s.attr[1], radius, 0)
}

// (*Shape) SetArc 设置形状"圆"的弧度 0..360
func (s *Shape) SetArc(arc float32) {
	if !s.requireKind(SHAPE_KIND_CIRCLE, "SetArc") {
		return
	}
	arc = math.Mod(arc, 360)
	if arc == 0 {
		arc = 360
	}
	if s.Render == nil || s.Render.Drawable == nil || s.attr[3] == arc {
		return
	}
	if t, ok := s.Render.Drawable.(common.Circle); ok {
		s.attr[3] = arc
		t.Arc = arc
		s.Render.Drawable = t
	}
}

// (*Shape) AddArc Loop to increase arc
func (s *Shape) AddArc(arc float32) {
	if !s.requireKind(SHAPE_KIND_CIRCLE, "AddArc") {
		return
	}
	arc = math.Mod(s.attr[3]+arc, 360)
	if arc == 0 {
		arc = 360
	}
	if s.Render == nil || s.Render.Drawable == nil || s.attr[3] == arc {
		return
	}
	if t, ok := s.Render.Drawable.(common.Circle); ok {
		s.attr[3] = arc
		t.Arc = arc
		s.Render.Drawable = t
	}
}
