package db

import (
	"github.com/google/uuid"
	"sync"
)

/*
DB is an interface that contains the methods to interact with the database
GetReceipt is a method that returns the points of the receipt
AddNewReceipt is a method that adds a new receipt to the database

*/
type DB interface {
	GetReceipt(id string) (int64, bool)
	AddNewReceipt(points int64) string
}


/*
Just trying to replicate the in memory database

InMemoryDB is a struct that contains the AllReceipts map
AllReceipts is a map that contains the id of the receipt and the points of the receipt
in a thread safe manner

InMemoryDB implements the DB interface
for AddNewReceipt, it generates a new UUID id and adds the points to the AllReceipts map

assumming that the receipt is valid, the points are added to the AllReceipts map
and the id generated is random and unique
*/

var lock = &sync.Mutex{}

type InMemoryDB struct {
	AllReceipts map[string]int64
}

func (db *InMemoryDB) GetReceipt(id string) (int64, bool) {
	lock.Lock()
	defer lock.Unlock()
	points, ok := db.AllReceipts[id]
	return points, ok
}

func (db *InMemoryDB) AddNewReceipt(points int64) string {
	lock.Lock()
	defer lock.Unlock()
	var id string = uuid.New().String()
	db.AllReceipts[id] = points
	return id
}