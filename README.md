# tcpproxy


## Install
```sh
go install -x -ldflags=-s github.com/wasmup/tcpproxy@latest

file $(which tcpproxy)
```


## Run
```sh
tcpproxy localhost:6379 remoteIP:6379
```