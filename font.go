package engoutil

import (
	"fmt"
	"image"
	"image/color"
	"io/ioutil"

	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
	"github.com/EngoEngine/gl"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

// font texture default width
const fontTextureDefWidth = 2048

type Padding struct {
	Top, Right, Bottom, Left float32
}

const (
	BG_FILL_FULL uint8 = iota
	BG_FILL_WRAP
)

type Font struct {
	URL   string
	Size  float64
	TTF   *truetype.Font
	face  font.Face
	scale float32
}

// Create is for loading fonts from the disk, given a location
func (f *Font) Create() error {
	// Read and parse the font
	ttfBytes, err := ioutil.ReadFile(f.URL)
	if err != nil {
		return err
	}

	ttf, err := freetype.ParseFont(ttfBytes)
	if err != nil {
		return err
	}
	f.TTF = ttf
	f.face = truetype.NewFace(f.TTF, &truetype.Options{
		Size:    f.Size,
		Hinting: font.HintingFull,
		DPI:     72 * float64(f.scale),
	})

	return nil
}

// CreatePreloaded is for loading fonts which have already been defined (and loaded) within Preload
func (f *Font) CreatePreloaded() error {
	fontres, err := engo.Files.Resource(f.URL)
	if err != nil {
		return err
	}

	fnt, ok := fontres.(common.FontResource)
	if !ok {
		return fmt.Errorf("preloaded font is not of type `*truetype.Font`: %s", f.URL)
	}

	f.TTF = fnt.Font
	f.face = truetype.NewFace(f.TTF, &truetype.Options{
		Size:    f.Size,
		Hinting: font.HintingFull,
		DPI:     72 * float64(f.scale),
	})
	return nil
}

func (f *Font) updateFontAtlas(text string) (atlas *FontAtlas) {

	var initialized bool
	if _, ok := atlasCache[*f]; ok {
		atlas = atlasCache[*f]
		initialized = true
	} else {
		atlas = &FontAtlas{
			XLocation: make(map[rune]float32),
			YLocation: make(map[rune]float32),
			Width:     make(map[rune]float32),
			Height:    make(map[rune]float32),
			// Use fixed width, when the capacity is insufficient, only expand the height
			TotalWidth: fontTextureDefWidth,
			Ascent:     float32(f.face.Metrics().Ascent.Ceil()),
			LeftSide:   make(map[rune]float32),
			RightSide:  make(map[rune]float32),
			OffsetY:    make(map[rune]float32),
			LineHeight: float32((f.face.Metrics().Ascent + f.face.Metrics().Descent).Ceil()),
		}
		atlasCache[*f] = atlas
	}

	if !initialized {
		// Generate basic character 0x20(Space) ~ 0x7E(~)
		basic := make([]rune, 95)
		for i := range basic {
			basic[i] = rune(i) + 32
		}
		text = string(basic) + text
	}

	var (
		// characters that need to be updated
		fresh         []rune
		w, h, advance float32
		prev          rune
		ascent        = make(map[rune]float32)
	)

	for _, char := range text {
		if _, has := atlas.Width[char]; has {
			continue
		}
		bounds, adv, ok := f.face.GlyphBounds(char)
		if !ok {
			continue
		}

		fresh = append(fresh, char)

		w = float32((bounds.Max.X - bounds.Min.X).Ceil())
		h = float32((bounds.Max.Y - bounds.Min.Y).Ceil())
		advance = float32(adv.Ceil())
		ascent[char] = float32(-bounds.Min.Y.Ceil())

		atlas.Width[char] = w
		atlas.Height[char] = h
		atlas.LeftSide[char] = float32(bounds.Min.X.Ceil())
		atlas.RightSide[char] = advance - float32(bounds.Max.X.Ceil())
		atlas.OffsetY[char] = atlas.Ascent - ascent[char]

		if prev > 0 {
			atlas.CurrentX += float32(f.face.Kern(prev, char).Ceil())
		}

		// partially overlapped characters
		if w > advance {
			if atlas.LeftSide[char] < 0 {
				advance -= atlas.LeftSide[char]
			} else if atlas.RightSide[char] < 0 {
				advance -= atlas.RightSide[char]
			}
		}

		// position correction of overlapped characters
		if atlas.LeftSide[char] < 0 {
			atlas.CurrentX -= atlas.LeftSide[char]
		}

		if atlas.CurrentX+advance > atlas.TotalWidth {
			atlas.CurrentX = 0
			atlas.CurrentY += atlas.LineHeight
			prev = 0
		}

		atlas.XLocation[char] = atlas.CurrentX
		atlas.YLocation[char] = atlas.CurrentY
		atlas.CurrentX += advance
		prev = char
	}
	if len(fresh) > 0 {
		atlas.TotalHeight = atlas.CurrentY + atlas.LineHeight
	}

	var actual *image.NRGBA

	// When the image is nil or needs to be updated
	if atlas.Image == nil {
		actual = image.NewNRGBA(image.Rect(0, 0, int(atlas.TotalWidth), int(atlas.TotalHeight)))
	} else if len(fresh) > 0 {
		actual = atlas.Image
		if int(atlas.TotalHeight) > atlas.Image.Rect.Max.Y {
			resizeNRGBAHeight(actual, atlas.TotalHeight)
		}
	}

	if len(fresh) > 0 {
		d := &font.Drawer{}
		d.Src = image.NewUniform(color.White)
		d.Face = f.face
		d.Dst = actual
		for _, char := range fresh {
			// draw on baseline
			d.Dot = fixed.P(int(atlas.XLocation[char]), int(atlas.YLocation[char]+atlas.Ascent))
			d.DrawString(string(char))
			// position correction
			atlas.XLocation[char] += atlas.LeftSide[char]
			atlas.YLocation[char] += atlas.OffsetY[char]
		}
	}

	if actual != nil {
		atlas.Image = actual
		atlas.Texture = common.NewTextureSingle(common.NewImageObject(actual)).Texture()
	}
	return
}

// A FontAtlas is a representation of some of the Font characters, as an image
type FontAtlas struct {
	Image   *image.NRGBA
	Texture *gl.Texture

	// XLocation contains the X-coordinate of the starting position of all characters
	XLocation map[rune]float32
	// YLocation contains the Y-coordinate of the starting position of all characters
	YLocation map[rune]float32
	// Width contains the width in pixels of all the characters, including the spacing between characters
	Width map[rune]float32
	// Height contains the height in pixels of all the characters
	Height map[rune]float32
	// TotalWidth is the total amount of pixels the `FontAtlas` is wide; useful for determining the `Viewport`,
	// which is relative to this value.
	TotalWidth float32
	// TotalHeight is the total amount of pixels the `FontAtlas` is high; useful for determining the `Viewport`,
	// which is relative to this value.
	TotalHeight float32
	// The position of the last character
	CurrentX, CurrentY float32

	// LineHeight is Ascent+Descent
	LineHeight float32
	Ascent     float32
	// left-side and right-side bearings
	LeftSide, RightSide map[rune]float32
	// OffsetY Ascent + bounds.Min.Y
	OffsetY map[rune]float32
}

// Text represents a string drawn onto the screen, as used by the `TextShader`.
type Text struct {
	// Font is the reference to the font you're using to render this. This includes the font size.
	Font *Font
	// Text is the actual text you want to draw. This may include newlines (\n).
	Text string
	// LineSpacing is the amount of additional spacing there is between the lines (when `Text` consists of multiple lines).
	LineSpacing float32
	// LetterSpacing is the amount of additional spacing there is between the characters.
	LetterSpacing float32

	// The shader uses the position here instead of the SpaceComponent.Position.
	// Because Padding changes size and position of SpaceComponent.
	Position engo.Point
	// font color
	Color *Color
	// BG fill style, BG_FILL_FULL or BG_FILL_WRAP
	BgStyle uint8
	// Only when the BgStyle is BG_TYPE_FULL.
	// This changes the size of the entity (common.SpaceComponent).
	// In order to make common.MouseComponent working.
	Padding Padding

	// The text rendered last time is used to reduce buffer update.
	buffered struct {
		text          string
		lineSpacing   float32
		letterSpacing float32
	}
	// The size calculated from the last rendering
	size [2]float32
	// The original size of the font. Used to change the size of the background,
	// this does not include padding.
	width, height float32
}

// Texture returns nil because the Text is generated from a FontAtlas. This implements the common.Drawable interface.
func (t Text) Texture() *gl.Texture { return nil }

// Width returns the width of the Text generated from a FontAtlas. This implements the common.Drawable interface.
func (t Text) Width() (width float32) {
	atlas := t.Font.updateFontAtlas(t.Text)

	var currentX float32

	for _, char := range []rune(t.Text) {
		// skip invisible characters
		if _, ok := atlas.Width[char]; !ok {
			continue
		}
		if char == '\n' {
			if width < currentX {
				width = currentX
			}
			currentX = 0
			continue
		}
		currentX += atlas.LeftSide[char] + atlas.Width[char] + atlas.RightSide[char] + t.LetterSpacing
		if width < currentX {
			width = currentX
		}
	}
	return
}

// Height returns the height the Text generated from a FontAtlas. This implements the common.Drawable interface.
func (t Text) Height() (height float32) {
	atlas := t.Font.updateFontAtlas(t.Text)

	var currentY float32

	for _, char := range []rune(t.Text) {
		// skip invisible characters
		if _, ok := atlas.Height[char]; !ok {
			continue
		}
		if char == '\n' {
			currentY += atlas.LineHeight + t.LineSpacing
			continue
		}
	}
	height = currentY + atlas.LineHeight
	return
}

func (t Text) View() (float32, float32, float32, float32) { return 0, 0, 1, 1 }

func (t Text) Close() {}

func (t Text) Length() int { return len([]rune(t.Text)) }

func (t Text) changed() bool {
	return t.buffered.text != t.Text || t.buffered.lineSpacing != t.LineSpacing || t.buffered.letterSpacing != t.LetterSpacing
}
