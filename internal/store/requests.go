package store

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/hashicorp/go-memdb"
)

type Request struct {
	ID       string `json:"id"`
	Location Point  `json:"location"`
	Complete bool   `json:"completed"`
}

var requestsTableSchema = &memdb.TableSchema{
	Name: tableRequests,
	Indexes: map[string]*memdb.IndexSchema{
		"id": {
			Name:    "id",
			Unique:  true,
			Indexer: &memdb.UUIDFieldIndex{Field: "ID"},
		},
		"complete": {
			Name:    "complete",
			Unique:  false,
			Indexer: &memdb.BoolFieldIndex{Field: "Complete"},
		},
	},
}

func (s *Store) GetRequests() ([]Request, error) {
	txn := s.database.Txn(false)
	defer txn.Abort()

	results, err := txn.Get(tableRequests, "id")
	if err != nil {
		return nil, fmt.Errorf("listing requests: %v", err)
	}

	var requests = make([]Request, 0)

	for obj := results.Next(); obj != nil; obj = results.Next() {
		request := obj.(Request)
		requests = append(requests, request)
	}

	return requests, nil
}

func (s *Store) CreateRequest(location Point) (string, error) {
	requestID := uuid.New().String()
	request := Request{
		ID:       requestID,
		Location: location,
		Complete: false,
	}

	txn := s.database.Txn(true)
	defer txn.Commit()

	if err := txn.Insert(tableRequests, request); err != nil {
		return "", fmt.Errorf("insert request: %v", err)
	}

	return requestID, nil
}

func (s *Store) IncompleteRequests() ([]Request, error) {
	txn := s.database.Txn(false)
	defer txn.Abort()

	results, err := txn.Get(tableRequests, "complete", false)
	if err != nil {
		return nil, fmt.Errorf("listing incomplete requests: %v", err)
	}

	var requests []Request

	for obj := results.Next(); obj != nil; obj = results.Next() {
		request := obj.(Request)
		requests = append(requests, request)
	}

	return requests, nil
}
