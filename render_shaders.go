package engoutil

import (
	"image"
	"log"

	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
	"github.com/EngoEngine/gl"
)

var (
	TextShader     = &textShader{cameraEnabled: true}
	TextHUDShader  = &textShader{}
	ShapeShader    = &shapeShader{cameraEnabled: true}
	ShapeHUDShader = &shapeShader{}

	atlasCache = make(map[Font]*FontAtlas)

	bufferSize = 10000

	shaders     = []common.Shader{TextShader, TextHUDShader, ShapeShader, ShapeHUDShader}
	shadersInit bool
)

func InitShaders() {
	if shadersInit {
		return
	}
	shadersInit = true
	for _, s := range shaders {
		common.AddShader(s)
	}
}

var _ common.Drawable = (*StippleLine)(nil)
var _ common.Drawable = (*StippleRect)(nil)

type Stipple struct {
	Factor  int32
	Pattern uint16
}

// StippleLine
type StippleLine struct {
	BorderWidth float32
	Points      []engo.Point
	Stipple     Stipple
}

func (StippleLine) Texture() *gl.Texture                       { return nil }
func (StippleLine) Width() float32                             { return 0 }
func (StippleLine) Height() float32                            { return 0 }
func (StippleLine) View() (float32, float32, float32, float32) { return 0, 0, 1, 1 }
func (StippleLine) Close()                                     {}

// StippleRect
type StippleRect struct {
	BorderWidth float32
	Stipple     Stipple
}

func (StippleRect) Texture() *gl.Texture                       { return nil }
func (StippleRect) Width() float32                             { return 0 }
func (StippleRect) Height() float32                            { return 0 }
func (StippleRect) View() (float32, float32, float32, float32) { return 0, 0, 1, 1 }
func (StippleRect) Close()                                     {}

func setBufferValue(buffer []float32, index int, value float32, changed *bool) {
	if buffer[index] != value {
		buffer[index] = value
		*changed = true
	}
}

func resizeNRGBAHeight(nrgba *image.NRGBA, h float32) {
	height := int(h)
	if height < nrgba.Rect.Max.Y {
		return
	}
	var b []uint8
	b, nrgba.Pix = nrgba.Pix, make([]uint8, len(nrgba.Pix)+height*nrgba.Rect.Max.X*4)
	copy(nrgba.Pix, b)
	nrgba.Rect.Max.Y = height
}

func unsupportedType(v interface{}) {
	warning("type %T not supported", v)
}

func warning(format string, a ...interface{}) {
	log.Printf("[WARNING] "+format, a...)
}
