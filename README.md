# Go Filesystem
Distribted Object Store implemented by golang, inspired by Haystack

![](https://img.shields.io/badge/language-golang-blue.svg)
![](https://img.shields.io/badge/license-MIT-000000.svg)
![](https://img.shields.io/github/tag/silentred/gofs.svg)
[![codebeat badge](https://codebeat.co/badges/d8a8e19c-ee50-419d-8bf1-60133dfe446e)](https://codebeat.co/projects/github-com-silentred-gofs-master)


## Roadmap
- Store
    - [x] needle
    - [ ] superblock
    - [ ] store
    - [x] index
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
