package engine

import (
	"image"
	"reflect"
	"strconv"
)

type ImageOptions struct {
	Width      int     `cache:"wi"`
	Height     int     `cache:"he"`
	Fit        bool    `cache:"ft"`
	Fill       bool    `cache:"fl"`
	Blur       float64 `cache:"bl"`
	Sharpen    float64 `cache:"sh"`
	Gamma      float64 `cache:"ga"`
	Contrast   float64 `cache:"co"`
	Brightness float64 `cache:"br"`
}

type Engine interface {
	Transform(img image.Image, o *ImageOptions) (image.Image, error)
	Resize(img image.Image, w, h int, fit bool) (image.Image, error)
}

func New() Engine {
	return ImageEngine{}
}

func (opt *ImageOptions) TagMap() map[string]string {
	val := reflect.ValueOf(opt).Elem()
	tagMap := make(map[string]string, 0)

	//loop over fields to extract tags and values
	for i := 0; i < val.NumField(); i++ {
		valueField := val.Field(i)
		typeField := val.Type().Field(i)
		var value string
		key := typeField.Tag.Get("cache")
		switch valueField.Kind() {
		case reflect.Int:
			value = strconv.Itoa(int(valueField.Int()))
		case reflect.Float64, reflect.Float32:
			value = strconv.FormatFloat(valueField.Float(), 'f', 1, 64)
		case reflect.Bool:
			if valueField.Bool() == true {
				value = "1"
			} else {
				value = "0"
			}
		default:
			value = valueField.String()
		}

		if key != "" {
			tagMap[key] = value
		}
	}

	return tagMap
}
