<p align="center">
<img
    src="https://redislabs.com/wp-content/uploads/2018/03/golang-redis.jpg"
    width="466" height="265" border="0">
<br>
</p>

<p align="center"><b>Becoming a full Redis implementation in Go</b></p>

This project started to see how easy it is to implement a full Redis clone in Go.
As one of the side effects, imagine you could write redis modules in Go, that would be awesome!

# Get involved!
This project is in *work-in-progress*, so share ideas, code and have fun.

The goal is to have all features and commands like the actual [redis](https://github.com/antirez/redis) written in C have.
We are searching contributors!


### Documentation

godoc: https://godoc.org/github.com/redis-go/redis

### Getting Started

You can already test out the API.

To install, run:
```bash
go get -u github.com/redis-go/redis
```


### Roadmap
- [x] Client connection / request / respond
- [x] RESP protocol
- [x] able to register commands
- [x] in-mem database
- [x] active key expirer
- [ ] Implementing data structures
  - [x] String
  - [x] List
  - [ ] Set
  - [ ] Sorted Set
  - [ ] Hash
  - [ ] ...
- [ ] Tests
  - [x] For existing commands
  - [x] For key expirer
- [ ] Alpha Release

### TODO beside Roadmap
- [ ] Persistence
- [ ] Redis config
  - [ ] Default redis config format
  - [ ] YAML support
  - [ ] Json support
- [ ] Pub/Sub
- [ ] Redis modules
- [ ] Benchmarks
- [ ] master slaves
- [ ] cluster
- [ ] ...
