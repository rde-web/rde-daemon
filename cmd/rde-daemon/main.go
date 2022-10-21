package main

import (
	"log"
	"net"

	. "rde-daemon/internal/daemon"
	. "rde-daemon/internal/protocol/rde_fs"

	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":9090")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption

	grpcServer := grpc.NewServer(opts...)
	RegisterRDEFSServer(grpcServer, FS{})
	grpcServer.Serve(lis)
}
