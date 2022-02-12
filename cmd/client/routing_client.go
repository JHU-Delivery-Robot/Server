package main

import (
	"context"
	"log"
	"time"

	pb "github.com/JHU-Delivery-Robot/Server/spec"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial(":9000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()
	c := pb.NewRouterClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.GetRoute(ctx, &pb.Coords{X: 2323, Y: 2})
	if err != nil {
		log.Fatalf("could not get path: %v", err)
	}
	log.Printf("Path: %s", r.GetPath())
}
