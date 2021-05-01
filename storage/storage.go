package storage

type Storage interface {
	Contains(string) bool
	Add(string)
	List() []string
}
