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
		Title:        "Text",
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
	scene.fps = engoutil.NewFPSComponent(390, 10, 16, 0xFFFFFFFF, 0x002B36FF)
	scene.fps.SetAnchor(1, 0)
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

	var exampleTexts = []string{
		"Hello, world",
		"你好，世界",
		"こんにちは、世界",
		"Здравствуй, мир",
		"Bonjour, le monde",
	}

	labels := make(engoutil.Shapes, len(exampleTexts))
	for i, s := range exampleTexts {
		labels[i] = engoutil.NewTextWithStyle(s, 10, float32(i)*40+10, 0, 0, 24, 0x000000FF, &engoutil.TextStyle{
			BG: 0xDEDEDEFF,
			Padding: engoutil.Padding{
				Top:    0,
				Right:  8,
				Bottom: 8,
				Left:   8,
			},
		})
		labels[i].OnHover(func(s *engoutil.Shape) {
			s.SetFillColor(0x000000FF)
			s.SetStrokeColor(0xFFFFFFFF)
		}, func(s *engoutil.Shape) {
			s.SetFillColor(0xDEDEDEFF)
			s.SetStrokeColor(0x000000FF)
		})
	}
	scene.canvas.GroupPush(labels)

	radio1 := newRadio("BG_FILL_FULL", 10, 220, true)
	radio2 := newRadio("BG_FILL_WRAP", 10, 250, false)
	radio1.onClick = func(r *Radio) {
		radio2.SetSelected(false)
		r.SetSelected(true)
		for _, label := range labels {
			label.SetBgStyle(engoutil.BG_FILL_FULL)
		}
	}
	radio2.onClick = func(r *Radio) {
		radio1.SetSelected(false)
		r.SetSelected(true)
		for _, label := range labels {
			label.SetBgStyle(engoutil.BG_FILL_WRAP)
		}
	}
	scene.canvas.GroupPush(radio1.shapes)
	scene.canvas.GroupPush(radio2.shapes)

	var letterSpacing float32
	textLetterSpacing := engoutil.NewText(fmt.Sprintf("LetterSpacing: %g", letterSpacing), 10, 295, 0, 0.5, 18, 0x000000FF)
	btn1 := newButton("-", 160, 280)
	btn2 := newButton("+", 235, 280)

	btn1.onClick = func() {
		letterSpacing -= 1
		if letterSpacing < -10 {
			letterSpacing = -10
		}
		if letterSpacing == -10 {
			btn1.SetDisable(true)
		} else if btn2.disabled {
			btn2.SetDisable(false)
		}
		for _, label := range labels {
			label.SetLetterSpacing(letterSpacing)
		}
		textLetterSpacing.SetText(fmt.Sprintf("LetterSpacing: %g", letterSpacing))
	}
	btn2.onClick = func() {
		letterSpacing += 1
		if letterSpacing > 10 {
			letterSpacing = 10
		}
		if letterSpacing == 10 {
			btn2.SetDisable(true)

		} else if btn1.disabled {
			btn1.SetDisable(false)
		}
		for _, label := range labels {
			label.SetLetterSpacing(letterSpacing)
		}
		textLetterSpacing.SetText(fmt.Sprintf("LetterSpacing: %g", letterSpacing))
	}
	scene.canvas.Push(textLetterSpacing)
	scene.canvas.GroupPush(btn1.shapes, btn2.shapes)

	radio3 := newRadio("Align left", 10, 325, true)
	radio4 := newRadio("Align right", 10, 355, false)
	radio3.onClick = func(r *Radio) {
		radio4.SetSelected(false)
		r.SetSelected(true)
		for _, label := range labels {
			label.MoveX(10)
			label.SetAnchor(0, 0)
		}
		scene.fps.Move(390, 10)
		scene.fps.SetAnchor(1, 0)
	}
	radio4.onClick = func(r *Radio) {
		radio3.SetSelected(false)
		r.SetSelected(true)
		for _, label := range labels {
			label.MoveX(390)
			label.SetAnchor(1, 0)
		}
		scene.fps.Move(10, 10)
		scene.fps.SetAnchor(0, 0)
	}
	scene.canvas.GroupPush(radio3.shapes)
	scene.canvas.GroupPush(radio4.shapes)

	scene.canvas.Draw()
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
