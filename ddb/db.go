package ddb

import (
	"database/sql"
	"downloader/log"

	_ "github.com/mattn/go-sqlite3"
)

type TaskBean struct {
	Name            string
	Url             string
	Create_time     string
	Completion_time string
}

var l = log.New()

func CreateDB() {
	db, err := openDB()
	if err != nil {
		l.PrintMulti(err)
		return
	}
	db.Exec("create database if exist ")
}

func openDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "dbfile.db")
	if err != nil {
		return nil, err
	}
	return db, nil
}

func InsertDB(t TaskBean) (int64, error) {
	db, err := openDB()
	if err != nil {
		l.PrintMulti(err)
		return -1, err
	}
	stmt, err := db.Prepare("insert into tb_task(name, url, create_time) values(?, ?, ?)")
	if err != nil {
		return -1, err
	}
	defer stmt.Close()
	defer db.Close()
	result, err := stmt.Exec(t.Name, t.Url, t.Create_time)
	if err != nil {
		return -1, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return -1, err
	}
	return id, nil
}
