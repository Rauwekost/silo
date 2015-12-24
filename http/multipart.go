package http

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/mholt/binding"
	"github.com/rauwekost/silo/metadata"
	"github.com/rauwekost/silo/service"
	"github.com/rauwekost/silo/storage"
)

type MultipartForm struct {
	Data *multipart.FileHeader `json:"data"`
}

func (f *MultipartForm) FieldMap(req *http.Request) binding.FieldMap {
	return binding.FieldMap{
		&f.Data: "data",
	}
}

func (f *MultipartForm) Upload(srv *service.Storage) (*metadata.Metadata, error) {
	var fh io.ReadCloser
	fh, err := f.Data.Open()
	if err != nil {
		return nil, err
	}
	defer fh.Close()

	dataBytes := bytes.Buffer{}
	_, err = dataBytes.ReadFrom(fh)
	if err != nil {
		return nil, err
	}

	//create metadata and save
	file := storage.NewContentFile(dataBytes.Bytes())
	meta := metadata.New(file, f.Data)
	srv.Save(file, meta)

	return meta, nil
}
