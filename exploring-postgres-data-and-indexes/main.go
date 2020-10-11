package main

import (
	"exploring-postgres-data-and-indexes/utils"
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	db := utils.DB{}
	db.Connect()
	defer db.Close()
	
	indexPath, err := db.GetUserNameIndexPath()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(indexPath)
	
	indexFile, err := ioutil.ReadFile(indexPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(indexFile)
}
