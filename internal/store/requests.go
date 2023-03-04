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

func (s *Store) GetAllRequests() ([]Request, error) {
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

func (s *Store) GetIncompleteRequests() ([]Request, error) {
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

func (s *Store) DeleteRequest(requestID string) (bool, error) {
	txn := s.database.Txn(true)
	defer txn.Commit()

	count, err := txn.DeleteAll(tableRequests, "requestID", requestID)
	if err != nil {
		return false, fmt.Errorf("deleting request %v: %v", requestID, err)
	}

	if count == 0 {
		return false, nil
	} else if count == 1 {
		return true, nil
	} else {
		return true, fmt.Errorf("found duplicates while deleting request %v: %v", requestID, err)
	}
}
