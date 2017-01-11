#Yet another fake server in go
Well I wanted to create something to start playing a bit with go, and I said to me...what's better than a fake server configurable via yaml? That is the story, nothing more and nothing less!

## Install dependencies
```go
go get github.com/gorilla/mux
go get gopkg.in/yaml.v2
```

## First config
Rename `config.yaml.dist` to `config.yaml`

## Build and run it
```go
go build && ./yafsig
```