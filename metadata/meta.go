package metadata

import (
	"mime/multipart"
	"time"

	log "github.com/rauwekost/silo/Godeps/_workspace/src/github.com/Sirupsen/logrus"
	"github.com/rauwekost/silo/storage"
)

var HeaderKeys = []string{
	"Age",
	"Content-Type",
	"Last-Modified",
	"Date",
	"Etag",
}

type Metadata struct {
	Name       string            `json:"name"`
	CreatedAt  time.Time         `json:"created_at"`
	ModifiedAt time.Time         `json:"modified_at"`
	Size       int               `json:"size"`
	Headers    map[string]string `json:"headers"`
	Checksum   string            `json:"checksum"`
}

func New(f storage.File, t *multipart.FileHeader) *Metadata {
	now := time.Now()
	m := Metadata{
		Name:       t.Filename,
		CreatedAt:  now,
		ModifiedAt: now,
		Size:       int(f.Size()),
		Checksum:   f.Checksum(),
	}

	headers := make(map[string]string, 0)
	for _, key := range HeaderKeys {
		if value, ok := t.Header[key]; ok && len(value) > 0 {
			headers[key] = value[0]
		}
	}
	m.Headers = headers
	return &m
}

func (m *Metadata) ContentType() string {
	typ, ok := m.Headers["Content-Type"]
	if !ok {
		log.Errorf("could not find content-type for: %s created-at: %s", m.Name, m.CreatedAt)
		return ""
	}

	return typ
}
