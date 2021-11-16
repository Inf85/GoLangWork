package db_connection

import ("database/sql"
	_ "github.com/go-sql-driver/mysql"
)


/**
  Set DataBase Connection
 */
func SetConnection() *sql.DB {
	db, err:=sql.Open("mysql",  "root@tcp(127.0.0.1:3306)/golang")
	if err != nil{
		panic(err)
		return nil
	}

	return db

}
