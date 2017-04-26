# Go Filesystem
Distribted Object Store implemented by golang, inspired by Haystack

## Roadmap
- Store
    - [x] needle
    - [ ] superblock
    - [ ] store
    - [ ] index
    - [ ] compact
    - [ ] metrics
- Meta
    - [ ] manage store meta data in etcd
    - [ ] manage bucket meta data
    - [ ] manage file meta data in KV
    - [ ] metrics
- Proxy
    - [ ] file CRUD (interact with Meta and Store)
- Pitchfork
    - use prometheus instead
- Cache
    - use nginx
- Cmd Tool
    - [ ] validate superblock
    - [ ] validate index
    - [ ] cluster status

## Design

For simplicity, combine Store, Meta, Proxy in one binary. Cache could handle most requests and protect Store, so both read and write are not frequent operations. Proxy could barely be the bottle neck, because Store reads and writes directly to file system.
