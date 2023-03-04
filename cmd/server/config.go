package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/go-playground/validator/v10"
)

type TLSCredentials struct {
	RootCA      string `json:"rootCA" validate:"required"`
	Certificate string `json:"certificate" validate:"required"`
	Key         string `json:"key" validate:"required"`
}

type Config struct {
	Credentials     TLSCredentials `json:"credentials" validate:"required"`
	GRPCListen      string         `json:"gRPCListen" validate:"required"`
	RESTListen      string         `json:"restListen" validate:"required"`
	OSRMAddress     string         `json:"OSRMAddress" validate:"required"`
	OSRMPRofileName string         `json:"OSRMProfileName" validate:"required"`
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
