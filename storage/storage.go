package storage

//Storage is an interface for multiple storage engines
type Storage interface {
	Save(path string, f File) error
	Exists(path string) bool
	Delete(path string) error
	Open(path string) (File, error)
	PathWithPrefix(path string) string
}
