package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
)
var db *sql.DB
/**
初始化连接数据库
*/
func init() {
	var err error
	db, err  = sql.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/blog")
	if err != nil {
		fmt.Errorf("Can't open database"+err.Error())
	}
}

/**
获取用户ID
*/
func getUserId(db *sql.DB, arg string) (int, error){
	var id int
	err := db.QueryRow("SELECT id FROM users WHERE email = ?",arg).Scan(&id)
	if err != sql.ErrNoRows {
		return id, errors.Wrap(err,"The user was not found")
	}
	return id, nil
}

func main() {
	defer db.Close()
	id , err1 := getUserId(db, "875198125@qq.com")

	if err1 != nil {
		fmt.Printf("Stack Trace:\n%+v\n", err1)
	}
	fmt.Println(id)
}