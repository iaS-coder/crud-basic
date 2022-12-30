package server

import (
	"crud/mysql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type user struct {
	ID    uint32 `json:"id"`
	Nome  string `json:"nome"`
	Email string `json:"email"`
}

// CreateUser new user in data base
func CreateUsers(w http.ResponseWriter, r *http.Request) {
	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.Write([]byte("Failed to write request body!"))
		return
	}

	var user user

	if err = json.Unmarshal(requestBody, &user); err != nil {
		w.Write([]byte("Error on convertion user to struct"))
		return
	}

	db, err := mysql.Connect()
	if err != nil {
		w.Write([]byte("Error to conect in data base!"))
		return
	}
	defer db.Close()

	statement, err := db.Prepare("insert into usuarios (nome, email) values (?, ?)")
	if err != nil {
		w.Write([]byte("Error to crate statement!"))
		return
	}
	defer statement.Close()

	insert, err := statement.Exec(user.Nome, user.Email)
	if err != nil {
		w.Write([]byte("Error to execute statement!"))
		return
	}

	idInsert, err := insert.LastInsertId()
	if err != nil {
		w.Write([]byte("Error to get inserted id!"))
		return
	}

	// STATUS CODES
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf("User was inserted with success! Id: %d", idInsert)))
}

// GetUsers return all users saved on data base
func GetUsers(w http.ResponseWriter, r *http.Request) {
	db, err := mysql.Connect()
	if err != nil {
		w.Write([]byte("Error to conect in data base!"))
		return
	}
	defer db.Close()

	lines, err := db.Query("select * from usuarios")
	if err != nil {
		w.Write([]byte("Error to get users data!"))
		return
	}
	defer lines.Close()

	var users []user
	for lines.Next() {
		var user user

		if err := lines.Scan(&user.ID, &user.Nome, &user.Email); err != nil {
			w.Write([]byte("Error to scan users data!"))
			return
		}

		users = append(users, user)
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(users); err != nil {
		w.Write([]byte("Error to convert users to json format!"))
		return
	}

}

// GetUsers return specific users saved on data base
func GetUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	ID, err := strconv.ParseUint(params["id"], 10, 32)
	if err != nil {
		w.Write([]byte("Error to convert param ID to int"))
		return
	}

	db, err := mysql.Connect()
	if err != nil {
		w.Write([]byte("Error to conect in data base!"))
		return
	}
	defer db.Close()

	line, err := db.Query("select * from usuarios where id = ?", ID)
	if err != nil {
		w.Write([]byte("Error to get user search by ID in data base!"))
		return
	}

	var user user
	if line.Next() {
		if err := line.Scan(&user.ID, &user.Nome, &user.Email); err != nil {
			w.Write([]byte("Error to scan users data search by ID!"))
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(user); err != nil {
		w.Write([]byte("Error to convert user to json format!"))
		return
	}
}

// UpdateUsers make change/update of data for user
func UpdateUsers(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	ID, err := strconv.ParseUint(params["id"], 10, 32)
	if err != nil {
		w.Write([]byte("Error to convert param ID to int"))
		return
	}

	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.Write([]byte("Failed to write request body!"))
		return
	}

	var user user

	if err = json.Unmarshal(requestBody, &user); err != nil {
		w.Write([]byte("Error on convertion user to struct"))
		return
	}

	db, err := mysql.Connect()
	if err != nil {
		w.Write([]byte("Error to conect in data base!"))
		return
	}
	defer db.Close()

	statement, err := db.Prepare("update usuarios set nome = ?, email = ? where id = ?")
	if err != nil {
		w.Write([]byte("Error to create statement!"))
		return
	}
	defer statement.Close()

	if _, err := statement.Exec(user.Nome, user.Email, ID); err != nil {
		w.Write([]byte("Error to update user data!"))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// DeleteUsers remove some user by ID from data base
func DeleteUsers(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	ID, err := strconv.ParseUint(params["id"], 10, 32)
	if err != nil {
		w.Write([]byte("Error to convert param ID to int"))
		return
	}

	db, err := mysql.Connect()
	if err != nil {
		w.Write([]byte("Error to conect in data base!"))
		return
	}
	defer db.Close()

	statement, err := db.Prepare("delete from usuarios where id = ?")
	if err != nil {
		w.Write([]byte("Error to create statement!"))
		return
	}
	defer statement.Close()

	if _, err := statement.Exec(ID); err != nil {
		w.Write([]byte("Error to delete user data!"))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
