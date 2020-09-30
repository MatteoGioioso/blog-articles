package utils

import (
	"crypto/md5"
	"encoding/binary"
	"encoding/hex"
	"io"
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

// Helper function for ...
func HexMD5(s string) string {
	hash := md5.New()
	io.WriteString(hash, s)
	return hex.EncodeToString(hash.Sum(nil))
}


// Very raw type conversion, in reality there are a lot more types,
// but for this scope this is going to be enough
func GetGoType(typeOid uint32) string {
	switch typeOid {
	case 23:
		return "int"
	case 1043:
		return "string"
	}

	return ""
}

func GetColumnName(buff []byte) string {
	count := 0

	for {
		b := buff[count]
		// column names are delimited by a 0 byte
		if b == 0 {
			break
		}

		// we count the number of bytes until we encounter the 0 byte
		// then we stop and return the casted slice containing the string
		count++
	}

	return string(buff[:count])
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
