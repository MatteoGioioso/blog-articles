package main

import (
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"os"
)

const localInstance = "postgres://postgres:123@localhost:5432/test?sslmode=disable"

func main() {
	m, err := migrate.New(
		"file://files",
		localInstance,
	)

	if err != nil {
		fmt.Println("Connection error: ", err)
		os.Exit(1)
	}

	if err := m.Up(); err != nil {
		fmt.Println("Migration error:  ", err)
		os.Exit(1)
	}
}
