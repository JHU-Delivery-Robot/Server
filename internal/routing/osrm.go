package routing

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

type osrmRouter struct {
	baseURL         string
	profileName     string
	maxResponseSize int64
	client          *http.Client
}

func NewOSRMRouter() osrmRouter {
	return osrmRouter{
		baseURL:         "http://osrm:5000",
		profileName:     "wheelchair_elektro",
		maxResponseSize: 1.2e+6,
		client: &http.Client{
			Timeout: time.Second * 10,
		},
	}
}

func (r *osrmRouter) Route(ctx context.Context, start *pb.Point, end *pb.Point) (*Route, error) {
	start_string := strconv.FormatFloat(start.Longitude, 'f', 5, 64) + "," + strconv.FormatFloat(start.Latitude, 'f', 5, 64)
	end_string := strconv.FormatFloat(end.Longitude, 'f', 5, 64) + "," + strconv.FormatFloat(end.Latitude, 'f', 5, 64)

	url := r.baseURL + "/route/v1/" + r.profileName + "/" + start_string + ";" + end_string
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

	response, err := r.client.Do(req)
	if err != nil || response.StatusCode != http.StatusOK {
		return nil, err
	}

	defer response.Body.Close()

	var osrm_response osrmRouteResponse
	if err := json.NewDecoder(io.LimitReader(response.Body, r.maxResponseSize)).Decode(&osrm_response); err != nil {
		log.Print(err)
		return nil, err
	}

	var route Route
	var points = osrm_response.Routes[0].Geometry.Coordinates
	for i := 0; i < len(points); i++ {
		var waypoint pb.Point
		waypoint.Longitude = points[i][0]
		waypoint.Latitude = points[i][1]
		route.Waypoints = append(route.Waypoints, &waypoint)
	}

	return &route, nil
}
