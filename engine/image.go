package engine

import (
	"image"

	"github.com/disintegration/imaging"
)

type ImageEngine struct {
}

var defaultResampling = imaging.Lanczos

func (e ImageEngine) Resize(img image.Image, w, h int, fit bool) (image.Image, error) {
	dst := imaging.Resize(img, w, h, defaultResampling)
	return dst, nil
}

func (e ImageEngine) Transform(img image.Image, o *ImageOptions) (image.Image, error) {
	var dst image.Image
	if o.Fit && o.Width > 0 && o.Height > 0 {
		dst = imaging.Fit(img, o.Width, o.Height, defaultResampling)
	} else if o.Fill && o.Width > 0 && o.Height > 0 {
		dst = imaging.Fill(img, o.Width, o.Height, imaging.Center, defaultResampling)
	} else {
		dst = imaging.Resize(img, o.Width, o.Height, defaultResampling)
	}

	if o.Blur > 0 {
		dst = imaging.Blur(dst, o.Blur)
	}

	if o.Sharpen > 0 {
		dst = imaging.Sharpen(dst, o.Sharpen)
	}

	if o.Gamma > 0 {
		dst = imaging.AdjustGamma(dst, o.Gamma)
	}

	if o.Brightness > 0 {
		dst = imaging.AdjustBrightness(dst, o.Brightness)
	}

	if o.Contrast > 0 {
		dst = imaging.AdjustContrast(dst, o.Contrast)
	}

	return dst, nil
}
