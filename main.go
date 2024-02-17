package main

import (
	"fmt"
	"io"
	"log"
	"log/slog"
	"net"
	"os"
	"runtime"
	"sync/atomic"
	"time"
)

func main() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{AddSource: true})))
	slog.Info(runtime.Version())

	a := os.Args[1:]
	info := len(a) > 2 && a[2] == "-log"
	if info {
		a = a[:2]
	}
	if len(a) != 2 {
		fmt.Println("Run:\ntcpproxy localhost:6379 remote:6379")
		fmt.Println("Run:\ntcpproxy localhost:6379 remote:6379 -log")
		return
	}
	src := a[0]
	dst := a[1]

	listener, err := net.Listen("tcp", src)
	if err != nil {
		slog.Error("listening", "src", src, "err", err)
		return
	}
	defer listener.Close()

	slog.Info("Proxy server started. Listening", "src", src)

	var read, write, count atomic.Int64
	if info {
		go func() {
			for {
				time.Sleep(1 * time.Second)
				fmt.Println("bytes read", read.Load(), "write", write.Load(), "count", count.Load())
			}
		}()
	}

	for {
		local, err := listener.Accept()
		if err != nil {
			slog.Error("accepting connection", "err", err)
			continue
		}

		go func(local net.Conn) {
			defer local.Close()
			remote, err := net.Dial("tcp", dst)
			if err != nil {
				slog.Error("connecting to remote Redis server", "err", err)
				return
			}
			defer remote.Close()
			if info {
				log.Println(local.RemoteAddr(), "<---->", local.LocalAddr(), "=Go=", remote.LocalAddr(), "<----->", remote.RemoteAddr())
			}
			count.Add(1)

			go func() {
				defer local.Close()
				defer remote.Close()

				n, err := io.Copy(remote, local)
				if err != nil {
					slog.Error("copying data from local to remote", "err", err)
				}
				write.Add(n)
			}()

			n, err := io.Copy(local, remote)
			if err != nil {
				slog.Error("copying data from remote to local", "err", err)
			}
			read.Add(n)
			count.Add(-1)
		}(local)
	}
}
