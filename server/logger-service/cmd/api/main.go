// logger-service/cmd/api/main.go
package main

import (
	"context"
	"encoding/json"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	pb "logger-service/logs" // modified import path

	"github.com/streadway/amqp"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
)

type Log struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	ServiceName string             `bson:"service_name" json:"service_name"`
	EventType   string             `bson:"event_type" json:"event_type"`
	UserID      string             `bson:"user_id" json:"user_id"`
	Description string             `bson:"description" json:"description"`
	Severity    int32              `bson:"severity" json:"severity"`
	Timestamp   time.Time          `bson:"timestamp" json:"timestamp"`
	Metadata    map[string]string  `bson:"metadata" json:"metadata"`
}

type LoggerServer struct {
	pb.UnimplementedLoggerServiceServer // modified to match generated code
	db                                  *mongo.Database
}

type Config struct {
	rabbitmq *amqp.Connection
	mongo    *mongo.Client
}

func (app *Config) setupRabbitMQ() error {
	ch, err := app.rabbitmq.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	exchanges := []string{"auth_logs", "map_coloring_logs", "map_storage_logs"}
	for _, exchange := range exchanges {
		err := ch.ExchangeDeclare(
			exchange,
			"fanout",
			true,
			false,
			false,
			false,
			nil,
		)
		if err != nil {
			return err
		}
	}

	queues := []string{"auth_logs_queue", "map_coloring_logs_queue", "map_storage_logs_queue"}
	for i, queue := range queues {
		q, err := ch.QueueDeclare(
			queue,
			true,
			false,
			false,
			false,
			nil,
		)
		if err != nil {
			return err
		}

		err = ch.QueueBind(
			q.Name,
			"",
			exchanges[i],
			false,
			nil,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func (app *Config) consumeMessages(ctx context.Context) {
	ch, err := app.rabbitmq.Channel()
	if err != nil {
		log.Printf("Failed to create channel: %v", err)
		return
	}
	defer ch.Close()

	queues := []string{"auth_logs_queue", "map_coloring_logs_queue", "map_storage_logs_queue"}

	for _, queue := range queues {
		go func(queueName string) {
			msgs, err := ch.Consume(
				queueName,
				"",
				false,
				false,
				false,
				false,
				nil,
			)
			if err != nil {
				log.Printf("Failed to consume from queue %s: %v", queueName, err)
				return
			}

			for {
				select {
				case msg := <-msgs:
					var logData Log
					if err := json.Unmarshal(msg.Body, &logData); err != nil {
						log.Printf("Error unmarshaling message: %v", err)
						msg.Nack(false, true)
						continue
					}

					_, err := app.mongo.Database("logs").Collection("logs").InsertOne(ctx, logData)
					if err != nil {
						log.Printf("Error saving to MongoDB: %v", err)
						msg.Nack(false, true)
						continue
					}

					msg.Ack(false)
				case <-ctx.Done():
					return
				}
			}
		}(queue)
	}
}

func (s *LoggerServer) LogEvent(ctx context.Context, req *pb.LogRequest) (*pb.LogResponse, error) {
	timestamp, err := time.Parse(time.RFC3339, req.Timestamp)
	if err != nil {
		timestamp = time.Now()
	}

	log := Log{
		ServiceName: req.ServiceName,
		EventType:   req.EventType,
		UserID:      req.UserId,
		Description: req.Description,
		Severity:    req.Severity,
		Timestamp:   timestamp,
		Metadata:    req.Metadata,
	}

	_, err = s.db.Collection("logs").InsertOne(ctx, log)
	if err != nil {
		return &pb.LogResponse{
			Success: false,
			Message: "Failed to save log",
		}, err
	}

	return &pb.LogResponse{
		Success: true,
		Message: "Log saved successfully",
	}, nil
}

func main() {
	// MongoDB connection
	mongoClient, err := connectToMongo()
	if err != nil {
		log.Fatal("Cannot connect to mongo:", err)
	}
	defer mongoClient.Disconnect(context.Background())

	// RabbitMQ connection
	rabbitConn, err := connectToRabbitMQ()
	if err != nil {
		log.Fatal("Cannot connect to RabbitMQ:", err)
	}
	defer rabbitConn.Close()

	app := Config{
		rabbitmq: rabbitConn,
		mongo:    mongoClient,
	}

	// Setup RabbitMQ exchanges and queues
	if err = app.setupRabbitMQ(); err != nil {
		log.Fatal("Cannot setup RabbitMQ:", err)
	}

	// Start consuming messages
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	app.consumeMessages(ctx)

	// Initialize gRPC server
	lis, err := net.Listen("tcp", ":50001") // hardcoded port for example
	if err != nil {
		log.Fatal("Failed to listen:", err)
	}

	grpcServer := grpc.NewServer()
	loggerServer := &LoggerServer{
		db: mongoClient.Database("logs"), // changed to "logs" to match your setup
	}
	pb.RegisterLoggerServiceServer(grpcServer, loggerServer)

	// Handle graceful shutdown
	go func() {
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
		<-signals
		log.Println("Received terminate signal. Shutting down...")
		grpcServer.GracefulStop()
		cancel()
	}()

	log.Printf("Starting gRPC server on port 50001")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal("Failed to serve:", err)
	}
}

func connectToMongo() (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func connectToRabbitMQ() (*amqp.Connection, error) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost")
	if err != nil {
		return nil, err
	}

	return conn, nil
}
