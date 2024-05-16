package db

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"sync"
)

type Chirp struct {
    ID   int    `json:"id"`
    Body string `json:"body"`
}

type DB struct {
	path string
	mux  *sync.RWMutex
}
type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
}


// NewDB creates a new database connection
// and creates the database file if it doesn't exist
func NewDB(path string) (*DB, error){

}

// CreateChirp creates a new chirp and saves it to disk
func (db *DB) CreateChirp(body string) (Chirp, error){

}

// GetChirps returns all chirps in the database
func (db *DB) GetChirps() ([]Chirp, error){

}

// ensureDB creates a new database file if it doesn't exist
func (db *DB) ensureDB() error{
//check path on db using Stat
 _,err := os.ReadFile(db.path)

 
 if errors.Is(err, fs.ErrNotExist) {
			//database.json does not exist //create file
 
emptyDb := DBStructure{Chirps: make(map[int]Chirp)}
// use json marshall indent to fill with empty data
 json, err := json.MarshalIndent(emptyDb, "", " ")
if err != nil {
	log.Fatal(err)
}
// if err := writefile; err != nil  {return error} 
if err := os.WriteFile(db.path, json, 0644); err != nil {
 return err
}

 
	}
 //database.json exists
 return nil
}



// loadDB reads the database file into memory
func (db *DB) loadDB() (DBStructure, error){
db.mux.RLock()
defer db.mux.RUnlock()

data, err := os.ReadFile(db.path)
if err != nil {
	return DBStructure{}, fmt.Errorf("error reading database file: %w", err)
}
var dbStructure DBStructure
if err := json.Unmarshal(data, &dbStructure); err != nil {
	return dbStructure, fmt.Errorf("error unmarshalling database: %w", err)
}
return dbStructure, nil
}

// writeDB writes the database file to disk
func (db *DB) writeDB(dbStructure DBStructure) error{

} 





func (db *DB) eeDB()error {
	_, err := os.Stat(db.path)

	if os.IsNotExist(err){
		
	}
}