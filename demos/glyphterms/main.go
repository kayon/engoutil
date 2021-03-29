package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"image/color"

	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"

	"github.com/kayon/engoutil"
)

//go:embed TimesNewRoman.ttf
var TimesNewRomanTTF []byte

const fontName = "Times New Roman.ttf"
const example = 'g'

var scene *Scene
var face font.Face

func init() {
	scene = &Scene{}
	ttf, err := freetype.ParseFont(TimesNewRomanTTF)
	if err != nil {
		panic(err)
	}
	face = truetype.NewFace(ttf, &truetype.Options{
		Size: 256,
	})

	if err = engo.Files.LoadReaderData(fontName, bytes.NewReader(TimesNewRomanTTF)); err != nil {
		panic(fmt.Sprintf("unable to load %q! Error was: ", fontName) + err.Error())
	}
}

func main() {
	engo.Run(engo.RunOptions{
		Title:        "glyphterms",
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
	scene.fps = engoutil.NewFPSComponent(390, 10, 8, 0x000000FF, 0)
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

	scene.draw()
}

func (scene *Scene) draw() {
	ascent, descent := float32(face.Metrics().Ascent.Ceil()), float32(face.Metrics().Descent.Ceil())
	lineHeight := ascent + descent
	bounds, adv, _ := face.GlyphBounds(example)
	advance := float32(adv.Ceil())
	leftSide := float32(bounds.Min.X.Floor())
	rightSide := advance - float32(bounds.Max.X.Ceil())
	// charAscent := float32(-bounds.Min.Y.Ceil())
	charDescent := +float32(bounds.Max.Y.Ceil())
	width := float32((bounds.Max.X - bounds.Min.X).Ceil())
	height := float32((bounds.Max.Y - bounds.Min.Y).Ceil())

	boxBounds := [4]float32{
		200 - (width+leftSide+rightSide)/2, 200 - lineHeight/2,
		200 + (width+leftSide+rightSide)/2, 200 + lineHeight/2,
	}
	boxBlankBounds := [4]float32{
		boxBounds[0] + leftSide, boxBounds[1] + 1,
		boxBounds[2] - rightSide, boxBounds[3] - height,
	}

	charBounds := [4]float32{
		boxBounds[0] + leftSide, 200 - height/2 + charDescent,
		boxBounds[0] + leftSide + width + rightSide, 200 + height/2 + charDescent,
	}
	_ = charBounds

	character := engoutil.NewTextWithStyle(string(example), 200, 200, 0.5, 0.5, 256, 0x000000FF, &engoutil.TextStyle{
		URL:     fontName,
		BgStyle: engoutil.BG_FILL_WRAP,
		BG:      0xFA2C4440,
	})

	box := engoutil.NewStippleRect(boxBounds[0], boxBounds[1], boxBounds[2]-boxBounds[0], boxBounds[3]-boxBounds[1], 1, 0xFFFF, 2, 0x000000FF)
	boxLeftSide := engoutil.NewRect(boxBounds[0], boxBounds[1], leftSide, lineHeight, 0, 0, 0x0079BA40)
	boxRightSide := engoutil.NewRect(boxBounds[2]-rightSide, boxBounds[1], rightSide, lineHeight, 0, 0, 0xFFCE4440)
	boxBlank := engoutil.NewRect(boxBlankBounds[0], boxBlankBounds[1], boxBlankBounds[2]-boxBlankBounds[0], boxBlankBounds[3]-boxBlankBounds[1], 1, 0, 0x00000020)
	blankCross := engoutil.NewStippleLine(engoutil.Points{
		{boxBlankBounds[0], boxBlankBounds[1]}, {boxBlankBounds[2], boxBlankBounds[3]},
		{boxBlankBounds[2], boxBlankBounds[1]}, {boxBlankBounds[0], boxBlankBounds[3]},
	}, 1, 0xFFFF, 2, 0x00000080)
	boxBlank.Hidden()
	blankCross.Hidden()

	baseline := engoutil.NewStippleLine(engoutil.Points{
		{75, boxBounds[1] + ascent}, {350, boxBounds[1] + ascent},
	}, 1, 0xF0F0, 2, 0x000000FF)
	origin := engoutil.NewCircle(boxBounds[0]-5-5, boxBounds[1]+ascent, 3, 0, 0, 0, 0x000000FF)
	lineGuideAscent := engoutil.NewStippleLine(engoutil.Points{
		// ascent vertical
		{boxBounds[0] - 5 - 5, boxBounds[1] + 10}, {boxBounds[0] - 5 - 5, boxBounds[1] + ascent - 10},
	}, 1, 0xFFFF, 2, 0x000000FF)
	lineGuideDescent := engoutil.NewStippleLine(engoutil.Points{
		// descent vertical
		{boxBounds[0] - 5 - 5, boxBounds[1] + ascent + 10}, {boxBounds[0] - 5 - 5, boxBounds[3] - 10},
	}, 1, 0xFFFF, 2, 0x000000FF)
	lineGuideLineHeight := engoutil.NewStippleLine(engoutil.Points{
		// line height vertical
		{boxBounds[2] + 5 + 5, boxBounds[1] + 10}, {boxBounds[2] + 5 + 5, boxBounds[3] - 10},
	}, 1, 0xFFFF, 2, 0x000000FF)
	lineGuide := engoutil.NewStippleLine(engoutil.Points{
		// left top
		{boxBounds[0] - 5 - 10, boxBounds[1]}, {boxBounds[0] - 5, boxBounds[1]},
		// left bottom
		{boxBounds[0] - 5 - 10, boxBounds[3]}, {boxBounds[0] - 5, boxBounds[3]},

		// right top
		{boxBounds[2] + 5, boxBounds[1]}, {boxBounds[2] + 15, boxBounds[1]},
		// right bottom
		{boxBounds[2] + 5, boxBounds[3]}, {boxBounds[2] + 15, boxBounds[3]},
	}, 1, 0xFFFF, 2, 0x000000FF)
	arrowAscent := engoutil.Shapes{
		// ascent top
		engoutil.NewPolygon(boxBounds[0]-5-10, boxBounds[1]+2, 10, 10, 0, engoutil.Points{
			{0, 1}, {1, 1}, {0.5, 0},
		}, 0, 0x000000FF),
		// ascent bottom
		engoutil.NewPolygon(boxBounds[0]-5-10, boxBounds[1]+ascent-10-5, 10, 10, 0, engoutil.Points{
			{0, 0}, {1, 0}, {0.5, 1},
		}, 0, 0x000000FF),
	}
	arrowDescent := engoutil.Shapes{
		// descent top
		engoutil.NewPolygon(boxBounds[0]-5-10, boxBounds[1]+ascent+5, 10, 10, 0, engoutil.Points{
			{0, 1}, {1, 1}, {0.5, 0},
		}, 0, 0x000000FF),

		// descent bottom
		engoutil.NewPolygon(boxBounds[0]-5-10, boxBounds[3]-10-2, 10, 10, 0, engoutil.Points{
			{0, 0}, {1, 0}, {0.5, 1},
		}, 0, 0x000000FF),
	}

	arrowLineHeight := engoutil.Shapes{
		// line height top
		engoutil.NewPolygon(boxBounds[2]+5, boxBounds[1]+2, 10, 10, 0, engoutil.Points{
			{0, 1}, {1, 1}, {0.5, 0},
		}, 0, 0x000000FF),
		// line height bottom
		engoutil.NewPolygon(boxBounds[2]+5, boxBounds[3]-10-2, 10, 10, 0, engoutil.Points{
			{0, 0}, {1, 0}, {0.5, 1},
		}, 0, 0x000000FF),
	}

	rectBoundingBox := engoutil.NewRect(10, 10, 30, 20, 0, 0, 0xFA2C4440)
	rectLeftSide := engoutil.NewRect(140, 10, 30, 20, 0, 0, 0x0079BA40)
	rectRightSide := engoutil.NewRect(240, 10, 30, 20, 0, 0, 0xFFCE4440)

	scene.canvas.Push(character, boxLeftSide, boxRightSide, boxBlank, blankCross, box, baseline, origin, lineGuide, rectBoundingBox, rectLeftSide, rectRightSide)
	scene.canvas.Push(lineGuide, lineGuideAscent, lineGuideDescent, lineGuideLineHeight)
	scene.canvas.GroupPush(arrowAscent, arrowDescent, arrowLineHeight)

	box.OnHover(func(shape *engoutil.Shape) {
		// when it's a solid line
		if _, pattern := shape.Stipple(); pattern == 0xFFFF {
			shape.SetStipple(1, 0x00FF)
		} else {
			// dotted line moving animation
			shape.MoveStippleLeft()
		}
		boxBlank.Show()
		blankCross.Show()
		boxLeftSide.SetFillColor(0x0079BAFF)
		boxRightSide.SetFillColor(0xFFCE44FF)
		character.SetStrokeColor(0xFFFFFFFF)
		character.SetFillColor(0xFA2C44FF)
		rectBoundingBox.SetFillColor(0xFA2C44FF)
		rectLeftSide.SetFillColor(0x0079BAFF)
		rectRightSide.SetFillColor(0xFFCE44FF)
	}, func(shape *engoutil.Shape) {
		boxBlank.Hidden()
		blankCross.Hidden()
		boxLeftSide.SetFillColor(0x0079BA40)
		boxRightSide.SetFillColor(0xFFCE4440)
		shape.SetStipple(1, 0xFFFF)
		character.SetStrokeColor(0x000000FF)
		character.SetFillColor(0xFA2C4440)
		rectBoundingBox.SetFillColor(0xFA2C4440)
		rectLeftSide.SetFillColor(0x0079BA40)
		rectRightSide.SetFillColor(0xFFCE4440)
	})

	textBaseline := engoutil.NewText("Baseline", 0, boxBounds[1]+ascent, 0, 0.5, 18, 0x000000FF)
	textAdvance := engoutil.NewText("Advance width", 200, boxBounds[3], 0.5, 0, 18, 0x000000FF)
	textAscent := engoutil.NewText("Ascent", boxBounds[0]-5-10, boxBounds[1]+ascent/2, 1, 0.5, 18, 0x000000FF)
	textDescent := engoutil.NewText("Descent", boxBounds[0]-5-10, boxBounds[3]-descent/2, 1, 0.5, 18, 0x000000FF)
	textLineHeight := engoutil.NewText("Line Height", boxBounds[2]+5+10, boxBounds[1]+lineHeight/2, 0, 0.5, 18, 0x000000FF)

	textBoundingBox := engoutil.NewText("Bounding box", 45, 20, 0, 0.5, 12, 0x000000FF)
	textLeftSide := engoutil.NewTextWithStyle("Left-side\nbearing", 175, 20, 0, 0.5, 12, 0x000000FF, &engoutil.TextStyle{
		LineSpacing: -8,
	})
	textRightSide := engoutil.NewTextWithStyle("Right-side\nbearing", 275, 20, 0, 0.5, 12, 0x000000FF, &engoutil.TextStyle{
		LineSpacing: -8,
	})

	textAscent.OnHover(func(s *engoutil.Shape) {
		s.SetStrokeColor(0xFF0000FF)
		lineGuideAscent.SetFillColor(0xFF0000FF)
		for _, item := range arrowAscent {
			item.SetFillColor(0xFF0000FF)
		}
	}, func(s *engoutil.Shape) {
		s.SetStrokeColor(0x000000FF)
		lineGuideAscent.SetFillColor(0x000000FF)
		for _, item := range arrowAscent {
			item.SetFillColor(0x000000FF)
		}
	})

	textBaseline.OnHover(func(s *engoutil.Shape) {
		s.SetStrokeColor(0xFF0000FF)
		baseline.SetFillColor(0xFF0000FF)
		baseline.MoveStippleLeft()
	}, func(s *engoutil.Shape) {
		s.SetStrokeColor(0x000000FF)
		baseline.SetFillColor(0x000000FF)
	})

	textDescent.OnHover(func(s *engoutil.Shape) {
		s.SetStrokeColor(0xFF0000FF)
		lineGuideDescent.SetFillColor(0xFF0000FF)
		for _, item := range arrowDescent {
			item.SetFillColor(0xFF0000FF)
		}
	}, func(s *engoutil.Shape) {
		s.SetStrokeColor(0x000000FF)
		lineGuideDescent.SetFillColor(0x000000FF)
		for _, item := range arrowDescent {
			item.SetFillColor(0x000000FF)
		}
	})

	textLineHeight.OnHover(func(s *engoutil.Shape) {
		s.SetStrokeColor(0xFF0000FF)
		lineGuideLineHeight.SetFillColor(0xFF0000FF)
		for _, item := range arrowLineHeight {
			item.SetFillColor(0xFF0000FF)
		}
	}, func(s *engoutil.Shape) {
		s.SetStrokeColor(0x000000FF)
		lineGuideLineHeight.SetFillColor(0x000000FF)
		for _, item := range arrowLineHeight {
			item.SetFillColor(0x000000FF)
		}
	})

	scene.canvas.Push(textBaseline, textAdvance, textAscent, textDescent, textLineHeight, textBoundingBox, textLeftSide, textRightSide)
	scene.canvas.Draw()
}
