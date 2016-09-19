package store

type Index struct {
	ID     uint64
	Offset uint32
	Size   uint32
}

type indexProvider interface {
	Last() *Index
	FindByID(id uint64) *Index
	Append(*Index) bool
}

// IndexManager manages index
type IndexManager struct {
	provider *indexProvider
}
