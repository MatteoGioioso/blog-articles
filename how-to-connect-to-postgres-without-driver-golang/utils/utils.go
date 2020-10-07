package utils

import (
	"crypto/md5"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"net"
)

var (
	Address      = "127.0.0.1:5432"
	User         = "postgres"
	Database     = "test"
	Password     = "123"
	Int32ByteLen = 4
	Int16ByteLen = 2
)

func HasParameter(buff []byte, index int) bool {
	ASCIIId := string(buff[index])

	// "S" stands for parameter
	if ASCIIId == "S" {
		return true
	}

	return false
}

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
func GetGoType(typeOid uint32) string {
	switch typeOid {
	case 25: // text
		return "string"
	case 23:
		return "int"
	case 1043: // varchar
		return "string"
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
