package store

var (
	// BasePath is program's data dir; TODO config
	BasePath = "/tmp"
)

// Store has multiple superblocks. Each superblock has one index.
// File list:
// basePath/gofs/Manifest includes all superblock:bucketName list
// basePath/gofs/sblock_00001 is superblock File
// basePath/gofs/sblock_00001.index is the index file of superblock
type Store struct {
	id       int
	basePath string
	blocks   []*Superblock
	status   StoreStatus
}

type StoreStatus struct {
	DiskUsage DiskStatus `json:"disk_usage"`
	BlockNum  int        `json:"block_cnt"`
}

type BlockStatus struct {
	NeedleNum int  `json:"needle_cnt"`
	Writeable bool `json:"writeable"`
}
