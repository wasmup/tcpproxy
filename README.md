# tcpproxy


## Install
```sh
go install -x -trimpath=true -ldflags=-s github.com/wasmup/tcpproxy@latest

file $(which tcpproxy)
```


## Run
```sh
tcpproxy localhost:6379 remoteIP:6379

tcpproxy localhost:8181 localhost:8080

tcpproxy localhost:8181 localhost:8080 -log 

```