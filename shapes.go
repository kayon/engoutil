package engoutil

import (
	"strings"

	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
)

type ShapeKind uint16

const (
	SHAPE_KIND_LINE ShapeKind = 1 << iota
	SHAPE_KIND_STIPPLE_LINE
	SHAPE_KIND_RECT
	SHAPE_KIND_STIPPLE_RECT
	SHAPE_KIND_CIRCLE
	SHAPE_KIND_POLYGON
	SHAPE_KIND_CURVE
	SHAPE_KIND_TEXT
	SHAPE_KIND_IMAGE
)

var shapeKindName = [9]string{"Line", "StippleLine", "Rect", "StippleRect", "Circle", "Polygon", "Curve", "Text", "Image"}

func (kind ShapeKind) String() string {
	var s []string
	for i, name := range shapeKindName {
		if kind&(1<<i) > 0 {
			s = append(s, name)
		}
	}
	if len(s) > 0 {
		return strings.Join(s, "|")
	}
	return "invalid shape"
}

type MouseAction uint8

const (
	MOUSE_NONE MouseAction = iota
	MOUSE_CLICKED
	MOUSE_DRAGGED
)

type Point [2]float32
type Points []Point

func (p Points) Points() (points []engo.Point) {
	points = make([]engo.Point, len(p))
	for i, v := range p {
		points[i] = engo.Point{X: v[0], Y: v[1]}
	}
	return
}

func (p Points) Equal(points []engo.Point) bool {
	if len(p) != len(points) {
		return true
	}
	for i, v := range points {
		if v.X != p[i][0] || v.Y != p[i][1] {
			return true
		}
	}
	return false
}

type Shapes []*Shape

// 基础形状
type Shape struct {
	kind   ShapeKind
	Entity *ecs.BasicEntity
	Render *common.RenderComponent
	Space  *common.SpaceComponent

	// Mouseable
	mouseAction MouseAction
	Mouse  *common.MouseComponent

	// attribute 0:x, 1:y
	// Line		2:offsetX, 3:offsetY, 4:sin, 5:cos
	// Rect		2:w, 3:h
	// Circle 	2:radius, 3:arc
	// Polygon	2:w, 3:h
	// Curve	2:w, 3:h
	// Text		2:ax, 3:ay, 4:scale
	attr [6]float32

	onUpdate func(*Shape, float32)
	onHover  [2]func(*Shape)
	onClick  func(*Shape)
	onDrag   func(*Shape, float32, float32)

	// SHAPE_KIND_STIPPLE_LINE, SHAPE_KIND_STIPPLE_RECT
	stipple *Stipple
}

// implementation of common.BasicFace
func (s *Shape) ID() uint64 {
	return s.Entity.ID()
}

// implementation of common.BasicFace
func (s *Shape) GetBasicEntity() *ecs.BasicEntity {
	return s.Entity
}

// implementation of common.RenderFace
func (s *Shape) GetRenderComponent() *common.RenderComponent {
	return s.Render
}

// implementation of common.SpaceFace
func (s *Shape) GetSpaceComponent() *common.SpaceComponent {
	return s.Space
}

// implementation of common.Mouseable
func (s *Shape) GetMouseComponent() *common.MouseComponent {
	return s.Mouse
}

// (*Shape) Hidden
func (s *Shape) Hidden() {
	s.Render.Hidden = true
}

// (*Shape) Show
func (s *Shape) Show() {
	s.Render.Hidden = false
}

// (*Shape) Visible
func (s *Shape) Visible() bool {
	return !s.Render.Hidden
}

// (*Shape) OnUpdate
func (s *Shape) OnUpdate(fn func(*Shape, float32)) {
	s.onUpdate = fn
}

// (*Shape) OnHover
func (s *Shape) OnHover(enter, leave func(*Shape)) {
	if s.Mouse == nil {
		s.Mouse = &common.MouseComponent{}
	}
	s.onHover[0] = enter
	s.onHover[1] = leave
}

// (*Shape) OnClick
func (s *Shape) OnClick(fn func(*Shape)) {
	if s.Mouse == nil {
		s.Mouse = &common.MouseComponent{}
	}
	s.onClick = fn
}

// (*Shape) OnDrag
func (s *Shape) OnDrag(fn func(*Shape, float32, float32)) {
	if s.Mouse == nil {
		s.Mouse = &common.MouseComponent{}
	}
	s.onDrag = fn
}

func (s *Shape) requireKind(kind ShapeKind, method string) bool {
	if kind&s.kind != s.kind {
		warning("(Shape) %s(), %s expected, got %s", method, kind, s.kind)
		return false
	}
	return true
}

func newShape(kind ShapeKind) (s *Shape) {
	entity := ecs.NewBasic()
	s = &Shape{
		kind:   kind,
		Entity: &entity,
		Render: &common.RenderComponent{},
		Space:  &common.SpaceComponent{},
	}
	return
}
