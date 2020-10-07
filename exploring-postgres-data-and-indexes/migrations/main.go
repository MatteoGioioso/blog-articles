package main

import (
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

const rdsInstance = "postgres://madeo:vba25tlm2@pa14sia00icfagv.cndlk5sduxcf.ap-southeast-1.rds.amazonaws.com:5432/pgtrigger"
const localInstance = "postgres://postgres:123@localhost:5432/test?sslmode=disable"

func main() {
	m, err := migrate.New(
		"file://files",
		rdsInstance)
	fmt.Println("Connection error: ", err)
	//err3 := m.Force(2)
	//fmt.Println(err3)

	err2 := m.Up()


	fmt.Println("Migration error:  ", err2)
}
