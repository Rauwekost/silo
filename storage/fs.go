package storage

import (
	"io/ioutil"
	"os"
	"strings"
)

var DefaultFilePermissions os.FileMode = 0755

//FileSystemStorage is a file system storage handler
type FileSystemStorage struct {
	*BaseStorage
}

//NewFileSystemStorage returns a file system storage engine
func NewFileSystemStorage(location string) Storage {
	return &FileSystemStorage{
		&BaseStorage{
			Location: location,
		},
	}
}

//Save a file
func (s *FileSystemStorage) Save(path string, f File) error {
	return s.SaveWithPermissions(path, f, DefaultFilePermissions)
}

//Save with file permissions
func (s *FileSystemStorage) SaveWithPermissions(path string, f File, perm os.FileMode) error {
	_, err := os.Stat(s.Location)
	if err != nil {
		return err
	}

	location := s.PathWithPrefix(path)
	basename := location[:strings.LastIndex(location, "/")+1]
	err = os.MkdirAll(basename, perm)
	if err != nil {
		return err
	}

	content, err := f.ReadAll()
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(location, content, perm)

	return err
}

//Open file and get contents
func (s *FileSystemStorage) Open(path string) (File, error) {
	f, err := os.Open(s.PathWithPrefix(path))
	if err != nil {
		return nil, err
	}
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	return NewContentFile(b), nil
}

//Delete a file
func (s *FileSystemStorage) Delete(path string) error {
	return os.Remove(s.PathWithPrefix(path))
}

//Exists checks if a file exists
func (s *FileSystemStorage) Exists(path string) bool {
	_, err := os.Stat(s.PathWithPrefix(path))
	return err == nil
}
