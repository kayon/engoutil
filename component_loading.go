package engoutil

import (
	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
	"github.com/EngoEngine/math"
)

var _ ecs.System = (*LoadingComponent)(nil)
var _ ecs.Initializer = (*LoadingComponent)(nil)

func NewLoadingComponent(x, y, size float32, fgColor, bgColor uint32) *LoadingComponent {
	return &LoadingComponent{
		position:    engo.Point{X: x, Y: y},
		size:        size,
		fgColor:     fgColor,
		bgColor:     bgColor,
		arcSpeed:    1.0,
		rotateSpeed: 0.5,
		increase:    1.0,
	}
}

type LoadingComponent struct {
	items       [2]*Shape
	position    engo.Point
	size        float32
	fgColor     uint32
	bgColor     uint32
	radian      float32
	increase    float32
	arcSpeed    float32
	rotateSpeed float32
	paused      bool
	update      bool
	hidden      bool
}

func (l *LoadingComponent) Move(x, y float32) {
	if l.position.X != x || l.position.Y != y {
		l.position.X = x
		l.position.Y = y
		l.update = true
	}
}

func (l *LoadingComponent) Speed() float32 {
	return l.arcSpeed
}

func (l *LoadingComponent) SetSpeed(speed float32) {
	l.arcSpeed = speed
	l.rotateSpeed = speed / 2
}

func (l *LoadingComponent) Pause(state bool) {
	l.paused = state
}

func (l *LoadingComponent) Toggle() {
	l.paused = !l.paused
}

func (l *LoadingComponent) State() bool {
	return !l.paused
}

func (l *LoadingComponent) Hide() {
	l.hidden = true
	l.update = true
}

func (l *LoadingComponent) Show() {
	l.hidden = false
	l.update = true
}

// Implement System interface
func (*LoadingComponent) Remove(ecs.BasicEntity) {}

// Implement System interface
func (l *LoadingComponent) Update(dt float32) {
	if l.update {
		for _, item := range l.items {
			item.Render.Hidden = l.hidden
			if !l.hidden {
				item.Move(l.position.X, l.position.Y)
			}
		}
		l.update = false
	}
	if !l.hidden && !l.paused {
		if l.radian >= 350 {
			l.increase = -1.0
		} else if l.radian <= 10 {
			l.increase = 1.0
		}
		l.radian += math.Mod(dt * l.arcSpeed * 360 * l.increase, 360)
		l.items[0].SetArc(l.radian)
		if l.increase > 0 {
			l.items[0].AddRotate(dt * l.rotateSpeed * 360)
		} else {
			// 当弧度减少时, 2转速
			l.items[0].AddRotate(dt * l.arcSpeed * 2 * 360)
		}
		l.items[0].Space.SetCenter(l.position)
	}
}

// Implement Initializer interface
func (l *LoadingComponent) New(w *ecs.World) {
	// FG
	l.items[0] = NewCircle(l.position.X, l.position.Y, l.size*0.5, 350, l.size*0.1, l.fgColor, 0)
	// BG
	l.items[1] = NewCircle(l.position.X, l.position.Y, l.size*0.5, 0, l.size*0.1, l.bgColor, 0)

	l.items[0].Render.SetZIndex(9997)
	l.items[1].Render.SetZIndex(9996)

	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			for _, item := range l.items {
				sys.AddByInterface(item)
			}
		}
	}
}
