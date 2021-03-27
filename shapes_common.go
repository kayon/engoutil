package engoutil

import (
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
	"github.com/EngoEngine/math"
)

const halfPI = math.Pi / 2

// (*Shape) Transform
// @overload Line.Transform(x1, y1, x2, y2)
// @overload Circle.Transform(cx, cy, radius)
// @overload Text.Transform(x, y)
func (s *Shape) Transform(x, y, width, height float32) {
	switch s.kind {
	case SHAPE_KIND_LINE:
		length := math.Hypot(width-x, height-y)
		theta := math.Atan2(height-y, width-x) - halfPI
		degrees := theta * 180 / math.Pi
		sin, cos := math.Sincos(theta)
		// position
		s.attr[0] = x
		s.attr[1] = y
		// offset: x, y
		s.attr[2] = float32(cos) * s.Space.Width / 2
		s.attr[3] = float32(sin) * s.Space.Width / 2
		s.attr[4] = sin
		s.attr[5] = cos
		x -= s.attr[2]
		y -= s.attr[3]
		s.Space.Position.X = x
		s.Space.Position.Y = y
		s.Space.Height = length
		s.Space.Rotation = degrees
	case SHAPE_KIND_STIPPLE_RECT, SHAPE_KIND_RECT, SHAPE_KIND_POLYGON, SHAPE_KIND_CURVE, SHAPE_KIND_IMAGE:
		s.Space.Position.X = x
		s.Space.Position.Y = y
		s.Space.Width = width
		s.Space.Height = height
		s.attr[0] = x
		s.attr[1] = y
		s.attr[2] = width
		s.attr[3] = height
	case SHAPE_KIND_CIRCLE:
		s.attr[0] = x
		s.attr[1] = y
		s.attr[2] = width
		size := width * 2
		s.Space.Position = engo.Point{X: x - width, Y: y - width}
		s.Space.Width = size
		s.Space.Height = size
	case SHAPE_KIND_TEXT:
		s.attr[0] = x
		s.attr[1] = y
		if t, ok := s.Render.Drawable.(*Text); ok {
			t.Position.X = x - t.Width()/s.attr[4]*s.attr[2]
			t.Position.Y = y - t.Height()/s.attr[4]*s.attr[3]
		}
	}
}

// (*Shape) Move 移动
func (s *Shape) Move(x, y float32) {
	if s.attr[0] == x && s.attr[1] == y {
		return
	}
	s.attr[0] = x
	s.attr[1] = y
	switch s.kind {
	// 圆形移动的是中心点
	case SHAPE_KIND_CIRCLE:
		s.Space.Position.X = x - s.attr[2] // -radius
		s.Space.Position.Y = y - s.attr[2] // -radius
	// 直线移动的是顶点
	case SHAPE_KIND_LINE:
		s.Space.Position.X = x - s.attr[2] // -offsetX
		s.Space.Position.Y = y - s.attr[3] // -offsetY
	case SHAPE_KIND_RECT:
	case SHAPE_KIND_POLYGON:
	case SHAPE_KIND_CURVE:
		s.Space.Position.X = x
		s.Space.Position.Y = y
	case SHAPE_KIND_TEXT:
		s.Transform(x, y, 0, 0)
	}
}

// (*Shape) MoveX 移动 X
func (s *Shape) MoveX(x float32) {
	if s.attr[0] != x {
		s.Move(x, s.attr[1])
	}
}

// (*Shape) MoveY 移动 Y
func (s *Shape) MoveY(y float32) {
	if s.attr[1] != y {
		s.Move(s.attr[0], y)
	}
}

// (*Shape) Rotate 设置旋转角度
func (s *Shape) Rotate(deg float32) {
	s.Space.Rotation = math.Mod(deg, 360)
}

// (*Shape) AddRotate 旋转增加角度
func (s *Shape) AddRotate(deg float32) {
	s.Space.Rotation = math.Mod(s.Space.Rotation + deg, 360)
}

// (*Shape) SetPoints
// 支持这些形状
// SHAPE_KIND_STIPPLE_LINE, SHAPE_KIND_POLYGON, SHAPE_KIND_CURVE
func (s *Shape) SetPoints(points Points) {
	if !s.requireKind(SHAPE_KIND_STIPPLE_LINE|SHAPE_KIND_POLYGON|SHAPE_KIND_CURVE, "SetPoints") {
		return
	}
	if s.Render == nil || s.Render.Drawable == nil {
		return
	}
	switch t := s.Render.Drawable.(type) {
	case StippleLine:
		if !points.Equal(t.Points) {
			t.Points = points.Points()
			s.Render.Drawable = t
		}
	case common.ComplexTriangles:
		if !points.Equal(t.Points) {
			t.Points = points.Points()
			s.Render.Drawable = t
		}
	case common.Curve:
		if !points.Equal(t.Points) {
			t.Points = points.Points()
			s.Render.Drawable = t
		}
	}
}

// (*Shape) SetStrokeWidth
func (s *Shape) SetStrokeWidth(width float32) {
	if s.Render == nil || s.Render.Drawable == nil {
		return
	}
	switch t := s.Render.Drawable.(type) {
	case StippleLine:
		if t.BorderWidth != width {
			t.BorderWidth = width
			s.Render.Drawable = t
		}
	case StippleRect:
		if t.BorderWidth != width {
			t.BorderWidth = width
			s.Render.Drawable = t
		}
	case common.Rectangle:
		switch s.kind {
		case SHAPE_KIND_LINE:
			s.attr[2] = float32(s.attr[4]) * s.Space.Width // sin * width
			s.attr[3] = float32(s.attr[5]) * s.Space.Width // cos * width
			s.Space.Position.X -= s.attr[2]
			s.Space.Position.Y -= s.attr[3]
		case SHAPE_KIND_RECT:
			if t.BorderWidth != width {
				t.BorderWidth = width
				s.Render.Drawable = t
			}
		}
	case common.Circle:
		if t.BorderWidth != width {
			t.BorderWidth = width
			s.Render.Drawable = t
		}
	case common.ComplexTriangles:
		if t.BorderWidth != width {
			t.BorderWidth = width
			s.Render.Drawable = t
		}
	case common.Curve:
		if t.LineWidth != width {
			t.LineWidth = width
			s.Render.Drawable = t
		}
	default:
		warning("Shape(%s) SetStrokeWidth(), type %T not supported", s.kind, t)
	}
}

// (*Shape) SetFillColor
func (s *Shape) SetFillColor(clr uint32) {
	if s.Render == nil {
		return
	}
	s.Render.Color = NewColor(clr)
	if !ColorEqualUint32(s.Render.Color, clr) {
		s.Render.Color = NewColor(clr)
	}
}

// (*Shape) SetStrokeColor
func (s *Shape) SetStrokeColor(clr uint32) {
	if s.kind == SHAPE_KIND_LINE {
		s.SetFillColor(clr)
		return
	}
	if s.Render == nil || s.Render.Drawable == nil {
		return
	}
	switch t := s.Render.Drawable.(type) {
	case StippleLine:
		s.SetFillColor(clr)
	case StippleRect:
		s.SetFillColor(clr)
	case common.Rectangle:
		if !ColorEqualUint32(t.BorderColor, clr) {
			t.BorderColor = NewColor(clr)
			s.Render.Drawable = t
		}
	case common.Circle:
		if !ColorEqualUint32(t.BorderColor, clr) {
			t.BorderColor = NewColor(clr)
			s.Render.Drawable = t
		}
	case common.ComplexTriangles:
		if !ColorEqualUint32(t.BorderColor, clr) {
			t.BorderColor = NewColor(clr)
			s.Render.Drawable = t
		}
	case common.Curve:
		s.SetFillColor(clr)
	case *Text:
		if !t.Color.EqualUint32(clr) {
			t.Color.Set(clr)
		}
	default:
		warning("Shape(%s) SetStrokeColor(), type %T not supported", s.kind, t)
	}
}
