package store

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/hashicorp/go-memdb"
)

type Assignment struct {
	ID        string  `json:"-"`
	RobotID   string  `json:"robotID"`
	RequestID string  `json:"requestID"`
	Route     []Point `json:"route"`
}

var assignmentsTableSchema = &memdb.TableSchema{
	Name: tableAssignments,
	Indexes: map[string]*memdb.IndexSchema{
		"id": {
			Name:    "id",
			Unique:  true,
			Indexer: &memdb.UUIDFieldIndex{Field: "ID"},
		},
		"robotID": {
			Name:    "robotID",
			Unique:  true,
			Indexer: &memdb.StringFieldIndex{Field: "RobotID"},
		},
		"requestID": {
			Name:    "requestID",
			Unique:  false,
			Indexer: &memdb.StringFieldIndex{Field: "RequestID"},
		},
	},
}

func (s *Store) GetRoute(robotID string) ([]Point, error) {
	txn := s.database.Txn(false)
	defer txn.Abort()

	result, err := txn.First(tableAssignments, "robotID", robotID)
	if err != nil {
		return nil, fmt.Errorf("getting route assignment: %v", err)
	}

	if result == nil {
		return nil, nil
	}

	assignment := result.(Assignment)

	return assignment.Route, nil
}

const RequestNone = "none"

func (s *Store) AddRoute(robotID string, route []Point, requestID string) error {
	txn := s.database.Txn(true)
	defer txn.Commit()

	assignment := Assignment{
		ID:        uuid.New().String(),
		RobotID:   robotID,
		Route:     route,
		RequestID: requestID,
	}

	if err := txn.Insert(tableAssignments, assignment); err != nil {
		return fmt.Errorf("insert assignment: %v", err)
	}

	return nil
}

func (s *Store) IsInProgress(requestID string) (bool, error) {
	txn := s.database.Txn(false)
	defer txn.Abort()

	result, err := txn.First(tableAssignments, "requestID", requestID)
	if err != nil {
		return true, fmt.Errorf("request progress check: %v", err)
	}

	if result == nil {
		return false, nil
	} else {
		return true, nil
	}
}
