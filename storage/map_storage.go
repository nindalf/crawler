package storage

type mapStorage struct {
	m map[string]bool
}

func NewMapStorage() Storage {
	m := make(map[string]bool)
	return mapStorage{m}
}

func (ms mapStorage) Contains(s string) bool {
	_, ok := ms.m[s]
	return ok
}

func (ms mapStorage) Add(s string) {
	ms.m[s] = true
}

func (ms mapStorage) List() []string {
	result := make([]string, 0, len(ms.m))
	for s := range ms.m {
		result = append(result, s)
	}
	return result
}
