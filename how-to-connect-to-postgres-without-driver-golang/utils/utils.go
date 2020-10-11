package utils

import (
	"crypto/md5"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"net"
	"strings"
)

var (
	Address      = "127.0.0.1:5432"
	User         = "postgres"
	Database     = "test"
	Password     = "123"
	Int32ByteLen = 4
	Int16ByteLen = 2
)

// Helper function for hashing our credentials
func HexMD5(s string) string {
	hash := md5.Sum([]byte(s))
	return hex.EncodeToString(hash[:])
}

// Very raw type conversion, in reality there are a lot more types,
// but for this scope this is going to be enough.
// PostgreSQL has a type table called pg_type (SELECT oid, typname FROM pg_type).
// The driver uses the knowledge about OIDs to figure out how to map data from database column types into primitive Go types.
// For this purpose, pgx internally uses the following map (key — type name, value — Object ID)
func GetTypeName(typeOid uint32) string {
	switch typeOid {
	case 25: // text
		return "text"
	case 23:
		return "int4"
	case 1043: // varchar
		return "varchar"
	}

	return ""
}

// Helper function to write on a connection and receive response data from it
func Execute(conn net.Conn, message []byte) ([]byte, error) {
	if _, err := conn.Write(message); err != nil {
		return nil, err
	}

	reply := make([]byte, 1024)

	if _, err := conn.Read(reply); err != nil {
		return nil, err
	}

	return reply, nil
}

func ReadAllBuffer(conn net.Conn, message []byte) ([]byte, error) {
	buf := make([]byte, 0, 4096)
	tmp := make([]byte, 4096)

	if _, err := conn.Write(message); err != nil {
		return nil, err
	}

	for {
		n, err := conn.Read(tmp)
		if err != nil {
			return nil, err
		}

		buf = append(buf, tmp[:n]...)
		fmt.Println(n)

		// Last message is Query idling "I"
		// Did not find any better method to stop reading from the buffer
		// till is completed
		if string(tmp[n-1]) == "I" {
			return buf, nil
		}
	}
}

func GetUint32Value(buff []byte, index *int) (value uint32) {
	value = binary.BigEndian.Uint32(buff[*index : *index+Int32ByteLen])
	newIndex := *index + Int32ByteLen
	*index = newIndex
	return
}

func GetUint16Value(buff []byte, index *int) (value uint16) {
	value = binary.BigEndian.Uint16(buff[*index : *index+Int16ByteLen])
	newIndex := *index + Int16ByteLen
	*index = newIndex
	return
}

func GetStringValue(buff []byte, length uint32, index *int) (value string) {
	value = string(buff[*index : *index+int(length)])
	*index = *index + int(length)
	return value
}

func GetStringValueWithoutLenButWithDivider(buff []byte, index *int) (value string) {
	count := 0

	for {
		b := buff[count]
		// Field is delimited by a 0 byte
		if b == 0 {
			break
		}

		// we count the number of bytes until we encounter the 0 byte
		// then we stop and return the casted slice containing the string
		count++
	}

	*index = *index + count + 1 // add the 0 byte

	return string(buff[:count])
}

func GetASCIIIdentifier(buff []byte, index *int) (id string) {
	id = string(buff[*index])
	*index = *index + 1
	return id
}

type TablePrinter struct {
	headerColLen []int
	dataColLen   []int
	fieldCount   int
	names        []string
	types        []string
}

func (t *TablePrinter) getHeaderInfo(names []string, types []string) {
	header := ""
	count := 0
	for i, name := range names {
		row := fmt.Sprintf("%v (%v) ", name, types[i])
		header = header + row
		t.headerColLen = append(t.headerColLen, len(row))
		count++
	}

	t.fieldCount = count
	t.names = names
	t.types = types
}

func (t *TablePrinter) formatHeader() string {
	header := ""
	count := 0
	for i, name := range t.names {
		row := fmt.Sprintf("%v (%v)%v| ", name, t.types[i], strings.Repeat(" ", t.dataColLen[i]))
		header = header + row
		count++
	}

	return header
}

func (t *TablePrinter) formatTableData(rows [][]string) string {
	t.dataColLen = make([]int, t.fieldCount)
	tableText := ""
	for _, row := range rows {
		rowText := ""
		for j, data := range row {
			columnLength := t.headerColLen[j]
			// Check which between the header or the data is longer
			if columnLength <= len(data) {
				columnLength = len(data)
				rowText = rowText + fmt.Sprintf("%v | ", data)
				t.dataColLen[j] = columnLength - t.headerColLen[j] + 2 // two are the extra spaces
				continue
			}
			columnLength = columnLength - len(data)
			t.dataColLen[j] = 1
			rowText = rowText + fmt.Sprintf("%v%v| ", data, strings.Repeat(" ", columnLength))
		}

		tableText = tableText + rowText + "\n"
	}
	return tableText
}

func (t *TablePrinter) PrintTable(names []string, types []string, rows [][]string) {
	t.getHeaderInfo(names, types)
	tableText := t.formatTableData(rows)
	header := t.formatHeader()

	fmt.Println("")
	fmt.Println("")
	fmt.Println(header)
	fmt.Println(strings.Repeat("-", len(header)))
	fmt.Println(tableText)
	fmt.Println("")
	fmt.Println("")
}
