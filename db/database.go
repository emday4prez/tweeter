package database

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"sort"
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
db := &DB{
		path: path,
		mux: &sync.RWMutex{},
}
   // create the file if it doesn't exist
    if err := db.ensureDB(); err != nil {
        return nil, err
    }
				return db, nil
}

// CreateChirp creates a new chirp and saves it to disk
func (db *DB) CreateChirp(body string) (Chirp, error){
	dbStruct, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}

	db.mux.Lock()
	defer db.mux.Unlock()

    nextId := 1
    for {
        if _, exists := dbStruct.Chirps[nextId]; !exists {
            break // Found an unused ID
        }
        nextId++
    }

    newChirp := Chirp{
        ID:   nextId,
        Body: body,
    }
    dbStruct.Chirps[nextId] = newChirp
    

    return newChirp, db.writeDB(dbStruct)  //
}

// GetChirps returns all chirps in the database
func (db *DB) GetChirps() ([]Chirp, error){
	dbStruct, err := db.loadDB()
	if err != nil {
		return nil, err
	}

	db.mux.Lock()
	defer db.mux.Unlock()

 chirps := make([]Chirp, 0, len(dbStruct.Chirps))
    for _, chirp := range dbStruct.Chirps {
        chirps = append(chirps, chirp)
    }
    return chirps, nil

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
    db.mux.Lock()
    defer db.mux.Unlock()
    // Sort Chirps
    chirps := make([]Chirp, 0, len(dbStructure.Chirps))
    for _, chirp := range dbStructure.Chirps {
        chirps = append(chirps, chirp)
    }

    // Sort the chirps by ID
    sort.Slice(chirps, func(i, j int) bool {
        return chirps[i].ID < chirps[j].ID
    })

    // Rebuild the map with sorted chirps  
    dbStructure.Chirps = make(map[int]Chirp)
    for _, chirp := range chirps {
        dbStructure.Chirps[chirp.ID] = chirp
    }

    jsonData, err := json.MarshalIndent(dbStructure, "", " ")
    if err != nil {
        return fmt.Errorf("error marshalling database data: %w", err)
    }
    return os.WriteFile(db.path, jsonData, 0644)
} 
