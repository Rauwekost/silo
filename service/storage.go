package service

import (
	"bytes"
	"fmt"
	"image"
	"sort"

	"golang.org/x/image/bmp"

	"image/gif"
	"image/jpeg"
	"image/png"

	log "github.com/Sirupsen/logrus"
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

	//check if the image can be found in cache
	log.Infof("cache-path: %s", s.PathForCached(sha, opt))
	cached, err := s.fs.Open(s.PathForCached(sha, opt))
	if err == nil {
		log.Debug("found cached version")

		//decode data into a image
		img, _, err := image.Decode(cached)
		if err != nil {
			return nil, err
		}

		//return cached image
		return img, nil
	}

	//decode  data into a image
	img, _, err := image.Decode(f)
	if err != nil {
		return nil, err
	}

	//transform image
	dst, err := s.imgp.Transform(img, opt)
	if err != nil {
		return nil, err
	}

	//save the cached image
	if err := s.saveImage(dst, s.PathForCached(sha, opt), m); err != nil {
		return nil, err
	}

	return dst, nil
}

func (s *Storage) saveImage(img image.Image, path string, m *metadata.Metadata) error {
	buf := new(bytes.Buffer)
	if err := encodeImage(buf, img, m.ContentType()); err != nil {
		return err
	}

	cf := storage.NewContentFile(buf.Bytes())
	if err := s.fs.Save(path, cf); err != nil {
		return err
	}

	return nil
}

func encodeImage(buf *bytes.Buffer, img image.Image, contentType string) error {
	switch contentType {
	case "image/jpeg":
		return jpeg.Encode(buf, img, nil)
	case "image/png":
		return png.Encode(buf, img)
	case "image/gif":
		return gif.Encode(buf, img, nil)
	case "image/bmp":
		return bmp.Encode(buf, img)
	default:
		return fmt.Errorf("unknown image format: %s", contentType)
	}
}

func (s *Storage) PathForSHA(sha string) string {
	return fmt.Sprintf("%s/%s/%s/%s", string(sha[0]), string(sha[1]), string(sha[2]), sha)
}

func (s *Storage) PathForCached(sha string, opt *engine.ImageOptions) string {
	path := s.PathForSHA(sha)
	tagMap := opt.TagMap()
	for _, k := range sortKeys(tagMap) {
		path = path + "-" + string(k) + tagMap[k]
	}
	return path
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

//sortKeys is a helper function that sorts keys alphabetically
func sortKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k, _ := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
