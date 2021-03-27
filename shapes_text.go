package engoutil

import (
	"github.com/EngoEngine/engo"
)

type TextStyle struct {
	URL           string
	Width         float32
	Height        float32
	BG            uint32
	BgStyle       uint8
	LineSpacing   float32
	LetterSpacing float32
	Padding       Padding
}

var defaultFontURL = FontDroidSans

func SetDefaultFont(urls string) {
	defaultFontURL = urls
}

// NewText use default font
// ax, ay is anchor point, range -1..1, use 0.5 center the text in x, y
func NewText(text string, x, y, ax, ay, size float32, color uint32) *Shape {
	return newText(text, x, y, ax, ay, size, color, nil)
}

func NewTextWithStyle(text string, x, y, ax, ay, size float32, color uint32, style *TextStyle) *Shape {
	return newText(text, x, y, ax, ay, size, color, style)
}

// (*Shape) SetText 改变文本
func (s *Shape) SetText(text string) {
	if !s.requireKind(SHAPE_KIND_TEXT, "SetText") {
		return
	}
	if s.Render == nil || s.Render.Drawable == nil {
		return
	}
	if t, ok := s.Render.Drawable.(*Text); ok {
		t.Text = text
		if s.attr[2] != 0 {
			s.SetAnchor(s.attr[2], s.attr[3])
		}
	}
}

// (*Shape) SetAnchor 设置锚点, 范围 -1..1
func (s *Shape) SetAnchor(ax, ay float32) {
	if !s.requireKind(SHAPE_KIND_TEXT, "SetAnchor") {
		return
	}
	if s.Render == nil || s.Render.Drawable == nil {
		return
	}
	if t, ok := s.Render.Drawable.(*Text); ok {
		s.attr[2] = ax
		s.attr[3] = ay
		t.Position.X = s.attr[0] - t.Width()/s.attr[4]*s.attr[2]
		t.Position.Y = s.attr[1] - t.Height()/s.attr[4]*s.attr[3]
	}
}

func (s *Shape) SetBgStyle(style uint8) {
	if !s.requireKind(SHAPE_KIND_TEXT, "SetBgStyle") {
		return
	}
	if s.Render == nil || s.Render.Drawable == nil {
		return
	}
	if t, ok := s.Render.Drawable.(*Text); ok {
		t.BgStyle = style
	}
}

func (s *Shape) SetLetterSpacing(spacing float32) {
	if !s.requireKind(SHAPE_KIND_TEXT, "SetLetterSpacing") {
		return
	}
	if s.Render == nil || s.Render.Drawable == nil {
		return
	}
	if t, ok := s.Render.Drawable.(*Text); ok {
		t.LetterSpacing = spacing
		s.SetAnchor(s.attr[2], s.attr[3])
	}
}

func (s *Shape) SetLineSpacing(spacing float32) {
	if !s.requireKind(SHAPE_KIND_TEXT, "SetLineSpacing") {
		return
	}
	if s.Render == nil || s.Render.Drawable == nil {
		return
	}
	if t, ok := s.Render.Drawable.(*Text); ok {
		t.LineSpacing = spacing
		s.SetAnchor(s.attr[2], s.attr[3])
	}
}

func newText(text string, x, y, ax, ay, size float32, color uint32, style *TextStyle) *Shape {
	if style == nil {
		style = &TextStyle{}
	}
	if style.URL == "" {
		style.URL = defaultFontURL
	}
	s := newShape(SHAPE_KIND_TEXT)
	s.attr[0] = x
	s.attr[1] = y
	s.attr[2] = ax
	s.attr[3] = ay
	s.attr[4] = engo.CanvasScale()

	f := &Font{
		URL:   style.URL,
		Size:  float64(size),
		scale: s.attr[4],
	}

	if err := f.CreatePreloaded(); err != nil {
		warning("%q CreatePreloaded. Error was: %s", style.URL, err.Error())
	}

	t := &Text{
		Font:          f,
		Text:          text,
		LineSpacing:   style.LineSpacing,
		LetterSpacing: style.LetterSpacing,
		Position:      engo.Point{X: x, Y: y},
		Color:         NewColor(color),
		BgStyle:       style.BgStyle,
		Padding:       style.Padding,
		width:         style.Width,
		height:        style.Height,
	}
	s.Render.Drawable = t
	s.Render.Color = NewColor(style.BG)
	s.Render.Scale.X = 1 / s.attr[4]
	s.Render.Scale.Y = 1 / s.attr[4]
	s.Render.SetShader(TextHUDShader)

	// Move position to anchor point
	t.Position.X = x - t.Width()/s.attr[4]*ax
	t.Position.Y = y - t.Height()/s.attr[4]*ay

	return s
}
