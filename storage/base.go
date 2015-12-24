package storage

import "path"

type BaseStorage struct {
	Location string
}

//NewBaseStorage creates a new base storage
func NewBaseStorage(location string) *BaseStorage {
	return &BaseStorage{
		Location: location,
	}
}

//Path makes the absolte path for the file
func (s *BaseStorage) PathWithPrefix(filepath string) string {
	return path.Join(s.Location, filepath)
}
