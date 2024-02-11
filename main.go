package main

import (
	"io"
	"log"
	"net"
	"os"
)

func main() {
	if len(os.Args) != 3 {
		log.Println("Run:\ntcpproxy localhost:6379 remote:6379")
		return
	}

	src := os.Args[1]
	listener, err := net.Listen("tcp", src)
	if err != nil {
		log.Println("Error listening on"+src, err)
		return
	}
	defer listener.Close()

	log.Println("Proxy server started. Listening on " + src)

	dst := os.Args[2]
	for {
		local, err := listener.Accept()
		if err != nil {
			log.Println("Error accepting connection:", err)
			continue
		}

		go func(local net.Conn) {
			defer local.Close()
			remote, err := net.Dial("tcp", dst)
			if err != nil {
				log.Println("Error connecting to remote Redis server:", err)
				return
			}
			defer remote.Close()
			log.Println(local.RemoteAddr(), "<---->", local.LocalAddr(), "=Go=", remote.LocalAddr(), "<----->", remote.RemoteAddr())

			go func() {
				_, err := io.Copy(remote, local)
				if err != nil {
					log.Println("Error copying data from local to remote:", err)
				}
			}()

			_, err = io.Copy(local, remote)
			if err != nil {
				log.Println("Error copying data from remote to local:", err)
			}
		}(local)
	}
}
