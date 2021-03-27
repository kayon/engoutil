package engoutil

// NewStippleLine
func NewStippleLine(points Points, factor int32, pattern uint16, lineWidth float32, clr uint32) *Shape {
	s := newShape(SHAPE_KIND_STIPPLE_LINE)
	s.stipple = &Stipple{
		Factor:  factor,
		Pattern: pattern,
	}
	s.Render.Drawable = StippleLine{BorderWidth: lineWidth, Points: points.Points(), Stipple: *s.stipple}
	s.Render.Color = NewColor(clr)
	s.Render.SetShader(ShapeHUDShader)
	return s
}

// NewStippleRect
func NewStippleRect(x, y, width, height float32, factor int32, pattern uint16, lineWidth float32, clr uint32) *Shape {
	s := newShape(SHAPE_KIND_STIPPLE_RECT)
	s.stipple = &Stipple{
		Factor:  factor,
		Pattern: pattern,
	}
	s.attr[0] = x
	s.attr[1] = y
	s.attr[2] = width
	s.attr[3] = height
	s.Render.Drawable = StippleRect{BorderWidth: lineWidth, Stipple: *s.stipple}
	s.Render.Color = NewColor(clr)
	s.Render.SetShader(ShapeHUDShader)
	s.Space.Position.X = x
	s.Space.Position.Y = y
	s.Space.Width = width
	s.Space.Height = height
	return s
}

// (*Shape) Stipple
func (s *Shape) Stipple() (int32, uint16) {
	if s.stipple == nil {
		return 0, 0
	}
	return s.stipple.Factor, s.stipple.Pattern
}

// (*Shape) SetStipple
func (s *Shape) SetStipple(factor int32, pattern uint16) {
	if !s.requireKind(SHAPE_KIND_STIPPLE_LINE|SHAPE_KIND_STIPPLE_RECT, "SetStipple") {
		return
	}
	if s.stipple.Factor == factor && s.stipple.Pattern == pattern {
		return
	}
	switch t := s.Render.Drawable.(type) {
	case StippleLine:
		t.Stipple.Factor = factor
		t.Stipple.Pattern = pattern
		s.Render.Drawable = t
	case StippleRect:
		t.Stipple.Factor = factor
		t.Stipple.Pattern = pattern
		s.Render.Drawable = t
	}
	s.stipple.Factor = factor
	s.stipple.Pattern = pattern
}

// (*Shape) MoveStippleLeft
func (s *Shape) MoveStippleLeft() {
	if !s.requireKind(SHAPE_KIND_STIPPLE_LINE|SHAPE_KIND_STIPPLE_RECT, "MoveStippleLeft") {
		return
	}
	if s.stipple.Pattern == 0xFFFF {
		return
	}
	pattern := s.stipple.Pattern
	if pattern&0x8000 > 0 {
		pattern = s.stipple.Pattern<<1 | 1
	} else {
		pattern <<= 1
	}
	s.SetStipple(s.stipple.Factor, pattern)
}

// (*Shape) MoveStippleRight
func (s *Shape) MoveStippleRight() {
	if !s.requireKind(SHAPE_KIND_STIPPLE_LINE|SHAPE_KIND_STIPPLE_RECT, "MoveStippleRight") {
		return
	}
	if s.stipple.Pattern == 0xFFFF {
		return
	}
	pattern := s.stipple.Pattern
	if pattern&0x1 > 0 {
		pattern = pattern>>1 | 0x8000
	} else {
		pattern >>= 1
	}
	s.SetStipple(s.stipple.Factor, pattern)
}
