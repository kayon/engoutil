package engoutil

import (
	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
)

var _ ecs.System = (*Canvas)(nil)
var _ ecs.Initializer = (*Canvas)(nil)

func NewCanvas() *Canvas {
	return &Canvas{}
}

// Using Canvas instead of RenderSystem
type Canvas struct {
	render  *common.RenderSystem
	mouse   *common.MouseSystem
	objects Shapes
	ready   Shapes
	ids     map[uint64]struct{}
	refresh bool
}

func (c *Canvas) Remove(basic ecs.BasicEntity) {
	delete(c.ids, basic.ID())
	c.render.Remove(basic)
	c.mouse.Remove(basic)
}

func (c *Canvas) New(w *ecs.World) {
	c.render = &common.RenderSystem{}
	c.mouse = &common.MouseSystem{}
	c.objects = make(Shapes, 0, 512)
	c.ids = make(map[uint64]struct{})
	w.AddSystem(c.render)
	w.AddSystem(c.mouse)
}

func (c *Canvas) Update(dt float32) {
	if c.refresh {
		n := len(c.objects)
		for i, v := range c.ready {
			c.ids[v.Entity.ID()] = struct{}{}
			v.Render.SetZIndex(v.Render.StartZIndex + float32(i+n))
			c.render.AddByInterface(v)
			if v.onClick != nil || v.onHover[0] != nil || v.onDrag != nil {
				c.mouse.AddByInterface(v)
			}
			c.objects = append(c.objects, v)
		}
		c.ready = make(Shapes, 256)
		c.refresh = false
	}

	for _, v := range c.objects {
		if v.onUpdate != nil {
			v.onUpdate(v, dt)
		}
		if v.Render.Hidden {
			continue
		}
		if v.onHover[0] != nil {
			if v.Mouse.Hovered {
				v.onHover[0](v)
			} else if v.Mouse.Leave {
				v.onHover[1](v)
			}
		}

		if v.onClick != nil || v.onDrag != nil {
			if v.Mouse.Clicked {
				v.mouseAction = MOUSE_CLICKED
			} else {
				if v.Mouse.Dragged {
					v.mouseAction = MOUSE_DRAGGED
					if v.onDrag != nil {
						v.onDrag(v, engo.Input.Mouse.X-v.Space.Position.X, engo.Input.Mouse.Y-v.Space.Position.Y)
					}
				} else if v.Mouse.Released {
					if v.mouseAction != MOUSE_DRAGGED && v.onClick != nil {
						v.onClick(v)
					}
					v.mouseAction = MOUSE_NONE
				}
			}
		}
	}
}

func (c *Canvas) Push(shapes ...*Shape) {
	for _, s := range shapes {
		if s == nil {
			continue
		}
		if _, ok := c.ids[s.Entity.ID()]; !ok {
			c.ready = append(c.ready, s)
		}
	}
}

func (c *Canvas) GroupPush(groups ...Shapes) {
	for _, g := range groups {
		c.Push(g...)
	}
}

// After that, the contents of push will be rendered in the next frame
func (c *Canvas) Draw() {
	c.refresh = true
}
