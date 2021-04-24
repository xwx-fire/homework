package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
)
var db *sql.DB

func init() {
	var err error
	db, err  = sql.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/blog")
	if err != nil {
		fmt.Errorf("Can't open database"+err.Error())
	}
}

func getUser(db *sql.DB, arg string) (int, error){
	var id int
	err := db.QueryRow("SELECT id FROM users WHERE email = ?",arg).Scan(&id)
	return id, errors.Wrap(err,"The user was not found")
}

func main() {
	defer db.Close()
	id , err1 := getUser(db, "875198125@qq.com")
	if !errors.Is(err1, sql.ErrNoRows) {
		fmt.Printf("Stack Trace:\n%+v\n", err1)
	}
	fmt.Println(id)
}