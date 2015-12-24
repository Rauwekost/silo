package engine

import "image"

type ImageOptions struct {
	Width      int
	Height     int
	Fit        bool
	Fill       bool
	Blur       float64
	Sharpen    float64
	Gamma      float64
	Contrast   float64
	Brightness float64
}

type Engine interface {
	Transform(img image.Image, o *ImageOptions) (image.Image, error)
	Resize(img image.Image, w, h int, fit bool) (image.Image, error)
}

func New() Engine {
	return ImageEngine{}
}
