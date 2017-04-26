package meta

type Key []byte

type KV interface {
	Get(Key) ([]byte, error)
	Put(Key, []byte) error
	Del(Key) error
}
