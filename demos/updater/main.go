package main

import (
	"fmt"
	"image/color"

	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"

	"github.com/kayon/engoutil"
)

var scene *Scene

func init() {
	scene = &Scene{}
}

func main() {
	engo.Run(engo.RunOptions{
		Title:        "Updater",
		Fullscreen:   false,
		Width:        400,
		Height:       400,
		NotResizable: true,
		FPSLimit:     60,
		MSAA:         4,
	}, scene)
}

type Scene struct {
	canvas *engoutil.Canvas
	fps    *engoutil.FPSComponent
}

func (scene *Scene) Preload() {
	engoutil.InitShaders()
	scene.canvas = engoutil.NewCanvas()
	scene.fps = engoutil.NewFPSComponent(10, 10, 16, 0x00212AFF, 0xFFCE42FF)
}

func (scene *Scene) Type() string {
	return "Scene"
}

func (scene *Scene) Setup(u engo.Updater) {
	common.SetBackground(color.White)

	w, _ := u.(*ecs.World)

	// Use canvas instead of common.RenderSystem
	w.AddSystem(scene.canvas)
	w.AddSystem(scene.fps)

	switchBtn := NewSwitch(185, 100, 30, 0xBABABAFF, 0x6ED55DFF)
	switchBtn.SetValue(true)
	switchBtn.onChange = func(state bool) {
		if state {
			scene.fps.Show()
		} else {
			scene.fps.Hide()
		}
	}
	scene.canvas.GroupPush(switchBtn.shapes)

	var duration float32 = 200
	textDuration := engoutil.NewText(fmt.Sprintf("Animation duration: %g msec", duration), 200, 160, 0.5, 0, 16, 0x000000FF)
	btn1 := newButton("-", 135, 200)
	btn2 := newButton("+", 205, 200)

	btn1.onClick = func() {
		duration -= 100
		if duration < 100 {
			duration = 100
		}
		if duration == 100 {
			btn1.SetDisable(true)
		} else if btn2.disabled {
			btn2.SetDisable(false)
		}
		switchBtn.SetAnimationDuration(duration)
		textDuration.SetText(fmt.Sprintf("Animation duration: %g msec", duration))
	}
	btn2.onClick = func() {
		duration += 100
		if duration > 2000 {
			duration = 2000
		}
		if duration == 2000 {
			btn2.SetDisable(true)

		} else if btn1.disabled {
			btn1.SetDisable(false)
		}
		switchBtn.SetAnimationDuration(duration)
		textDuration.SetText(fmt.Sprintf("Animation duration: %g msec", duration))
	}
	scene.canvas.Push(textDuration)
	scene.canvas.GroupPush(btn1.shapes, btn2.shapes)

	radio1 := newRadio("Disable", 115, 260, false)
	radio2 := newRadio("Enable", 205, 260, true)
	radio1.onClick = func(r *Radio) {
		radio2.SetSelected(false)
		r.SetSelected(true)
		switchBtn.SetDisable(true)
	}
	radio2.onClick = func(r *Radio) {
		radio1.SetSelected(false)
		r.SetSelected(true)
		switchBtn.SetDisable(false)
	}
	scene.canvas.GroupPush(radio1.shapes)
	scene.canvas.GroupPush(radio2.shapes)

	scene.canvas.Draw()
}

type Switch struct {
	shapes      engoutil.Shapes
	fg, bg      uint32
	left, right float32
	anim        struct {
		step      float32
		direction float32
		x         float32
	}
	disable bool
	value   bool
	hidden  bool
	onChange func(state bool)
}

func (s *Switch) SetAnimationDuration(msec float32) {
	// Moving distance per millisecond
	s.anim.step = (s.right - s.left) / msec
}

func (s *Switch) Value() bool {
	return s.value
}

// (*Switch) SetValue 设置值
func (s *Switch) SetValue(v bool) {
	if s.value != v {
		s.value = v
		s.update()
		if s.onChange != nil {
			s.onChange(s.value)
		}
	}
}

func (s *Switch) SetDisable(state bool) {
	if s.disable != state {
		s.disable = state
		bg := s.fg
		fg := uint32(0xFFFFFFFF)
		if s.value {
			bg = s.bg
		}
		if s.disable {
			// improve opacity
			bg &= 0xFFFFFF00 + 0x80
			fg &= 0xFFFFFF00 + 0x80
		}
		s.shapes[0].SetFillColor(bg)
		s.shapes[1].SetFillColor(bg)
		s.shapes[2].SetFillColor(bg)
		s.shapes[3].SetFillColor(fg)
	}
}

func (s *Switch) update() {
	if s.value {
		s.shapes[0].SetFillColor(s.bg)
		s.shapes[1].SetFillColor(s.bg)
		s.shapes[2].SetFillColor(s.bg)
		s.anim.direction = 1
	} else {
		s.shapes[0].SetFillColor(s.fg)
		s.shapes[1].SetFillColor(s.fg)
		s.shapes[2].SetFillColor(s.fg)
		s.anim.direction = -1
	}
}

func NewSwitch(x, y, size float32, fg, bg uint32) (s *Switch) {
	s = &Switch{fg: fg, bg: bg}
	radius := size / 2
	width := radius * 2
	overstep := radius * 0.2
	barRadius := radius - overstep
	s.left = x
	s.right = x + width

	s.shapes = engoutil.Shapes{
		// center rect
		engoutil.NewRect(x, y, width, size, 0, 0, s.fg),
		// left round
		engoutil.NewCircle(x, y+radius, radius, 180, 0, 0, s.fg),
		// right round
		engoutil.NewCircle(x+width, y+radius, radius, 180, 0, 0, s.fg),
		// circle
		engoutil.NewCircle(x, y+radius, barRadius, 0, 0, 0, 0xFFFFFFFF),
		// full size rectangle for clicking
		engoutil.NewRect(x-radius, y, width+radius*2, size, 0, 0, 0),
	}

	// left round
	s.shapes[1].Space.Rotation = 90
	s.shapes[1].Space.SetCenter(engo.Point{X: x, Y: y + radius})
	// right round
	s.shapes[2].Space.Rotation = -90
	s.shapes[2].Space.SetCenter(engo.Point{X: x + width, Y: y + radius})

	s.SetAnimationDuration(200)
	s.anim.x = x
	// Circle moving animation
	s.shapes[3].OnUpdate(func(circle *engoutil.Shape, dt float32) {
		if s.anim.direction != 0 {
			s.anim.x += s.anim.step * dt * 1000 * s.anim.direction
			switch s.anim.direction {
			case 1:
				if s.anim.x >= s.right {
					s.anim.direction = 0
					s.anim.x = s.right
				}
			case -1:
				if s.anim.x <= s.left {
					s.anim.direction = 0
					s.anim.x = s.left
				}
			}
			circle.MoveX(s.anim.x)
		}
	})
	// Full size rectangle
	s.shapes[4].OnClick(func(*engoutil.Shape) {
		if !s.disable {
			s.SetValue(!s.value)
		}
	})
	return
}

type Radio struct {
	shapes   engoutil.Shapes
	selected bool
	onClick  func(r *Radio)
}

func (r *Radio) SetSelected(state bool) {
	r.selected = state
	if state {
		r.shapes[1].Show()
	} else {
		r.shapes[1].Hidden()
	}
}

func newRadio(text string, x, y float32, selected bool) *Radio {
	radio := &Radio{
		shapes: engoutil.Shapes{
			engoutil.NewCircle(x+10, y+10, 8, 0, 2, 0x000000FF, 0),
			engoutil.NewCircle(x+10, y+10, 5, 0, 0, 0, 0x000000FF),
			engoutil.NewText(text, x+24, y+10, 0, 0.5, 16, 0x000000FF),
		},
		selected: selected,
	}
	radio.SetSelected(selected)
	radio.shapes[0].OnHover(func(s *engoutil.Shape) {
		if radio.selected {
			return
		}
		s.SetStrokeColor(0x999999FF)
		radio.shapes[2].SetStrokeColor(0x999999FF)
	}, func(s *engoutil.Shape) {
		s.SetStrokeColor(0x000000FF)
		radio.shapes[2].SetStrokeColor(0x000000FF)
	})
	radio.shapes[2].OnHover(func(s *engoutil.Shape) {
		if radio.selected {
			return
		}
		s.SetStrokeColor(0x999999FF)
		radio.shapes[0].SetStrokeColor(0x999999FF)
	}, func(s *engoutil.Shape) {
		s.SetStrokeColor(0x000000FF)
		radio.shapes[0].SetStrokeColor(0x000000FF)
	})
	radio.shapes[0].OnClick(func(*engoutil.Shape) {
		if radio.selected || radio.onClick == nil {
			return
		}
		radio.shapes[0].SetStrokeColor(0x000000FF)
		radio.shapes[2].SetStrokeColor(0x000000FF)
		radio.onClick(radio)
	})
	radio.shapes[2].OnClick(func(*engoutil.Shape) {
		if radio.selected || radio.onClick == nil {
			return
		}
		radio.shapes[0].SetStrokeColor(0x000000FF)
		radio.shapes[2].SetStrokeColor(0x000000FF)
		radio.onClick(radio)
	})
	return radio
}

type Button struct {
	shapes   engoutil.Shapes
	onClick  func()
	disabled bool
}

func (btn *Button) SetDisable(state bool) {
	btn.disabled = state
	if state {
		btn.shapes[1].SetStipple(1, 0xFFFF)
		btn.shapes[0].SetStrokeColor(0xCCCCCCFF)
		btn.shapes[1].SetStrokeColor(0xCCCCCCFF)
		btn.shapes[0].SetFillColor(0)
	} else {
		btn.shapes[0].SetStrokeColor(0x000000FF)
		btn.shapes[1].SetStrokeColor(0x000000FF)
	}
}

func newButton(text string, x, y float32) *Button {
	btn := &Button{
		shapes: engoutil.Shapes{
			engoutil.NewTextWithStyle(text, x+30, y+15, 0.5, 0.5, 18, 0x000000FF, &engoutil.TextStyle{
				// Manually set the size to fill the background color
				Width:  60,
				Height: 30,
			}),
			engoutil.NewStippleRect(x, y, 60, 30, 1, 0xFFFF, 2, 0x000000FF),
		},
	}
	btn.shapes[1].OnHover(func(s *engoutil.Shape) {
		if btn.disabled {
			return
		}
		if _, pattern := s.Stipple(); pattern == 0xFFFF {
			s.SetStipple(1, 0xFF00)
		} else {
			s.MoveStippleRight()
		}
		btn.shapes[0].SetFillColor(0xEEEEEEFF)
	}, func(s *engoutil.Shape) {
		s.SetStipple(1, 0xFFFF)
		btn.shapes[0].SetFillColor(0)
	})
	btn.shapes[1].OnClick(func(*engoutil.Shape) {
		if btn.disabled || btn.onClick == nil {
			return
		}
		btn.onClick()
	})
	return btn
}