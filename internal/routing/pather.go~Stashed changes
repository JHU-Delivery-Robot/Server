package routing

import (
	"encoding/json"
)

type OSRMRoute struct {
	Longitude float64
	Latitude  float64
}

type osrmRouteResponse struct {
	Code   string  `json:"code"`
	Routes []route `json:"routes"`
}

type route struct {
	Geometry geometry `json:"geometry"`
	Distance float32  `json:"distance"`
	Duration float32  `json:"duration"`
}

type geometry struct {
	Coordinates coordinates `json:"coordinates"`
}

type coordinates []float64

func GetOSRMRoute(outputJson string) OSRMRoute {
	var jsonParse osrmRouteResponse
	json.Unmarshal([]byte(outputJson), &jsonParse)
	var information OSRMRoute

	information.Longitude = jsonParse.Routes[1].Geometry.Coordinates[0]
	information.Latitude = jsonParse.Routes[1].Geometry.Coordinates[1]

	return information
}
