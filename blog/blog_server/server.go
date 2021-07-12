package main

import (
	"blog_server/blogpb"
	"blog_server/repository"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"os"
	"os/signal"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func connect() (*mongo.Client, *mongo.Collection) {
	fmt.Println("mongo DB started")
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}

	err = client.Connect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	return client, client.Database("mydb").Collection("blog")
}

func main() {
	// mongo
	client, collection := connect()

	fmt.Println("Blog Server Started!")
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	sm := repository.NewBlogItemServerMongo(collection)
	var opts []grpc.ServerOption
	s := grpc.NewServer(opts...)
	blogpb.RegisterBlogServiceServer(s, sm)
	// Register reflection service on gRPC server
	reflection.Register(s)

	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve : %v", err)
		}
	}()

	// wait for control C to stop
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	// block until a signal is received
	<-ch
	fmt.Println("stopping the server")
	s.Stop()
	fmt.Println("Closing the lister")
	lis.Close()
	fmt.Println("closing the mongodb connection")
	client.Disconnect(context.TODO())
	fmt.Println("End of program")
}
