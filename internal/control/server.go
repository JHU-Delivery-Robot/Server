package control

import (
	"context"
	"io/ioutil"
	"log"

	"github.com/JHU-Delivery-Robot/Server/internal/routing"
	pb "github.com/JHU-Delivery-Robot/Server/spec"
)

type Server struct {
	pb.UnimplementedRouterServer
}

func (s *Server) GetRoute(ctx context.Context, in *pb.Coords) (*pb.RoutePoint, error) {
	log.Printf("Received Coords: (%v, %v)\n", in.GetX(), in.GetY())

	data, err := ioutil.ReadFile("testRoute.txt")
	if err == nil {
		var nextRoute routing.OSRMRoute = routing.GetOSRMRoute(string(data))
		return &pb.RoutePoint{Longitude: nextRoute.Longitude, Latitude: nextRoute.Latitude}, nil
	} else {
		return nil, err
	}

}

func (s *Server) SendStatus(ctx context.Context, in *pb.RobotStatus) (*pb.StatusConfirmation, error) {
	log.Print("Recieved robot status\n")
	return &pb.StatusConfirmation{Errors: 0}, nil
}
