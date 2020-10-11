package utils

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"os"
	"strconv"
)

type DB struct {
	conn *pgx.Conn
}

func (db *DB) Connect() {
	connStr := "postgres://postgres:123@localhost:5432/test?sslmode=disable"
	conn, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		fmt.Println("Unable to Connect to database", err)
		os.Exit(1)
	}

	db.conn = conn
}

func (db DB) getDbOID() (oid string, err error) {
	rows, err := db.conn.Query(context.Background(), "SELECT oid FROM pg_database WHERE datname='test'")
	if err != nil {
		return "", err
	}
	defer rows.Close()

	for rows.Next() {
		var rawOid uint32
		if err := rows.Scan(&rawOid); err != nil {
			return "", err
		}

		oid = strconv.Itoa(int(rawOid))
	}

	return
}

func (db DB) getTableOid() (oid string, err error) {
	rows, err := db.conn.Query(context.Background(), "SELECT oid FROM pg_class WHERE relname='users'")
	if err != nil {
		return "", err
	}
	defer rows.Close()

	for rows.Next() {
		var rawOid uint32
		if err := rows.Scan(&rawOid); err != nil {
			return "", err
		}

		oid = strconv.Itoa(int(rawOid))
	}

	return
}

func (db DB) getUsernameIndexOid() (oid string, err error) {
	rows, err := db.conn.Query(
		context.Background(),
		`
			SELECT i.oid
			FROM   pg_index as idx
				   JOIN   pg_class as i
						  ON     i.oid = idx.indexrelid
				   JOIN   pg_am as am
						  ON     i.relam = am.oid
			WHERE i.relname='users_username_key';`,
	)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	for rows.Next() {
		var rawOid uint32
		if err := rows.Scan(&rawOid); err != nil {
			return "", err
		}

		oid = strconv.Itoa(int(rawOid))
	}

	return
}

func (db *DB) GetUsersTableFilePath() (string, error) {
	path := "pg_data/base"
	dbOid, err := db.getDbOID()
	if err != nil {
		return "", err
	}
	path = path + "/" + dbOid

	tableOid, err := db.getTableOid()
	if err != nil {
		return "", err
	}
	path = path + "/" + tableOid

	_, err = db.getUsernameIndexOid()
	if err != nil {
		return "", err
	}

	return path, nil
}

func (db DB) GetUserNameIndexPath() (string, error) {
	path := "pg_data/base"
	dbOid, err := db.getDbOID()
	if err != nil {
		return "", err
	}
	path = path + "/" + dbOid
	
	usernameIndexOid, err := db.getUsernameIndexOid()
	if err != nil {
		return "", err
	}
	
	return path + "/" + usernameIndexOid, nil
}

func (db *DB) Close() {
	db.conn.Close(context.Background())
}
