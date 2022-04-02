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

	data, err := ioutil.ReadFile("./testRoute.txt")

	if err == nil {
		var nextRoute routing.OSRMRoute = routing.GetOSRMRoute(data)
		var points []*pb.Point
		for i := 0; i < len(nextRoute.Latitude); i++ {

			points = append(points, &pb.Point{Longitude: nextRoute.Longitude[i], Latitude: nextRoute.Latitude[i]})
		}
		return &pb.RoutePoint{Waypoints: points}, nil
	} else {
		return nil, err
	}

}

func (s *Server) SendStatus(ctx context.Context, in *pb.RobotStatus) (*pb.StatusConfirmation, error) {
	log.Print("Recieved robot status\n")
	return &pb.StatusConfirmation{Errors: 0}, nil
}
