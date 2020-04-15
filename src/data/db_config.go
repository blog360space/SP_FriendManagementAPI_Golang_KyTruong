package data

import (
	"database/sql"
	"fmt"
)

func InitDB() (*sql.DB, error) {
	dbDriver := "mysql"
	dbUser := "root"
	dbPass := "123456@x@X"
	dbName := "friendMgmt"
	dbPort := "3306"
	dbHost := "fullstack-mysql"
	dbUrl := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", dbUser, dbPass, dbHost, dbPort, dbName)
	db, err := sql.Open(dbDriver, dbUrl)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
