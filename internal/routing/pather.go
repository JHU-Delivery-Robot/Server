package routing

type OSRMRoute struct {
	Longitude float32
	Latitude  float32
}

type osrmRouteResponse struct {
	Code   string  `json:"code"`
	Routes []route `json"routes"`
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

func (OSRMRoute) GetOSRMRoute() {

}
