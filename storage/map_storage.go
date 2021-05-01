package storage

type MapStorage struct {
	m map[string]bool
}

func NewMapStorage() Storage {
	m := make(map[string]bool)
	return MapStorage{m}
}

func (ms MapStorage) Contains(s string) bool {
	_, ok := ms.m[s]
	return ok
}

func (ms MapStorage) Add(s string) {
	ms.m[s] = true
}

func (ms MapStorage) List() []string {
	result := make([]string, 0, len(ms.m))
	for s := range ms.m {
		result = append(result, s)
	}
	return result
}
