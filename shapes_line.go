package engoutil

import (
	"github.com/EngoEngine/engo/common"
)

func NewLine(x1, y1, x2, y2, lineWidth float32, clr uint32) *Shape {
	s := newShape(SHAPE_KIND_LINE)
	s.Render.Drawable = common.Rectangle{}
	s.Render.Color = NewColor(clr)
	s.Render.SetShader(common.LegacyHUDShader)
	s.Space.Width = lineWidth
	s.Transform(x1, y1, x2, y2)
	return s
}
