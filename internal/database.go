package db

import (
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
 

//check error using os.IsNotExist(err) for "file not found"
 
	//database.json does not exist
//create file
//create empty db structure (with the type and make a map of chirps for the Chirps key)
 
// use json marshall indent to fill with empty data
 


// if err := writefile; err != nil  {return error}  


 //database.json exists
 
}



// loadDB reads the database file into memory
func (db *DB) loadDB() (DBStructure, error){

}

// writeDB writes the database file to disk
func (db *DB) writeDB(dbStructure DBStructure) error{

} 





func (db *DB) eeDB()error {
	_, err := os.Stat(db.path)

	if os.IsNotExist(err){
		
	}
}