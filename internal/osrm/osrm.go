package osrm

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	pb "github.com/JHU-Delivery-Robot/Server/protocol"
)

type coordinates []float64

type geometry struct {
	Coordinates []coordinates `json:"coordinates"`
}

type route struct {
	Geometry geometry `json:"geometry"`
	Distance float32  `json:"distance"`
	Duration float32  `json:"duration"`
}

type osrmRouteResponse struct {
	Code    string  `json:"code"`
	Message string  `json:"message"`
	Routes  []route `json:"routes"`
}

type Route struct {
	Waypoints []*pb.Point
}

// Max size of response to read from OSRM
const maxResponseSize int64 = 1.2e+6

var osrmClient = &http.Client{
	Timeout: time.Second * 10,
}

const osrmBaseURL = "http://osrm:5000"
const osrmProfileName = "wheelchair_elektro"

func GetRoute(ctx context.Context, start *pb.Point, end *pb.Point) (*Route, error) {
	start_string := strconv.FormatFloat(start.Longitude, 'f', 5, 64) + "," + strconv.FormatFloat(start.Latitude, 'f', 5, 64)
	end_string := strconv.FormatFloat(end.Longitude, 'f', 5, 64) + "," + strconv.FormatFloat(end.Latitude, 'f', 5, 64)

	url := osrmBaseURL + "/route/v1/" + osrmProfileName + "/" + start_string + ";" + end_string
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	req = req.WithContext(ctx)
	query := req.URL.Query()
	query.Add("overview", "full")
	query.Add("geometries", "geojson")
	req.URL.RawQuery = query.Encode()

	response, err := osrmClient.Do(req)
	if err != nil || response.StatusCode != http.StatusOK {
		return nil, err
	}

	defer response.Body.Close()

	var osrm_response osrmRouteResponse
	if err := json.NewDecoder(io.LimitReader(response.Body, maxResponseSize)).Decode(&osrm_response); err != nil {
		log.Print(err)
		return nil, err
	}

	var route Route

	for i := 0; i < len(osrm_response.Routes[0].Geometry.Coordinates); i++ {
		var waypoint pb.Point
		waypoint.Longitude = osrm_response.Routes[0].Geometry.Coordinates[i][0]
		waypoint.Latitude = osrm_response.Routes[0].Geometry.Coordinates[i][1]
		route.Waypoints = append(route.Waypoints, &waypoint)
	}

	return &route, nil
}
