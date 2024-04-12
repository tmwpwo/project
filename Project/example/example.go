package main

import (
	"database/sql"
	"fmt"
)

func main() {

	secret := "secretKey"

	db, err := sql.Open("mysql", "user:password@/dbname")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer db.Close()

	userInput := "1 OR 1=1"
	query := fmt.Sprintf("INSERT INTO users (name) VALUES ('%s')", userInput) // SQL injection vulnerability
	_, err = db.Exec(query)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("User added successfully.")
}
