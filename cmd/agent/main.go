package main

import (
	"context"
	"fmt"
	"log"

	"github.com/simonarnost/logpulse/gen/logpulsepb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const collectorAddr = "localhost:50051"

func mustDial() *grpc.ClientConn {
	conn, err := grpc.NewClient(collectorAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("initializing client connection failed: %v", err)
	}
	return conn
}

func mustOpenStream(client logpulsepb.CollectorServiceClient) logpulsepb.CollectorService_IngestClient {
	stream, err := client.Ingest(context.Background())
	if err != nil {
		log.Fatalf("opening ingest stream failed: %v", err)
	}
	return stream
}

func mustFinish(stream logpulsepb.CollectorService_IngestClient) *logpulsepb.IngestResponse {
	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("stream close and recv failed: %v", err)
	}
	return res
}

func main() {
	conn := mustDial()
	// Process exits; close error not actionable.
	defer conn.Close()
	client := logpulsepb.NewCollectorServiceClient(conn)
	stream := mustOpenStream(client)

	// todo -- rewrite mocked logic
	for i := 0; i < 3; i++ {
		if err := stream.Send(&logpulsepb.LogEntry{
			Level:   logpulsepb.LogLevel_INFO,
			Message: fmt.Sprintf("agent sending log message %d", i),
			Service: "agent007",
		}); err != nil {
			log.Fatalf("sending message through stream failed: %v", err)
		}
	}

	res := mustFinish(stream)
	log.Printf("collector received %d log messages", res.LogCount)
}
