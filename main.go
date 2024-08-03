package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
)

type User struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Phone string `json:"phone"`
	Age   int    `json:"age"`
}

var db *sql.DB

func connectDb() {
	var err error
	dsn := "root:admin@tcp(localhost:3306)/practice"
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Error opening database: ", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal("Error pinging database: ", err)
	}

	fmt.Println("Connected to the database")
}

func saveUser(c echo.Context) error {
	user := new(User)
	if err := c.Bind(user); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	query := "INSERT INTO users(name, phone, age) VALUES (?, ?, ?)"
	_, err := db.Exec(query, user.Name, user.Phone, user.Age)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to save user"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "User saved successfully"})
}



func getUser(c echo.Context) error {
	query := "SELECT id, name, phone, age FROM users"
	rows, err := db.Query(query)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get users"})
	}
	defer rows.Close()

	users := []User{}
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.Id, &user.Name, &user.Phone, &user.Age); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to scan"})
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "error in itteration"})
	}

	return c.JSON(http.StatusOK, users)
}

func main() {
	fmt.Println("Connecting to the database")
	connectDb()

	e := echo.New()

	e.POST("/create-user", saveUser)
	e.GET("/all-user", getUser)

	
	e.Logger.Fatal(e.Start(":8080"))
}
