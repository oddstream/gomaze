package maze

import (
	_ "embed" // go:embed only allowed in Go files that import "embed"
	"log"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
)

//go:embed assets/Acme-Regular.ttf
var acmeFontBytes []byte

// AcmeFonts contains references to small, normal, large and huge Acme fonts
type AcmeFonts struct {
	small  font.Face
	normal font.Face
	large  font.Face
	huge   font.Face
}

// NewAcmeFonts loads some fonts and returns a pointer to an object referencing them
func NewAcmeFonts() *AcmeFonts {

	tt, err := truetype.Parse(acmeFontBytes)
	if err != nil {
		log.Fatal(err)
	}

	af := &AcmeFonts{}

	af.small = truetype.NewFace(tt, &truetype.Options{
		Size:    16,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	af.normal = truetype.NewFace(tt, &truetype.Options{
		Size:    32,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	af.large = truetype.NewFace(tt, &truetype.Options{
		Size:    48,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	af.huge = truetype.NewFace(tt, &truetype.Options{
		Size:    128,
		DPI:     72,
		Hinting: font.HintingFull,
	})

	return af
}
