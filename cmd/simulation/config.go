package main

import (
	"encoding/json"
	"fmt"
	"math"
	"os"

	pb "github.com/JHU-Delivery-Robot/Server/protocol"
	"github.com/go-playground/validator/v10"
)

type Waypoint struct {
	X *float64 `json:"x" validate:"required"`
	Y *float64 `json:"y" validate:"required"`
}

type Coordinate struct {
	Longitude *float64 `json:"longitude" validate:"required"`
	Latitude  *float64 `json:"latitude" validate:"required"`
}

type Config struct {
	Route  []Waypoint `json:"route" validate:"required,dive"`
	Origin Coordinate `json:"origin" validate:"required"`
}

// mean radius from https://en.wikipedia.org/wiki/Earth
const earth_mean_radius_m = 6371000.0

func (c Config) reproject(w Waypoint) pb.Point {
	x := *w.X / earth_mean_radius_m
	y := *w.Y / earth_mean_radius_m

	theta := math.Atan(math.Sinh(x) / math.Cos(y))
	phi := math.Asin(math.Sin(y) / math.Cosh(x))

	return pb.Point{Latitude: 180.0 * (phi / math.Pi), Longitude: 180.0*(theta/math.Pi) + *c.Origin.Longitude}
}

func (c Config) GetRoute() []*pb.Point {
	var points = make([]*pb.Point, len(c.Route))

	for i := 0; i < len(c.Route); i++ {
		waypoint := c.reproject(c.Route[i])
		points[i] = &waypoint
	}

	return points
}

func LoadConfig(path string) (Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return Config{}, fmt.Errorf("unable to open config file: %w", err)
	}
	defer file.Close()

	var config Config
	if err := json.NewDecoder(file).Decode(&config); err != nil {
		return Config{}, fmt.Errorf("unable to parse config file: %w", err)
	}

	validate := validator.New()
	if err := validate.Struct(config); err != nil {
		return Config{}, fmt.Errorf("config file failed validation: %w", err)
	}

	return config, nil
}
