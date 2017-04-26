package store

// Store has multiple superblocks. Each superblock has one index.
// File list:
// basePath/gofs/Manifest includes all superblock:bucketName list
// basePath/gofs/sblock_00001 is superblock File
// basePath/gofs/sblock_00001.index is the index file of superblock
type Store struct {
	basePath string
}
