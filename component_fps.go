package engoutil

import (
	"fmt"

	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
	"github.com/EngoEngine/math"
)

var _ ecs.System = (*FPSComponent)(nil)
var _ ecs.Initializer = (*FPSComponent)(nil)

func NewFPSComponent(x, y, size float32, fgColor, bgColor uint32) (f *FPSComponent) {
	if size < 1 {
		size = 1
	}
	f = &FPSComponent{
		position: engo.Point{X: x, Y: y},
		color:    [2]uint32{fgColor, bgColor},
		size:     size,
		format:   "FPS:%g",
	}
	return
}

type FPSComponent struct {
	position engo.Point
	text     *Shape
	color    [2]uint32
	size     float32
	elapsed  float32
	format   string
	anchor   [2]float32
	hidden   bool
	update   bool
}

func (f *FPSComponent) Hide() {
	f.hidden = true
	f.update = true
}

func (f *FPSComponent) Show() {
	f.hidden = false
	f.update = true
}

func (f *FPSComponent) SetAnchor(ax, ay float32) {
	f.anchor[0], f.anchor[1] = ax, ay
	f.update = true
}

func (f *FPSComponent) Move(x, y float32) {
	f.position.X = x
	f.position.Y = y
	f.update = true
}

func (f *FPSComponent) SetFormat(format string) {
	f.format = format
}

func (*FPSComponent) Remove(ecs.BasicEntity) {}

func (f *FPSComponent) Update(dt float32) {
	if f.update {
		f.text.Move(f.position.X, f.position.Y)
		f.text.SetAnchor(f.anchor[0], f.anchor[1])
		f.text.Render.Hidden = f.hidden
		f.update = false
		f.elapsed = 1
	}

	f.elapsed += dt
	if f.elapsed >= 1 {
		if !f.hidden {
			s := f.String()
			f.text.SetText(s)
		}
		f.elapsed--
	}
}

func (f *FPSComponent) String() string {
	if engo.Time == nil {
		return ""
	}
	return fmt.Sprintf(f.format, engo.Time.FPS())
}

func (f *FPSComponent) New(w *ecs.World) {
	s := f.String()

	f.text = NewTextWithStyle(s, f.position.X, f.position.Y, 0, 0, f.size, f.color[0], &TextStyle{
		URL: Font04b08,
		BG:  f.color[1],
		Padding: Padding{
			Top:    math.Ceil(f.size * 0.25),
			Right:  math.Ceil(f.size * 0.35),
			Bottom: 0,
			Left:   math.Ceil(f.size * 0.35),
		},
	})
	f.text.Render.StartZIndex = 10002

	f.update = true
	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			sys.AddByInterface(f.text)
		}
	}
}
