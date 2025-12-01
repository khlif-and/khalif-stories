package utils

import (
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"mime/multipart"

	"github.com/generaltso/vibrant"

)

func ExtractDominantColor(file multipart.File) (string, error) {
	img, _, err := image.Decode(file)
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

	r, g, b := bestSwatch.Color.RGB()

	return fmt.Sprintf("#%02x%02x%02x", r, g, b), nil
}