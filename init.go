package engoutil

import (
	"bytes"
	"fmt"

	"github.com/EngoEngine/engo"

	"engo/compon/assets"
)

const Font04b08 = "04b08.ttf"
const FontDroidSans = "DroidSansFallbackFull.ttf"

func init() {
	var err error
	if err = engo.Files.LoadReaderData(Font04b08, bytes.NewReader(assets.Font04b08)); err != nil {
		panic(fmt.Sprintf("unable to load %q! Error was: ", Font04b08) + err.Error())
	}
	if err = engo.Files.LoadReaderData(FontDroidSans, bytes.NewReader(assets.FontDroidSans)); err != nil {
		panic(fmt.Sprintf("unable to load %q! Error was: ", FontDroidSans) + err.Error())
	}
}
