package service

import (
	"fmt"
	"image"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"github.com/rauwekost/silo/engine"
	"github.com/rauwekost/silo/metadata"
	"github.com/rauwekost/silo/storage"
	"github.com/rauwekost/silo/store"
)

//Storage service
type Storage struct {
	fs   storage.Storage
	db   *store.Store
	imgp engine.Engine
}

var Extensions = map[string]string{
	"image/jpeg": "jpg",
	"image/png":  "png",
	"image/bmp":  "bmp",
	"image/gif":  "gif",
}

//New return a new service instance
func New(fs storage.Storage, db *store.Store) (*Storage, error) {
	return &Storage{
		fs:   fs,
		db:   db,
		imgp: engine.New(),
	}, nil
}

//Save a file
func (s *Storage) Save(f storage.File, m *metadata.Metadata) error {
	sha := f.Checksum()
	err := s.fs.Save(s.PathForSHA(sha), f)
	if err != nil {
		return err
	}
	err = s.db.Put([]byte(sha), m)
	if err != nil {
		return err
	}
	return nil
}

//Find a file
func (s *Storage) Find(sha string) (storage.File, *metadata.Metadata, error) {
	m, err := s.getFileMeta(sha)
	if err != nil {
		return nil, nil, err
	}
	fmt.Printf("%+v\n", m)

	f, err := s.getFileContent(sha)
	if err != nil {
		return nil, nil, err
	}

	return f, m, nil
}

//Display a image, resampled images are not stored on disk and are generated on the fly.
//we can cache these images in a proxy
func (s *Storage) Display(sha string, opt *engine.ImageOptions) (image.Image, error) {
	//find the file
	f, m, err := s.Find(sha)
	if err != nil {
		return nil, err
	}

	//check if the extension is supported
	_, ok := Extensions[m.ContentType()]
	if !ok {
		return nil, fmt.Errorf("Unsupported content-type")
	}

	//decode data into a image
	img, _, err := image.Decode(f)
	if err != nil {
		return nil, err
	}

	//transform image
	dst, err := s.imgp.Transform(img, opt)
	if err != nil {
		return nil, err
	}

	return dst, nil
}

func (s *Storage) PathForSHA(sha string) string {
	return fmt.Sprintf("%s/%s/%s/%s", string(sha[0]), string(sha[1]), string(sha[2]), sha)
}

//get File Meta data from database
func (s *Storage) getFileMeta(sha string) (*metadata.Metadata, error) {
	var m metadata.Metadata

	err := s.db.Get([]byte(sha), &m)
	if err != nil {
		return nil, err
	}

	return &m, nil
}

func (s *Storage) getFileContent(sha string) (storage.File, error) {
	f, err := s.fs.Open(s.PathForSHA(sha))
	if err != nil {
		return nil, err
	}

	return f, nil
}
