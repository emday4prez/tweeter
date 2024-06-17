package database

import (
	"encoding/json"
	"errors"
	"os"
	"sync"
)

type DB struct {
	path string
	mu   *sync.RWMutex
}

type Chirp struct {
	ID   int    `json:"id"`
	Body string `json:"body"`
}

type DBStructure struct {
	Chirps       map[int]Chirp        `json:"chirps"`
	Users        map[int]User         `json:"users"`
	RefreshTokens map[string]RefreshToken `json:"refresh_tokens"` 
}

type User struct {
	ID           int      `json:"id"`
	Email        string   `json:"email"`
	Password     string   `json:"password"`
}

type RefreshToken struct {
	Token  string `json:"token"`
	UserID int    `json:"user_id"`
	Expiry string `json:"expiry"`
}

func NewDB(path string) (*DB, error) {
	db := &DB{
		path: path,
		mu:   &sync.RWMutex{},
	}
	err := db.ensureDB()
	return db, err
}

func (db *DB) CreateChirp(body string) (Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}

	id := len(dbStructure.Chirps) + 1
	chirp := Chirp{
		ID:   id,
		Body: body,
	}
	dbStructure.Chirps[id] = chirp

	err = db.writeDB(dbStructure)
	if err != nil {
		return Chirp{}, err
	}

	return chirp, nil
}

func (db *DB) GetChirps() ([]Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return nil, err
	}

	chirps := make([]Chirp, 0, len(dbStructure.Chirps))
	for _, chirp := range dbStructure.Chirps {
		chirps = append(chirps, chirp)
	}

	return chirps, nil
}

func (db *DB) GetChirpById(id int) (Chirp, error) {
    var matchingChirp Chirp
	dbStructure, err := db.loadDB()
	if err != nil {
		return matchingChirp, err
	}

 
	for _, chirp := range dbStructure.Chirps {
        if chirp.ID == id {
            return chirp, nil
        }
	 
	}

	return matchingChirp, errors.New("no chirps with that id")
}

func (db *DB) GetUserByEmail(email string) (User, error) {
    var matchingUser User
	dbStructure, err := db.loadDB()
	if err != nil {
		return matchingUser, err
	}

 
	for _, user := range dbStructure.Users {
        if user.Email == email {
            return user, nil
        }
	 
	}

	return matchingUser, errors.New("no users with that email")
}

func (db *DB) FindToken(token string)(RefreshToken, error){
	    var matchingToken RefreshToken
	dbStructure, err := db.loadDB()
	if err != nil {
		return matchingToken, err
	}

 
    if foundToken, exists := dbStructure.RefreshTokens[token]; exists {
        return foundToken, nil
    }

	return matchingToken, errors.New("token not found")
}

func (db *DB) CreateUser(email string, password string) (User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	id := len(dbStructure.Users) + 1
	user := User{
		ID:   id,
		Email: email,
		Password: password,
	}
	dbStructure.Users[id] = user

	err = db.writeDB(dbStructure)
	if err != nil {
		return User{}, err
	}

	return user, nil
}


func (db *DB) UpdateUser(id int, email string, password string) error {
    dbStructure, err := db.loadDB()
    if err != nil {
        return err
    }

    if user, exists := dbStructure.Users[id]; exists { 
        user.Email = email
        user.Password = password
        dbStructure.Users[id] = user // Reassign the modified user
    } else {
        return errors.New("user not found")
    }

    return db.writeDB(dbStructure)
}

func (db *DB) createDB() error {
	dbStructure := DBStructure{
		Chirps: map[int]Chirp{},
        Users: map[int]User{},
	}
	return db.writeDB(dbStructure)
}

func (db *DB) ensureDB() error {
	_, err := os.ReadFile(db.path)
	if errors.Is(err, os.ErrNotExist) {
		return db.createDB()
	}
	return err
}

func (db *DB) loadDB() (DBStructure, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	dbStructure := DBStructure{}
	dat, err := os.ReadFile(db.path)
	if errors.Is(err, os.ErrNotExist) {
		return dbStructure, err
	}
	err = json.Unmarshal(dat, &dbStructure)
	if err != nil {
		return dbStructure, err
	}

	return dbStructure, nil
}

func (db *DB) writeDB(dbStructure DBStructure) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	dat, err := json.Marshal(dbStructure)
	if err != nil {
		return err
	}

	err = os.WriteFile(db.path, dat, 0600)
	if err != nil {
		return err
	}
	return nil
}