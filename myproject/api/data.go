package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
	_ "github.com/go-sql-driver/mysql"
)

// Data represents the data structure
type Data struct {
	Message string `json:"message"`
}

// InitHandlers initializes API handlers
func InitHandlers() {
	http.HandleFunc("/api/data", GetData)
	http.HandleFunc("/api/users", GetUserData)
	http.HandleFunc("/api/users/create", CreateUser)
	http.HandleFunc("/api/users/update", UpdateUser)
}

type UserData struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// GetData is the handler function for the /api/data endpoint
func GetData(w http.ResponseWriter, r *http.Request) {
	// Create a Data object
	data := Data{
		Message: "Hello, World!",
	}

	// Convert the data object to JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Set the content type header
	w.Header().Set("Content-Type", "application/json")

	// Write the JSON response
	w.Write(jsonData)
}

func GetDataFromDB() ([]UserData, error) {

	db, err := sql.Open("mysql", "root:password@tcp(127.0.0.1:3306)/go_learn")
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT id,name,email FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []UserData
	for rows.Next() {
		var user UserData
		if err := rows.Scan(&user.ID, &user.Name, &user.Email); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

// ----------------------- Get data from DB ------------------------------------
func GetUserData(w http.ResponseWriter, r *http.Request) {
	users, err := GetDataFromDB()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(users)
}

// ----------------------- END Get data from DB ------------------------------------

// ----------------------- Create users ------------------------------------
type User struct {
	Name  string `json:"name" validate:"required"`
	Email string `json:"email" validate:"required,email"`
}

var validate = validator.New()

func CreateUser(w http.ResponseWriter, r *http.Request) {

	// Parse the JSON request body into a User struct
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := validate.Struct(user); err != nil {
		// Handle validation errors
		errors := err.(validator.ValidationErrors)
		http.Error(w, errors.Error(), http.StatusBadRequest)
		return
	}

	// Insert the user into the database
	db, err := sql.Open("mysql", "root:password@tcp(127.0.0.1:3306)/go_learn")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer db.Close()

	_, err = db.Exec("INSERT INTO users (name, email) VALUES (?, ?)", user.Name, user.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Respond with success message
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User created successfully"})
}

// ----------------------- end Create users ------------------------------------

// ----------------------- Update users ------------------------------------
func UpdateUser(w http.ResponseWriter, r *http.Request) {

	// Parse the JSON request body into a User struct
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := validate.Struct(user); err != nil {
		// Handle validation errors
		errors := err.(validator.ValidationErrors)
		http.Error(w, errors.Error(), http.StatusBadRequest)
		return
	}

	if err := validate.Struct(user); err != nil {
		// Handle validation errors
		errors := err.(validator.ValidationErrors)
		http.Error(w, errors.Error(), http.StatusBadRequest)
		return
	}

	// Insert the user into the database
	db, err := sql.Open("mysql", "root:password@tcp(127.0.0.1:3306)/go_learn")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer db.Close()
	id := 1
	newName := user.Name
	newEmail := user.Email

	// Prepare the SQL statement
	stmt, err := db.Prepare("UPDATE users SET name = ?, email= ? WHERE id = ?")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	// Execute the SQL statement
	result, err := stmt.Exec(newName, newEmail, id)
	if err != nil {
		log.Fatal(err)
	}

	// Check the number of rows affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": fmt.Sprintf("User updated successfully. Affected row: %d", rowsAffected)})

}

// ----------------------- end Update users ------------------------------------
