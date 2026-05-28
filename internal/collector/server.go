package collector

import (
	"io"
	"log"

	"github.com/simonarnost/logpulse/gen/logpulsepb"
	"google.golang.org/grpc"
)

type Server struct {
	logpulsepb.UnimplementedCollectorServiceServer
}

func (s *Server) Ingest(stream grpc.ClientStreamingServer[logpulsepb.LogEntry, logpulsepb.IngestResponse]) error {
	var count int64
	for {
		entry, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&logpulsepb.IngestResponse{LogCount: count})
		}
		if err != nil {
			return err
		}
		count++
		log.Printf("received: service=%s level=%s msg=%s", entry.Service, entry.Level, entry.Message)
	}
}
