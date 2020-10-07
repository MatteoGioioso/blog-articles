package main

import (
	"fmt"
	"io/ioutil"
)

// SELECT oid FROM pg_databases WHERE databse='test';

// SELECT oid,relname FROM pg_class WHERE relname='users';

func main() {
	path := "pg_data/base/16384/16399"
	file, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println(err)
	}
	
	fmt.Println(file)
}
