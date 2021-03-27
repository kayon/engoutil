package main

import (
	"fmt"

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
		Title:        "Component - Loading",
		Fullscreen:   false,
		Width:        400,
		Height:       400,
		NotResizable: true,
		FPSLimit:     60,
		MSAA:         4,
	}, scene)
}

type Scene struct {
	canvas  *engoutil.Canvas
	loading [3]*engoutil.LoadingComponent
	fps     *engoutil.FPSComponent
}

func (scene *Scene) Preload() {
	engoutil.InitShaders()
	scene.canvas = engoutil.NewCanvas()
	scene.loading[0] = engoutil.NewLoadingComponent(200, 200, 100, 0x039175FF, 0xD6D6D6FF)
	scene.loading[1] = engoutil.NewLoadingComponent(100, 150, 50, 0xFB4E66FF, 0xD6D6D6FF)
	scene.loading[2] = engoutil.NewLoadingComponent(300, 150, 50, 0x19C0FBFF, 0xD6D6D6FF)
	scene.fps = engoutil.NewFPSComponent(200, 100, 16, 0xFFFFFFFF, 0)
	scene.fps.SetAnchor(0.5, 0)
}

func (scene *Scene) Type() string {
	return "Scene"
}

func (scene *Scene) Setup(u engo.Updater) {
	common.SetBackground(engoutil.NewColor(0x212121FF))

	w, _ := u.(*ecs.World)

	// Use canvas instead of common.RenderSystem
	w.AddSystem(scene.canvas)
	w.AddSystem(scene.fps)
	for _, item := range scene.loading {
		w.AddSystem(item)
	}

	var pause bool
	var speed float32 = 1.0
	const speedMin = 0.1
	const speedMax = 3.0

	label := engoutil.NewText("Loading spinner", 200, 10, 0.5, 0, 36, 0xFFFFFFFF)
	speedText := engoutil.NewTextWithStyle(fmt.Sprintf("SPEED: %.1f", speed), 200, 70, 0.5, 0, 24, 0xFFFFFFFF, &engoutil.TextStyle{
		URL: engoutil.Font04b08,
	})
	scene.canvas.Push(label, speedText)

	buttonPlay := engoutil.Shapes{
		engoutil.NewCircle(200, 330, 25, 0, 2, 0xFFFFFFFF, 0),
		// pause
		engoutil.NewRect(192, 320, 6, 20, 0, 0, 0xFFFFFFFF),
		engoutil.NewRect(202, 320, 6, 20, 0, 0, 0xFFFFFFFF),
		// play
		engoutil.NewPolygon(194, 320, 16, 20, 0, engoutil.Points{
			{0, 0}, {1, 0.5}, {0, 1},
		}, 0, 0xFFFFFFFF),
	}
	buttonBack := engoutil.Shapes{
		engoutil.NewCircle(135, 330, 25, 0, 2, 0xFFFFFFFF, 0),
		engoutil.NewPolygon(125, 320, 16, 20, 0, engoutil.Points{
			{0.5, 0}, {0, 0.5}, {0.5, 1},
			{1, 0}, {0.5, 0.5}, {1, 1},
		}, 0, 0xFFFFFFFF),
	}
	buttonForward := engoutil.Shapes{
		engoutil.NewCircle(265, 330, 25, 0, 2, 0xFFFFFFFF, 0),
		engoutil.NewPolygon(259, 320, 16, 20, 0, engoutil.Points{
			{0, 0}, {0.5, 0.5}, {0, 1},
			{0.5, 0}, {1, 0.5}, {0.5, 1},
		}, 0, 0xFFFFFFFF),
	}

	buttonPlay[3].Hidden()
	buttonPlay[0].OnHover(func(s *engoutil.Shape) {
		s.SetStrokeColor(0x29CD34FF)
		buttonPlay[1].SetFillColor(0x29CD34FF)
		buttonPlay[2].SetFillColor(0x29CD34FF)
		buttonPlay[3].SetFillColor(0x29CD34FF)
	}, func(s *engoutil.Shape) {
		s.SetStrokeColor(0xFFFFFFFF)
		buttonPlay[1].SetFillColor(0xFFFFFFFF)
		buttonPlay[2].SetFillColor(0xFFFFFFFF)
		buttonPlay[3].SetFillColor(0xFFFFFFFF)
	})
	buttonPlay[0].OnClick(func(*engoutil.Shape) {
		if !pause {
			buttonPlay[1].Hidden()
			buttonPlay[2].Hidden()
			buttonPlay[3].Show()
		} else {
			buttonPlay[1].Show()
			buttonPlay[2].Show()
			buttonPlay[3].Hidden()
		}
		pause = !pause
		for _, item := range scene.loading {
			item.Pause(pause)
		}
	})
	buttonBack[0].OnHover(func(s *engoutil.Shape) {
		// disable
		if speed == speedMin {
			return
		}
		s.SetStrokeColor(0x29CD34FF)
		buttonBack[1].SetFillColor(0x29CD34FF)
	}, func(s *engoutil.Shape) {
		// disable
		if speed == speedMin {
			return
		}
		s.SetStrokeColor(0xFFFFFFFF)
		buttonBack[1].SetFillColor(0xFFFFFFFF)
	})
	buttonForward[0].OnHover(func(s *engoutil.Shape) {
		// disable
		if speed == speedMax {
			return
		}
		s.SetStrokeColor(0x29CD34FF)
		buttonForward[1].SetFillColor(0x29CD34FF)
	}, func(s *engoutil.Shape) {
		// disable
		if speed == speedMax {
			return
		}
		s.SetStrokeColor(0xFFFFFFFF)
		buttonForward[1].SetFillColor(0xFFFFFFFF)
	})

	var buttonBackDisable, buttonForwardDisable bool

	var setButtonDisable = func(btn engoutil.Shapes, disable bool) {
		var stateColor uint32 = 0xFFFFFF50
		if !disable {
			stateColor = 0xFFFFFFFF
		}
		btn[0].SetStrokeColor(stateColor)
		for i := 1; i < len(btn); i++ {
			btn[i].SetFillColor(stateColor)
		}
	}

	buttonBack[0].OnClick(func(*engoutil.Shape) {
		speed -= 0.1
		if speed < speedMin {
			speed = speedMin
		}
		// disable
		if speed == speedMin {
			buttonBackDisable = true
			setButtonDisable(buttonBack, true)
		} else if buttonForwardDisable {
			buttonForwardDisable = false
			setButtonDisable(buttonForward, false)
		}
		for _, item := range scene.loading {
			item.SetSpeed(speed)
		}
		speedText.SetText(fmt.Sprintf("SPEED: %.1f", speed))
	})
	buttonForward[0].OnClick(func(*engoutil.Shape) {
		speed += 0.1
		if speed > speedMax {
			speed = speedMax
		}
		// disable
		if speed == speedMax {
			buttonForwardDisable = true
			setButtonDisable(buttonForward, true)
		} else if buttonBackDisable {
			buttonBackDisable = false
			setButtonDisable(buttonBack, false)
		}
		for _, item := range scene.loading {
			item.SetSpeed(speed)
		}
		speedText.SetText(fmt.Sprintf("SPEED: %.1f", speed))
	})

	// background
	bg := engoutil.Shapes{
		engoutil.NewCircle(200, 200, 60, 0, 0, 0, 0xFFFFFFFF),
		engoutil.NewCircle(100, 150, 35, 0, 0, 0, 0xFFFFFFFF),
		engoutil.NewCircle(300, 150, 35, 0, 0, 0, 0xFFFFFFFF),
	}

	scene.canvas.GroupPush(bg, buttonPlay, buttonBack, buttonForward)

	scene.canvas.Draw()
}
