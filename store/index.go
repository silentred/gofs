package store

type Index struct {
	ID     uint64
	Offset uint32
	Size   uint32
}

type indexProvider interface {
	NextID() uint32
	FindByID(id uint64) *Index
	Append(*Index) bool
}

// IndexManager manages index
type IndexManager struct {
	provider *indexProvider
}

func NewProvider(provider string) indexProvider {
	return nil
}

// ===== File provider =====
// every superblock has one Index File
