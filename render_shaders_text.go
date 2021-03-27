package engoutil

import (
	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
	"github.com/EngoEngine/engo/math"
	"github.com/EngoEngine/gl"
)

type textShader struct {
	program *gl.Program

	indicesRectangles    []uint16
	indicesRectanglesVBO *gl.Buffer

	inPosition  int
	inTexCoords int

	matrixProjection *gl.UniformLocation
	matrixView       *gl.UniformLocation
	matrixModel      *gl.UniformLocation
	uf_Color         *gl.UniformLocation
	uf_Target        *gl.UniformLocation
	uf_TexSize       *gl.UniformLocation

	projectionMatrix []float32
	viewMatrix       []float32
	modelMatrix      []float32

	camera        *common.CameraSystem
	cameraEnabled bool

	lastBuffer  *gl.Buffer
	lastTexture *gl.Texture
}

func (l *textShader) Setup(*ecs.World) error {
	var err error
	l.program, err = common.LoadShader(`
attribute vec2 in_Position;
attribute vec2 in_TexCoords;

uniform mat3 matrixProjection;
uniform mat3 matrixView;
uniform mat3 matrixModel;

varying vec2 var_TexCoords;

void main() {
  var_TexCoords = in_TexCoords;

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

varying vec2 var_TexCoords;

uniform sampler2D uf_Texture;
uniform vec4 uf_Color;
uniform int uf_Target;
uniform vec2 uf_TexSize;

void main (void) {
  if (uf_Target == 1) {
    gl_FragColor = uf_Color;
  } else if (uf_Target == 2) {
	float alpha = 4.0*texture2D(uf_Texture, var_TexCoords).a;
	alpha -= texture2D(uf_Texture, var_TexCoords + vec2(uf_TexSize.x, 0)).a;
	alpha -= texture2D(uf_Texture, var_TexCoords + vec2(-uf_TexSize.x, 0)).a;
	alpha -= texture2D(uf_Texture, var_TexCoords + vec2(0, uf_TexSize.y)).a;
    alpha -= texture2D(uf_Texture, var_TexCoords + vec2(0, -uf_TexSize.y)).a;
    gl_FragColor = vec4(uf_Color.xyz, alpha);
  } else {
	gl_FragColor = uf_Color * texture2D(uf_Texture, var_TexCoords);
  }
}`)

	if err != nil {
		return err
	}

	// Create and populate indices buffer
	l.indicesRectangles = make([]uint16, 6*bufferSize)
	for i, j := 0, 0; i < bufferSize*6; i, j = i+6, j+4 {
		l.indicesRectangles[i+0] = uint16(j + 0)
		l.indicesRectangles[i+1] = uint16(j + 1)
		l.indicesRectangles[i+2] = uint16(j + 2)
		l.indicesRectangles[i+3] = uint16(j + 0)
		l.indicesRectangles[i+4] = uint16(j + 2)
		l.indicesRectangles[i+5] = uint16(j + 3)
	}
	l.indicesRectanglesVBO = engo.Gl.CreateBuffer()
	engo.Gl.BindBuffer(engo.Gl.ELEMENT_ARRAY_BUFFER, l.indicesRectanglesVBO)
	engo.Gl.BufferData(engo.Gl.ELEMENT_ARRAY_BUFFER, l.indicesRectangles, engo.Gl.STATIC_DRAW)

	// Define things that should be read from the texture buffer
	l.inPosition = engo.Gl.GetAttribLocation(l.program, "in_Position")
	l.inTexCoords = engo.Gl.GetAttribLocation(l.program, "in_TexCoords")

	// Define things that should be set per draw
	l.matrixProjection = engo.Gl.GetUniformLocation(l.program, "matrixProjection")
	l.matrixView = engo.Gl.GetUniformLocation(l.program, "matrixView")
	l.matrixModel = engo.Gl.GetUniformLocation(l.program, "matrixModel")
	l.uf_Color = engo.Gl.GetUniformLocation(l.program, "uf_Color")
	l.uf_Target = engo.Gl.GetUniformLocation(l.program, "uf_Target")
	l.uf_TexSize = engo.Gl.GetUniformLocation(l.program, "uf_TexSize")

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

func (l *textShader) Pre() {
	engo.Gl.Enable(engo.Gl.BLEND)
	engo.Gl.BlendFunc(engo.Gl.SRC_ALPHA, engo.Gl.ONE_MINUS_SRC_ALPHA)

	// Bind shader and buffer, enable attributes
	engo.Gl.UseProgram(l.program)
	engo.Gl.BindBuffer(engo.Gl.ELEMENT_ARRAY_BUFFER, l.indicesRectanglesVBO)
	engo.Gl.EnableVertexAttribArray(l.inPosition)
	engo.Gl.EnableVertexAttribArray(l.inTexCoords)

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

func (l *textShader) updateBuffer(ren *common.RenderComponent, space *common.SpaceComponent) {
	txt, ok := ren.Drawable.(*Text)
	if !ok {
		unsupportedType(ren.Drawable)
		return
	}

	if len(ren.BufferContent) < 16*txt.Length()+16 {
		ren.BufferContent = make([]float32, 16*txt.Length()+16)
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

func (l *textShader) generateBufferContent(ren *common.RenderComponent, space *common.SpaceComponent, buffer []float32) (changed bool) {
	txt, ok := ren.Drawable.(*Text)
	if !ok {
		unsupportedType(ren.Drawable)
		return
	}

	var (
		size               [2]float32
		atlas              = txt.Font.updateFontAtlas(txt.Text)
		currentX, currentY float32
		count              = 1
		index              int
		w, h, x, y         float32
		offsetX, offsetY   float32
	)

	if txt.changed() {
		for i := 0; i < len(buffer); i++ {
			buffer[i] = 0
		}
		for _, char := range []rune(txt.Text) {
			// skip invisible characters
			if _, ok := atlas.Width[char]; !ok {
				continue
			}
			if char == '\n' {
				currentX = 0
				currentY += atlas.LineHeight + txt.LineSpacing
				size[1] += atlas.LineHeight + txt.LineSpacing
				continue
			}
			index = count * 16
			w, h = atlas.Width[char], atlas.Height[char]
			x, y = atlas.XLocation[char], atlas.YLocation[char]
			offsetX = atlas.LeftSide[char]
			offsetY = atlas.OffsetY[char]

			// These five are at 0, 0:
			setBufferValue(buffer, 0+index, currentX+offsetX, &changed)
			setBufferValue(buffer, 1+index, currentY+offsetY, &changed)
			setBufferValue(buffer, 2+index, x/atlas.TotalWidth, &changed)
			setBufferValue(buffer, 3+index, y/atlas.TotalHeight, &changed)

			// These five are at 1, 0:
			setBufferValue(buffer, 4+index, currentX+offsetX+w, &changed)
			setBufferValue(buffer, 5+index, currentY+offsetY, &changed)
			setBufferValue(buffer, 6+index, (x+w)/atlas.TotalWidth, &changed)
			setBufferValue(buffer, 7+index, y/atlas.TotalHeight, &changed)

			// These five are at 1, 1:
			setBufferValue(buffer, 8+index, currentX+offsetX+w, &changed)
			setBufferValue(buffer, 9+index, currentY+h+offsetY, &changed)
			setBufferValue(buffer, 10+index, (x+w)/atlas.TotalWidth, &changed)
			setBufferValue(buffer, 11+index, (y+h)/atlas.TotalHeight, &changed)

			// These five are at 0, 1:
			setBufferValue(buffer, 12+index, currentX+offsetX, &changed)
			setBufferValue(buffer, 13+index, currentY+h+offsetY, &changed)
			setBufferValue(buffer, 14+index, x/atlas.TotalWidth, &changed)
			setBufferValue(buffer, 15+index, (y+h)/atlas.TotalHeight, &changed)

			currentX += offsetX + w + atlas.RightSide[char] + txt.LetterSpacing
			if size[0] < currentX {
				size[0] = currentX
			}
			count++
		}
		size[1] = currentY + atlas.LineHeight
		txt.size = size
		txt.buffered.text = txt.Text
		txt.buffered.lineSpacing = txt.LineSpacing
		txt.buffered.letterSpacing = txt.LetterSpacing
	} else {
		size = txt.size
	}

	w, h = size[0], size[1]
	if txt.width != 0 && txt.height != 0 {
		w, h = txt.width*txt.Font.scale, txt.height*txt.Font.scale
	}
	x, y = 0, 0
	x = (size[0] - w) / 2
	y = (size[1] - h) / 2
	y -= txt.Padding.Top
	h += txt.Padding.Top
	w += txt.Padding.Right
	h += txt.Padding.Bottom
	x -= txt.Padding.Left
	w += txt.Padding.Left

	// background rectangle
	setBufferValue(buffer, 0, x, &changed)
	setBufferValue(buffer, 1, y, &changed)
	// setBufferValue(buffer, 2, 0, &changed)
	// setBufferValue(buffer, 3, 0, &changed)

	setBufferValue(buffer, 4, x+w, &changed)
	setBufferValue(buffer, 5, y, &changed)
	// setBufferValue(buffer, 6, 0, &changed)
	// setBufferValue(buffer, 7, 0, &changed)

	setBufferValue(buffer, 8, x+w, &changed)
	setBufferValue(buffer, 9, y+h, &changed)
	// setBufferValue(buffer, 10, 0, &changed)
	// setBufferValue(buffer, 11, 0, &changed)

	setBufferValue(buffer, 12, x, &changed)
	setBufferValue(buffer, 13, y+h, &changed)
	// setBufferValue(buffer, 14, 0, &changed)
	// setBufferValue(buffer, 15, 0, &changed)

	// for MouseComponent working
	space.Position.X = txt.Position.X + x
	space.Position.Y = txt.Position.Y + y
	space.Width = w * ren.Scale.X
	space.Height = h * ren.Scale.X

	return
}

func (l *textShader) Draw(ren *common.RenderComponent, space *common.SpaceComponent) {
	txt, ok := ren.Drawable.(*Text)
	if !ok {
		unsupportedType(ren.Drawable)
	}

	if l.lastBuffer != ren.Buffer || ren.Buffer == nil {
		l.updateBuffer(ren, space)

		engo.Gl.BindBuffer(engo.Gl.ARRAY_BUFFER, ren.Buffer)
		engo.Gl.VertexAttribPointer(l.inPosition, 2, engo.Gl.FLOAT, false, 16, 0)
		engo.Gl.VertexAttribPointer(l.inTexCoords, 2, engo.Gl.FLOAT, false, 16, 8)
		l.lastBuffer = ren.Buffer
	}

	atlas := txt.Font.updateFontAtlas(txt.Text)

	if atlas.Texture != l.lastTexture {
		engo.Gl.BindTexture(engo.Gl.TEXTURE_2D, atlas.Texture)
		l.lastTexture = atlas.Texture
	}

	engo.Gl.TexParameteri(engo.Gl.TEXTURE_2D, engo.Gl.TEXTURE_WRAP_S, engo.Gl.CLAMP_TO_EDGE)
	engo.Gl.TexParameteri(engo.Gl.TEXTURE_2D, engo.Gl.TEXTURE_WRAP_T, engo.Gl.CLAMP_TO_EDGE)

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

	l.modelMatrix[6] = txt.Position.X * engo.GetGlobalScale().X
	l.modelMatrix[7] = txt.Position.Y * engo.GetGlobalScale().Y

	engo.Gl.UniformMatrix3fv(l.matrixModel, false, l.modelMatrix)

	// draw background
	if _, _, _, alpha := ren.Color.RGBA(); alpha > 0 {
		bg := ParseColor(ren.Color).Vec4()
		engo.Gl.Uniform1i(l.uf_Target, 1)
		engo.Gl.Uniform4f(l.uf_Color, bg[0], bg[1], bg[2], bg[3])
		if txt.BgStyle == BG_FILL_WRAP {
			engo.Gl.DrawElements(engo.Gl.TRIANGLES, 6*txt.Length(), engo.Gl.UNSIGNED_SHORT, 12)
		} else {
			engo.Gl.DrawElements(engo.Gl.TRIANGLES, 6, engo.Gl.UNSIGNED_SHORT, 0)
		}
	}

	// TODO: draw text outline
	// engo.Gl.Uniform1i(l.uf_Target, 2)
	// engo.Gl.Uniform2f(l.uf_TexSize, 1/txt.size[0], 1/txt.size[1])
	// engo.Gl.DrawElements(engo.Gl.TRIANGLES, 6*txt.Length(), engo.Gl.UNSIGNED_SHORT, 12)

	// draw text
	fg := txt.Color.Vec4()
	engo.Gl.Uniform1i(l.uf_Target, 0)
	engo.Gl.Uniform4f(l.uf_Color, fg[0], fg[1], fg[2], fg[3])
	engo.Gl.DrawElements(engo.Gl.TRIANGLES, 6*txt.Length(), engo.Gl.UNSIGNED_SHORT, 12)
}

func (l *textShader) Post() {
	l.lastBuffer = nil
	l.lastTexture = nil

	// Cleanup
	engo.Gl.DisableVertexAttribArray(l.inPosition)
	engo.Gl.DisableVertexAttribArray(l.inTexCoords)

	engo.Gl.BindTexture(engo.Gl.TEXTURE_2D, nil)
	engo.Gl.BindBuffer(engo.Gl.ARRAY_BUFFER, nil)
	engo.Gl.BindBuffer(engo.Gl.ELEMENT_ARRAY_BUFFER, nil)

	engo.Gl.Disable(engo.Gl.BLEND)
}

func (l *textShader) SetCamera(c *common.CameraSystem) {
	if l.cameraEnabled {
		l.camera = c
	}
}
