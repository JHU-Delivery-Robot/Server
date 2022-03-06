package main

import (
	"context"
	"log"
	"time"

	pb "github.com/JHU-Delivery-Robot/Server/spec"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.Dial(":9000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()
	c := pb.NewRouterClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r1, err1 := c.SendStatus(ctx, &pb.RobotStatus{Id: pb.RobotStatus_VIRTUAL, Message: "first status", Location: &pb.Point{Longitude: 0, Latitude: 0}})
	if err1 != nil {
		log.Fatalf("could not get status: %v", err1)
	}
	log.Printf("Errors: %d", r1.GetErrors())

	r2, err2 := c.GetRoute(ctx, &pb.Coords{X: 2323, Y: 2})
	if err2 != nil {
		log.Fatalf("could not get path: %v", err2)
	}
	log.Printf("Lat: %f Lon: %f", r2.GetLatitude(), r2.GetLongitude())
}
