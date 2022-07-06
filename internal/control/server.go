package control

import (
	"context"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/JHU-Delivery-Robot/Server/internal/routing"
	pb "github.com/JHU-Delivery-Robot/Server/spec"
)

type Server struct {
	pb.UnimplementedRouterServer
}

func (s *Server) GetRoute(ctx context.Context, in *pb.Coords) (*pb.RoutePoint, error) {
	log.Printf("Received OSRM Coords: (%v, %v)\n", in.GetX(), in.GetY())

	//data, err := ioutil.ReadFile("./testRoute.txt")

	resp, err := http.Get("http://osrm:5000/route/v1/driving/-76.6215169429779,39.32603594879762;-76.62238597869873,39.328870107766434?overview=full&steps=true&geometries=geojson")

	if err == nil {
		defer resp.Body.Close()
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

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
