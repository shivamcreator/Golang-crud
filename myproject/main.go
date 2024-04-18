// main.go

package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

// import (
//     "database/sql"
//     "encoding/json"
//     "fmt"
//     "log"
//     "net/http"

//     "github.com/go-playground/validator/v10"
//     _ "github.com/go-sql-driver/mysql"
// )

type UserData struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}
type Data struct {
	Message string `json:"message"`
}

func main() {
	http.HandleFunc("/users", GetDataFromDB)
	fmt.Println("Server listening on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func GetDataFromDB() ([]UserData, error) {

	db, err := sql.Open("mysql", "root:password@tcp(127.0.0.1:3306)/my_database")
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

// func handleUsers(w http.ResponseWriter, r *http.Request) {
// 	// Create a Data object
// 	// data := Data{
// 	// 	Message: "Hello, World!",
// 	// }

// 	// // Convert the data object to JSON
// 	// jsonData, err := json.Marshal(data)
// 	// if err != nil {
// 	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
// 	// 	return
// 	// }

// 	// // Set the content type header
// 	// w.Header().Set("Content-Type", "application/json")

// 	// // Write the JSON response
// 	// w.Write(jsonData)
// 	// Connect to MySQL database
// 	db, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/my_database")
// 	// db, err := sql.Open("mysql", "root:password@tcp(127.0.0.1:3306)/go_learn")

// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	defer db.Close()

// 	// Query database for users
// 	rows, err := db.Query("SELECT id, username, email FROM users")
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	defer rows.Close()

// 	// Process query results
// 	var users []User
// 	for rows.Next() {
// 		var user User
// 		if err := rows.Scan(&user.ID, &user.Username, &user.Email); err != nil {
// 			http.Error(w, err.Error(), http.StatusInternalServerError)
// 			return
// 		}
// 		users = append(users, user)
// 	}

// 	// Send JSON response
// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(users)
// }
