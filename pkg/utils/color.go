package utils

import (
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"

	"github.com/generaltso/vibrant"

)

// ExtractDominantColor sekarang menerima io.Reader (lebih generic)
func ExtractDominantColor(r io.Reader) (string, error) {
	img, _, err := image.Decode(r)
	if err != nil {
		return "", err
	}

	palette, err := vibrant.NewPaletteFromImage(img)
	if err != nil {
		return "", err
	}

	swatches := palette.ExtractAwesome()
	var bestSwatch *vibrant.Swatch

	if sw, ok := swatches["Vibrant"]; ok {
		bestSwatch = sw
	} else if sw, ok := swatches["LightVibrant"]; ok {
		bestSwatch = sw
	} else if sw, ok := swatches["DarkVibrant"]; ok {
		bestSwatch = sw
	} else if sw, ok := swatches["Muted"]; ok {
		bestSwatch = sw
	} else {
		for _, sw := range swatches {
			bestSwatch = sw
			break
		}
	}

	if bestSwatch == nil {
		return "#000000", nil
	}

	rVal, gVal, bVal := bestSwatch.Color.RGB()
	return fmt.Sprintf("#%02x%02x%02x", rVal, gVal, bVal), nil
}