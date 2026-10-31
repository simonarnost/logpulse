package main

import (
	"bufio"
	"context"
	"io"
	"log"
	"os"

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

// entrySender is the subset of the ingest stream streamLines needs; easy to fake in tests.
type entrySender interface {
	Send(*logpulsepb.LogEntry) error
}

// streamLines sends each non-empty line from r as a LogEntry.
func streamLines(r io.Reader, stream entrySender) error {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		if err := stream.Send(&logpulsepb.LogEntry{
			Level:   logpulsepb.LogLevel_INFO,
			Message: line,
			Service: "agent007",
		}); err != nil {
			return err
		}
	}
	// Scan returns false on both EOF and error; Err is nil on clean EOF.
	return scanner.Err()
}

func main() {
	conn := mustDial()
	// Process exits; close error not actionable.
	defer conn.Close()
	client := logpulsepb.NewCollectorServiceClient(conn)
	stream := mustOpenStream(client)

	if err := streamLines(os.Stdin, stream); err != nil {
		log.Fatalf("streaming lines failed: %v", err)
	}

	res := mustFinish(stream)
	log.Printf("collector received %d log messages", res.LogCount)
}
