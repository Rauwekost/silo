package service

import (
	"testing"

	"github.com/rauwekost/silo/engine"
)

func getInstance() (*Storage, error) {
	return New(nil, nil)
}

func TestPathForCached(t *testing.T) {
	srv, _ := getInstance()
	sha := "SHA"
	opt := engine.ImageOptions{
		Width:      100,
		Height:     0,
		Fit:        true,
		Fill:       true,
		Blur:       0,
		Sharpen:    0,
		Gamma:      0,
		Contrast:   0,
		Brightness: 0,
	}

	cachedUrl := srv.PathForCached(sha, &opt)
	expect := "S/H/A/SHAb0.00c0.00fltruefttrueg0.00h0s0.00w100"
	if cachedUrl != expect {
		t.Fatalf("cachedUrl invalid expected: %s got: %s", expect, cachedUrl)
	}
}
