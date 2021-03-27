package engoutil

import (
	"fmt"
	"image/color"
	"math"
)

var _ color.Color = (*Color)(nil)

// Color implementation image.Color Interface
type Color struct {
	r, g, b, a uint32
	raw        uint32
}

func (c Color) String() string {
	return fmt.Sprintf("%X", c.raw)
}

// implementation image.Color interface
func (c Color) RGBA() (uint32, uint32, uint32, uint32) {
	return c.r, c.g, c.b, c.a
}

// (*Color) Set set a new color
func (c *Color) Set(i uint32) {
	r := uint32(byte(i >> 24))
	g := uint32(byte(i >> 16))
	b := uint32(byte(i >> 8))
	a := uint32(byte(i))

	c.r = r | r<<8
	c.g = g | g<<8
	c.b = b | b<<8
	c.a = a | a<<8
	c.raw = i
}

func (c *Color) EqualUint32(u uint32) bool {
	if c == nil {
		return false
	}
	return c.raw == u
}

func (c Color) Vec4() [4]float32 {
	r, g, b, a := c.Unpack()
	return [4]float32{float32(r) / 0xFF, float32(g) / 0xFF, float32(b) / 0xFF, float32(a) / 0xFF}
}

func (c Color) Unpack() (byte, byte, byte, byte) {
	return byte(c.raw >> 24), byte(c.raw >> 16), byte(c.raw >> 8), byte(c.raw)
}

// (Color) ToColor returns a color.Color
func (c Color) ToColor() color.Color {
	return color.RGBA{
		R: byte(c.raw >> 24),
		G: byte(c.raw >> 16),
		B: byte(c.raw >> 8),
		A: byte(c.raw),
	}
}

func (c Color) ToFloat32() float32 {
	colorR, colorG, colorB, colorA := c.RGBA()
	colorR >>= 8
	colorG >>= 8
	colorB >>= 8
	colorA >>= 8

	red := colorR
	green := colorG << 8
	blue := colorB << 16
	alpha := colorA << 24

	return math.Float32frombits((alpha | blue | green | red) & 0xfeffffff)
}

// (Color) Alpha 返回透明度
func (c Color) Alpha() byte {
	return byte(c.raw & 0xFF)
}

// (*Color) SetAlpha 设置透明度
func (c *Color) SetAlpha(i byte) {
	a := uint32(i)
	c.a = a | a<<8
	c.raw &= 0xFFFFFF00 + a
}

// NewColor creating color with uint32
func NewColor(i uint32) *Color {
	r := uint32(byte(i >> 24))
	g := uint32(byte(i >> 16))
	b := uint32(byte(i >> 8))
	a := uint32(byte(i))
	return &Color{
		r:   r | r<<8,
		g:   g | g<<8,
		b:   b | b<<8,
		a:   a | a<<8,
		raw: i,
	}
}

// ParseColor creating color with color.Color
func ParseColor(clr color.Color) *Color {
	r, g, b, a := clr.RGBA()
	return &Color{
		r:   r,
		g:   g,
		b:   b,
		a:   a,
		raw: (r & 0xFF << 24) | (g & 0xFF << 16) | (b & 0xFF << 8) | (a & 0xFF),
	}
}

// ColorEqualUint32 判断 color.Color 和 uint32 是否相等
func ColorEqualUint32(c color.Color, u uint32) bool {
	if c == nil {
		return false
	}
	r, g, b, a := c.RGBA()
	i := (r & 0xFF << 24) | (g & 0xFF << 16) | (b & 0xFF << 8) | (a & 0xFF)
	return i == u
}