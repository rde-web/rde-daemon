package main

import (
	"errors"
	"io"
	"log"
	"net"

	. "rde-daemon/internal/daemon"
	. "rde-daemon/internal/protocol/rde_fs"

	"google.golang.org/grpc"
)

func main() {
	const daemonAddr string = "localhost:46046"
	lis, err := net.Listen("tcp", daemonAddr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	defer lis.Close()
	var opts []grpc.ServerOption = make([]grpc.ServerOption, 0)
	opts = append(opts)

	grpcServer := grpc.NewServer(opts...)
	RegisterRDEFSServer(grpcServer, FS{})
	go grpcServer.Serve(lis)

	proxyLis, err := net.Listen("tcp", ":9090")
	if err != nil {
		log.Fatalf("failed to listen PROXY: %v", err)
	}
	defer proxyLis.Close()
	for {
		conn, err := proxyLis.Accept()
		if err != nil {
			println("err accept PROXY", err.Error())
			break
		}
		sock, err := net.Dial("tcp", daemonAddr)
		if err != nil {
			println("err dial DAEMON", err.Error())
			break
		}
		for {
			var buffp []byte = make([]byte, 1024)
			nr, err := conn.Read(buffp)
			if err != nil {
				if errors.Is(err, io.EOF) {
					continue
				}
				println("err read PROXY", err.Error())
				break
			}
			println("readed PRXOY", nr, string(buffp))

			if _, err := sock.Write(buffp); err != nil {
				println("err write DAEMON", err.Error())
				lis.Close()
				break
			}
			var buffd []byte = make([]byte, 1024)
			if nr, err = sock.Read(buffd); err != nil {
				if errors.Is(err, io.EOF) {
					continue
				}
				println("err read DAEMON", err.Error())
				break
			}
			println("readed DAEMON", nr, string(buffd))
			if _, err := conn.Write(buffd); err != nil {
				println("err write PROXY", err.Error())
				break
			}

		}
		conn.Close()
		sock.Close()
	}
}
