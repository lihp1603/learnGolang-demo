package main

import (
	"database/sql"
	"fmt"
	//  _操作其实只是引入该包。当导入一个包时，它所有的init()函数就会被执行，但有些时候并非真的需要使用这些包，仅仅是希望它的init()函数被执 行而已
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := sql.Open("mysql", "root:1234@tcp(127.0.0.1:3306)/lhp_learndemo")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer db.Close()
	// 测试网络连接是否可用
	if err = db.Ping(); err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println("connect mysql ok")

	var create_id int = 1903
	rows, err := db.Query("select req_no from ocv_creative where creative_id=?;", create_id)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		var req_no string
		err = rows.Scan(&req_no)
		if err != nil {
			fmt.Println(err.Error())
		}
		fmt.Println(req_no)
	}
	err = rows.Err()
	if err != nil {
		fmt.Println(err.Error())
	}
}
