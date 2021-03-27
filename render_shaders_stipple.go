package engoutil

import (
	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
	"github.com/EngoEngine/gl"
	"github.com/EngoEngine/math"

	gl2 "github.com/go-gl/gl/v2.1/gl"
)

type shapeShader struct {
	program *gl.Program

	inPosition int

	matrixProjection *gl.UniformLocation
	matrixView       *gl.UniformLocation
	matrixModel      *gl.UniformLocation
	uf_Color         *gl.UniformLocation

	projectionMatrix []float32
	viewMatrix       []float32
	modelMatrix      []float32

	camera        *common.CameraSystem
	cameraEnabled bool

	lastBuffer *gl.Buffer
}

func (l *shapeShader) Setup(w *ecs.World) error {
	var err error
	l.program, err = common.LoadShader(`
attribute vec2 in_Position;

uniform mat3 matrixProjection;
uniform mat3 matrixView;
uniform mat3 matrixModel;

void main() {

  vec3 matr = matrixProjection * matrixView * matrixModel * vec3(in_Position, 1.0);
  gl_Position = vec4(matr.xy, 0, matr.z);
}
`, `
#ifdef GL_ES
#define LOWP lowp
precision mediump float;
#else
#define LOWP
#endif

uniform vec4 uf_Color;

void main (void) {
  gl_FragColor = uf_Color;
}`)

	if err != nil {
		return err
	}

	// Define things that should be read from the texture buffer
	l.inPosition = engo.Gl.GetAttribLocation(l.program, "in_Position")

	// Define things that should be set per draw
	l.matrixProjection = engo.Gl.GetUniformLocation(l.program, "matrixProjection")
	l.matrixView = engo.Gl.GetUniformLocation(l.program, "matrixView")
	l.matrixModel = engo.Gl.GetUniformLocation(l.program, "matrixModel")
	l.uf_Color = engo.Gl.GetUniformLocation(l.program, "uf_Color")

	l.projectionMatrix = make([]float32, 9)
	l.projectionMatrix[8] = 1

	l.viewMatrix = make([]float32, 9)
	l.viewMatrix[0] = 1
	l.viewMatrix[4] = 1
	l.viewMatrix[8] = 1

	l.modelMatrix = make([]float32, 9)
	l.modelMatrix[0] = 1
	l.modelMatrix[4] = 1
	l.modelMatrix[8] = 1

	return nil
}

func (l *shapeShader) Pre() {
	engo.Gl.Enable(engo.Gl.BLEND)
	// TODO: modify in the future, engo.Gl not have LINE_STIPPLE
	gl2.Enable(gl2.LINE_STIPPLE)
	engo.Gl.BlendFunc(engo.Gl.SRC_ALPHA, engo.Gl.ONE_MINUS_SRC_ALPHA)

	// Bind shader and buffer, enable attributes
	engo.Gl.UseProgram(l.program)
	engo.Gl.EnableVertexAttribArray(l.inPosition)

	if engo.ScaleOnResize() {
		l.projectionMatrix[0] = 1 / (engo.GameWidth() / 2)
		l.projectionMatrix[4] = 1 / (-engo.GameHeight() / 2)
	} else {
		l.projectionMatrix[0] = 1 / (engo.CanvasWidth() / (2 * engo.CanvasScale()))
		l.projectionMatrix[4] = 1 / (-engo.CanvasHeight() / (2 * engo.CanvasScale()))
	}

	if l.cameraEnabled {
		l.viewMatrix[1], l.viewMatrix[0] = math.Sincos(l.camera.Angle() * math.Pi / 180)
		l.viewMatrix[3] = -l.viewMatrix[1]
		l.viewMatrix[4] = l.viewMatrix[0]
		l.viewMatrix[6] = -l.camera.X()
		l.viewMatrix[7] = -l.camera.Y()
		l.viewMatrix[8] = l.camera.Z()
	} else {
		l.viewMatrix[6] = -1 / l.projectionMatrix[0]
		l.viewMatrix[7] = 1 / l.projectionMatrix[4]
	}

	engo.Gl.UniformMatrix3fv(l.matrixProjection, false, l.projectionMatrix)
	engo.Gl.UniformMatrix3fv(l.matrixView, false, l.viewMatrix)
}

func (l *shapeShader) updateBuffer(ren *common.RenderComponent, space *common.SpaceComponent) {
	if len(ren.BufferContent) == 0 {
		ren.BufferContent = make([]float32, l.computeBufferSize(ren.Drawable)) // because we add at most this many elements to it
	}

	if changed := l.generateBufferContent(ren, space, ren.BufferContent); !changed {
		return
	}

	if ren.Buffer == nil {
		ren.Buffer = engo.Gl.CreateBuffer()
	}
	engo.Gl.BindBuffer(engo.Gl.ARRAY_BUFFER, ren.Buffer)
	engo.Gl.BufferData(engo.Gl.ARRAY_BUFFER, ren.BufferContent, engo.Gl.STATIC_DRAW)
}

func (l *shapeShader) computeBufferSize(draw common.Drawable) int {
	switch shape := draw.(type) {
	case StippleLine:
		return len(shape.Points) * 2
	case StippleRect:
		return 16
	default:
		return 0
	}
}

func (l *shapeShader) generateBufferContent(ren *common.RenderComponent, space *common.SpaceComponent, buffer []float32) bool {
	var changed bool

	switch shape := ren.Drawable.(type) {
	case StippleLine:
		for i, v := range shape.Points {
			setBufferValue(buffer, i*2, v.X, &changed)
			setBufferValue(buffer, i*2+1, v.Y, &changed)
		}
	case StippleRect:
		w, h := space.Width, space.Height
		// top
		// setBufferValue(buffer, 0, 0, &changed)
		// setBufferValue(buffer, 1, 0, &changed)

		setBufferValue(buffer, 2, w, &changed)
		// setBufferValue(buffer, 3, 0, &changed)

		// right
		setBufferValue(buffer, 4, w, &changed)
		// setBufferValue(buffer, 5, 0, &changed)

		setBufferValue(buffer, 6, w, &changed)
		setBufferValue(buffer, 7, h, &changed)

		// bottom
		setBufferValue(buffer, 8, w, &changed)
		setBufferValue(buffer, 9, h, &changed)

		// setBufferValue(buffer, 10, 0, &changed)
		setBufferValue(buffer, 11, h, &changed)

		// left
		// setBufferValue(buffer, 12, 0, &changed)
		setBufferValue(buffer, 13, h, &changed)

		// setBufferValue(buffer, 14, 0, &changed)
		// setBufferValue(buffer, 15, 0, &changed)
	}

	return changed
}

func (l *shapeShader) Draw(ren *common.RenderComponent, space *common.SpaceComponent) {
	if l.lastBuffer != ren.Buffer || ren.Buffer == nil {
		l.updateBuffer(ren, space)

		engo.Gl.BindBuffer(engo.Gl.ARRAY_BUFFER, ren.Buffer)
		engo.Gl.VertexAttribPointer(l.inPosition, 2, engo.Gl.FLOAT, false, 8, 0)

		l.lastBuffer = ren.Buffer
	}

	if space.Rotation != 0 {
		sin, cos := math.Sincos(space.Rotation * math.Pi / 180)

		l.modelMatrix[0] = ren.Scale.X * engo.GetGlobalScale().X * cos
		l.modelMatrix[1] = ren.Scale.X * engo.GetGlobalScale().X * sin
		l.modelMatrix[3] = ren.Scale.Y * engo.GetGlobalScale().Y * -sin
		l.modelMatrix[4] = ren.Scale.Y * engo.GetGlobalScale().Y * cos
	} else {
		l.modelMatrix[0] = ren.Scale.X * engo.GetGlobalScale().X
		l.modelMatrix[1] = 0
		l.modelMatrix[3] = 0
		l.modelMatrix[4] = ren.Scale.Y * engo.GetGlobalScale().Y
	}

	if _, ok := ren.Drawable.(StippleLine); !ok {
		l.modelMatrix[6] = space.Position.X * engo.GetGlobalScale().X
		l.modelMatrix[7] = space.Position.Y * engo.GetGlobalScale().Y
	} else {
		l.modelMatrix[6] = 0
		l.modelMatrix[7] = 0
	}

	engo.Gl.UniformMatrix3fv(l.matrixModel, false, l.modelMatrix)
	color := ParseColor(ren.Color).Vec4()
	engo.Gl.Uniform4f(l.uf_Color, color[0], color[1], color[2], color[3])

	switch shape := ren.Drawable.(type) {
	case StippleLine:
		// TODO: modify in the future, engo.Gl not have LineStipple
		gl2.LineStipple(shape.Stipple.Factor, shape.Stipple.Pattern)
		engo.Gl.LineWidth(shape.BorderWidth)
		engo.Gl.DrawArrays(engo.Gl.LINES, 0, len(shape.Points)*2)
	case StippleRect:
		// TODO: modify in the future, engo.Gl not have LineStipple
		gl2.LineStipple(shape.Stipple.Factor, shape.Stipple.Pattern)
		engo.Gl.LineWidth(shape.BorderWidth)
		engo.Gl.DrawArrays(engo.Gl.LINES, 0, 16)
	}
}

func (l *shapeShader) Post() {
	l.lastBuffer = nil

	// Cleanup
	engo.Gl.DisableVertexAttribArray(l.inPosition)

	engo.Gl.BindBuffer(engo.Gl.ARRAY_BUFFER, nil)

	engo.Gl.Disable(engo.Gl.BLEND)
	// TODO: modify in the future, engo.Gl not have LINE_STIPPLE
	engo.Gl.Disable(gl2.LINE_STIPPLE)
}

func (l *shapeShader) SetCamera(c *common.CameraSystem) {
	if l.cameraEnabled {
		l.camera = c
	}
}
