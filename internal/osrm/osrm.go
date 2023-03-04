package osrm

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/JHU-Delivery-Robot/Server/internal/store"
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

type routeResponse struct {
	Code    string  `json:"code"`
	Message *string `json:"message"`
	Routes  []route `json:"routes"`
}

func (r routeResponse) String() string {
	if r.Message != nil {
		return fmt.Sprintf("%s: %s", r.Code, *r.Message)
	}

	return r.Code
}

type Client struct {
	baseURL         string
	profileName     string
	maxResponseSize int64
	client          *http.Client
}

func New(osrmAddress string, profileName string) Client {
	return Client{
		baseURL:         osrmAddress,
		profileName:     profileName,
		maxResponseSize: 1.2e+6,
		client: &http.Client{
			Timeout: time.Second * 10,
		},
	}
}

func (c *Client) Route(ctx context.Context, start store.Point, end store.Point) ([]store.Point, error) {
	start_string := strconv.FormatFloat(start.Longitude, 'f', 5, 64) + "," + strconv.FormatFloat(start.Latitude, 'f', 5, 64)
	end_string := strconv.FormatFloat(end.Longitude, 'f', 5, 64) + "," + strconv.FormatFloat(end.Latitude, 'f', 5, 64)

	url := c.baseURL + "/route/v1/" + c.profileName + "/" + start_string + ";" + end_string
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

	response, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusBadRequest {
		return nil, fmt.Errorf("HTTP %s", response.Status)
	}

	defer response.Body.Close()

	var osrm_response routeResponse
	if err := json.NewDecoder(io.LimitReader(response.Body, c.maxResponseSize)).Decode(&osrm_response); err != nil {
		log.Println(err)
		return nil, err
	}

	if response.StatusCode == http.StatusBadRequest {
		return nil, fmt.Errorf("OSRM error %s", osrm_response)
	}

	var points = osrm_response.Routes[0].Geometry.Coordinates
	var waypoints = make([]store.Point, len(points))

	for i := 0; i < len(points); i++ {
		waypoint := store.Point{
			Longitude: points[i][0],
			Latitude:  points[i][1],
		}
		waypoints[i] = waypoint
	}

	return waypoints, nil
}
