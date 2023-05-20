package image

import (
	"bytes"
	"image"
	"image/color"
	_ "image/jpeg"
	_ "image/png"
	"math"
)

type Color struct {
	Red   uint8
	Green uint8
	Blue  uint8
}

func BytesToImage(b []byte) (image.Image, error) {
	img, _, err := image.Decode(bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	return img, nil
}

func CalcAvgColor(img image.Image) (result Color) {
	imgSize := img.Bounds().Size()

	var redSum float64
	var greenSum float64
	var blueSum float64

	for x := 0; x <= imgSize.X; x++ {
		for y := 0; y <= imgSize.Y; y++ {
			pixel := img.At(x, y)
			col := color.RGBAModel.Convert(pixel).(color.RGBA)

			redSum += float64(col.R)
			greenSum += float64(col.G)
			blueSum += float64(col.B)
		}
	}

	imgArea := float64(imgSize.X * imgSize.Y)

	result.Red = uint8(math.Round(redSum / imgArea))
	result.Green = uint8(math.Round(greenSum / imgArea))
	result.Blue = uint8(math.Round(blueSum / imgArea))

	return
}
