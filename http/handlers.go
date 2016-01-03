package http

import (
	"io"
	"net/http"
	"strconv"

	"github.com/rauwekost/silo/Godeps/_workspace/src/github.com/mholt/binding"
	"github.com/rauwekost/silo/engine"
)

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) *Error {
	JSON(w, 200, map[string]interface{}{"status": "ok"})
	return nil
}

func (s *Server) versionHandler(w http.ResponseWriter, r *http.Request) *Error {
	JSON(w, 200, map[string]interface{}{"version": "ok"})
	return nil
}

func (s *Server) uploadHandler(w http.ResponseWriter, r *http.Request) *Error {
	multipartForm := new(MultipartForm)
	errs := binding.Bind(r, multipartForm)
	if errs.Handle(w) {
		return nil
	}

	f, err := multipartForm.Upload(s.StorageService)
	if err != nil {
		return ErrInternalServer("failed to upload")
	}

	JSON(w, 200, f)
	return nil
}

func (s *Server) displayHandler(w http.ResponseWriter, r *http.Request) *Error {
	id := r.URL.Query().Get("id")
	if id == "" {
		return ErrInvalidRequest("no id provided")
	}
	options, err := s.ImageOptionsFromRequest(r)
	if err != nil {
		return ErrInvalidRequest(err.Error())
	}

	img, err := s.StorageService.Display(id, options)
	if err != nil {
		return ErrInvalidRequest(err.Error())
	}

	writeImage(w, &img)
	return nil
}

func (s *Server) downloadHandler(w http.ResponseWriter, r *http.Request) *Error {
	id := r.URL.Query().Get("id")
	if id == "" {
		return ErrInvalidRequest("no id provided")
	}
	f, m, err := s.StorageService.Find(id)
	if err != nil {
		return ErrNotFound(err.Error())
	}

	w.Header().Set("Content-Disposition", "attachment; filename="+m.Name)
	w.Header().Set("Content-Type", m.ContentType())
	w.Header().Set("Content-Lenght", strconv.Itoa(m.Size))
	io.Copy(w, f)

	return nil
}

func (s *Server) ImageOptionsFromRequest(r *http.Request) (*engine.ImageOptions, error) {
	w, _ := strconv.Atoi(r.URL.Query().Get("w"))
	h, _ := strconv.Atoi(r.URL.Query().Get("h"))
	fit, _ := strconv.ParseBool(r.URL.Query().Get("fit"))
	fill, _ := strconv.ParseBool(r.URL.Query().Get("fill"))
	blur, _ := strconv.ParseFloat(r.URL.Query().Get("blur"), 64)
	sharpen, _ := strconv.ParseFloat(r.URL.Query().Get("sharpen"), 64)
	gamma, _ := strconv.ParseFloat(r.URL.Query().Get("gamma"), 64)
	brightness, _ := strconv.ParseFloat(r.URL.Query().Get("brightness"), 64)
	contrast, _ := strconv.ParseFloat(r.URL.Query().Get("contrast"), 64)

	return &engine.ImageOptions{
		Width:      w,
		Height:     h,
		Fit:        fit,
		Fill:       fill,
		Blur:       blur,
		Sharpen:    sharpen,
		Gamma:      gamma,
		Brightness: brightness,
		Contrast:   contrast,
	}, nil
}
