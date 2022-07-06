package routing

import (
	"encoding/json"
)

type OSRMRoute struct {
	Longitude []float64
	Latitude  []float64
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
	Coordinates []coordinates `json:"coordinates"`
}

type coordinates []float64

func GetOSRMRoute(outputJson []byte) OSRMRoute {
	var jsonParse osrmRouteResponse
	json.Unmarshal(outputJson, &jsonParse)

	var information OSRMRoute

	for i := 0; i < len(jsonParse.Routes[0].Geometry.Coordinates); i++ {
		information.Longitude = append(information.Longitude, jsonParse.Routes[0].Geometry.Coordinates[i][0])
		information.Latitude = append(information.Latitude, jsonParse.Routes[0].Geometry.Coordinates[i][1])
	}

	return information
}
