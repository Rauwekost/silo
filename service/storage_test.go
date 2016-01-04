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
	expect := "S/H/A/SHA-bl0.0-br0.0-co0.0-fl1-ft1-ga0.0-he0-sh0.0-wi100"
	if cachedUrl != expect {
		t.Fatalf("cachedUrl invalid expected: %s got: %s", expect, cachedUrl)
	}
}
