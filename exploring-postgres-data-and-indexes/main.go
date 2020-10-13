package main

import (
	"exploring-postgres-data-and-indexes/utils"
	"fmt"
	"io/ioutil"
	"os"
)

type ItemData struct {
	lpOff uint
	lpFlags uint
	lpLen uint
}

func main() {
	db := utils.DB{}
	db.Connect()
	defer db.Close()
	
	path, err := db.GetUsersTableFilePath()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("path: ", path)
	tableFile, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	parser := utils.Parser{Buff: tableFile}
	//dump := hexdump.Config{Width: 32}
	//fmt.Println(dump.Dump(parser.Buff))
	// Table header
	parser.ShiftBy(12)
	pdLower := parser.GetInt16()
	pdUpper := parser.GetInt16()
	fmt.Println(pdUpper, pdLower)
	pdSpecial := parser.GetInt16()
	fmt.Println(pdSpecial)
	pdVersion := parser.GetInt16()
	fmt.Println(pdVersion)
	prune := parser.GetInt32()
	fmt.Println(prune)
	fmt.Printf("%08b\n", []byte{0xc8, 0x9f, 0x6c, 0x00})
	fmt.Println(parser.GetInt16())
	fmt.Println(parser.GetInt16())
	

	//offset1 := parser.GetInt16()
	//length1 := parser.GetInt16()
	//fmt.Println(offset1, length1)
	//
	//offset2 := parser.GetInt16()
	//length2 := parser.GetInt16()
	//fmt.Println(offset2, length2)
	//
	//offset3 := parser.GetInt16()
	//length3 := parser.GetInt16()
	//fmt.Println(offset3, length3)
}
