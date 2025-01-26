// logger.go
package main

import (
	"context"
	"fmt"
	"log"
	pb "map-storage-service/logs" // You'll need to copy the proto files from logger service
	"time"

	"google.golang.org/grpc"
)

type LoggerClient struct {
	client pb.LoggerServiceClient
	conn   *grpc.ClientConn
}

var loggerClient *LoggerClient

func connectToLogger() error {
	conn, err := grpc.Dial("logger-service:50001", grpc.WithInsecure())
	if err != nil {
		return fmt.Errorf("failed to connect to logger service: %v", err)
	}

	client := pb.NewLoggerServiceClient(conn)
	loggerClient = &LoggerClient{
		client: client,
		conn:   conn,
	}

	return nil
}

func (l *LoggerClient) Close() {
	if l.conn != nil {
		l.conn.Close()
	}
}

func (l *LoggerClient) LogEvent(eventType, userId, description string, metadata map[string]string) error {
	if l == nil || l.client == nil {
		return fmt.Errorf("logger client not initialized")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	_, err := l.client.LogEvent(ctx, &pb.LogRequest{
		ServiceName: "map_storage",
		EventType:   eventType,
		Description: description,
		Severity:    1,
		Timestamp:   time.Now().Format(time.RFC3339),
		Metadata:    metadata,
	})

	if err != nil {
		log.Printf("Failed to send log: %v", err)
		return err
	}

	return nil
}
