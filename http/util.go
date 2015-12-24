package http

import (
	"bytes"
	"encoding/json"
	"image"
	"image/jpeg"
	"net/http"
	"strconv"

	log "github.com/Sirupsen/logrus"
)

//HttpHandler wrapper that expects a error in return
type HttpHandler func(http.ResponseWriter, *http.Request) *Error

//HttpHandler's ServeHTTP method
func (fn HttpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if e := fn(w, r); e != nil { // e is *Error, not os.Error.
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(e.Status)
		json.NewEncoder(w).Encode(e)
	}
}

//JSON helper to write json to Response writer
func JSON(w http.ResponseWriter, status int, obj interface{}) {
	j, err := json.Marshal(obj)
	if err != nil {
		log.Errorf("Unable to marshal %#v to JSON: %v", obj, err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(j)
}

func writeImage(w http.ResponseWriter, img *image.Image) {
	buffer := new(bytes.Buffer)
	if err := jpeg.Encode(buffer, *img, nil); err != nil {
		log.Println("unable to encode image.")
	}

	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("Content-Length", strconv.Itoa(len(buffer.Bytes())))
	if _, err := w.Write(buffer.Bytes()); err != nil {
		log.Warn("unable to write image.")
	}
}
