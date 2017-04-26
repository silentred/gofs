package store

// Index item in index file
type Index struct {
	ID     uint64
	Offset uint32
	Size   uint32
}

type indexProvider interface {
	// LoadIndex during recovery from failure
	LoadIndex(string) error
	// Get ID of next needle
	NextID() uint32
	FindByID(id uint64) *Index
	// Persistent the index
	Append(Index) bool
}

// IndexManager manages index. Read, write index, give id to needle
type IndexManager struct {
	provider *indexProvider
}

func NewProvider(provider string) indexProvider {
	return nil
}

// ===== File provider =====
// every superblock has one Index File
