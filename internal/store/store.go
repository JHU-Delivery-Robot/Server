package store

import (
	"fmt"

	"github.com/hashicorp/go-memdb"
)

type Store struct {
	// map from identifiers to Robots
	database *memdb.MemDB
}

type Point struct {
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
}

const (
	tableAssignments = "assignments"
	tableRequests    = "requests"
	tableRobots      = "robots"
)

func New() (*Store, error) {
	schema := &memdb.DBSchema{
		Tables: map[string]*memdb.TableSchema{
			tableAssignments: assignmentsTableSchema,
			tableRequests:    requestsTableSchema,
			tableRobots:      robotsTableSchema,
		},
	}

	err := schema.Validate()
	if err != nil {
		return nil, fmt.Errorf("validating database schema: %v", err)
	}

	db, err := memdb.NewMemDB(schema)
	if err != nil {
		return nil, fmt.Errorf("creating memdb database: %v", err)
	}

	return &Store{
		database: db,
	}, nil
}
