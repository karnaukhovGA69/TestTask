package main

import "database/sql"

type DB interface {
	Get()
	Put()
}

type PostgresDB struct {
	db *sql.DB
}

type DBelg struct {
	longToShort map[string]string
	shortToLong map[string]string
}

func (ps *PostgresDB) Get() {

}

func (ps *PostgresDB) Put() {

}

func (gleb *DBelg) Get() {

}

func (gleb *DBelg) Put() {

}

func MakeDB(name string) DB {
	switch name {
	case "Postgres":
		return &PostgresDB{db: &sql.DB{}}
	case "DBelg":
		return &DBelg{longToShort: make(map[string]string), shortToLong: make(map[string]string)}
	default:
		return nil
	}

}
