package main

import (
	"log"
	"net"

	"github.com/simonarnost/logpulse/gen/logpulsepb"
	"github.com/simonarnost/logpulse/internal/collector"
	"google.golang.org/grpc"
)

const (
	clientAddr = ":50051"
	protocol   = "tcp"
)

func initializeListener() net.Listener {
	listener, err := net.Listen(protocol, clientAddr)
	if err != nil {
		log.Fatalf("initializing listener failed: %v", err)
	}
	return listener
}

func main() {
	listener := initializeListener()
	grpcServer := grpc.NewServer()
	logpulsepb.RegisterCollectorServiceServer(grpcServer, &collector.Server{})
	log.Printf("collector listening on %s", listener.Addr())
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("starting grpc server failed: %v", err)
	}
}
