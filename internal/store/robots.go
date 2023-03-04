package store

import (
	"fmt"

	"github.com/hashicorp/go-memdb"
)

type Robot struct {
	ID       string `json:"id"`
	Status   string `json:"status"`
	Location Point  `json:"location"`
}

var robotsTableSchema = &memdb.TableSchema{
	Name: tableRobots,
	Indexes: map[string]*memdb.IndexSchema{
		"id": {
			Name:    "id",
			Unique:  true,
			Indexer: &memdb.StringFieldIndex{Field: "ID"},
		},
	},
}

func (s *Store) GetRobots() ([]Robot, error) {
	txn := s.database.Txn(false)
	defer txn.Abort()

	results, err := txn.Get(tableRobots, "id")
	if err != nil {
		return nil, fmt.Errorf("listing robots: %v", err)
	}

	var robots []Robot

	for obj := results.Next(); obj != nil; obj = results.Next() {
		robot := obj.(Robot)
		robots = append(robots, robot)
	}

	return robots, nil
}

func (s *Store) UpsertRobot(robot Robot) error {
	txn := s.database.Txn(true)
	defer txn.Commit()

	if err := txn.Insert(tableRobots, robot); err != nil {
		return fmt.Errorf("upserting robot: %v", err)
	}

	return nil
}
