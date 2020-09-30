package main

import (
	"encoding/binary"
	"fmt"
	"github.com/jackc/pgio"
	"net"
	"os"
	"pg-go/utils"
)

func makeStartupMessageRaw() []byte {
	buff := make([]byte, 0, 1024)

	// this is initially 0
	sp := len(buff)

	// Allocate space for the length which will be calculated at the end of the encoding
	// not 100% sure why is -1, it just creates 255, 255, 255
	buff = pgio.AppendInt32(buff, -1)

	//Attach protocol version
	buff = pgio.AppendUint32(buff, 196608)

	// Attach params
	buff = append(buff, "user"...)
	buff = append(buff, 0)
	buff = append(buff, utils.User...)
	buff = append(buff, 0)
	buff = append(buff, "database"...)
	buff = append(buff, 0)
	buff = append(buff, utils.Database...)
	buff = append(buff, 0)
	buff = append(buff, 0)

	// Append at the beginning of the buffer the total length of the message
	lengthOfTheMessage := int32(len(buff[sp:]))
	pgio.SetInt32(buff[sp:], lengthOfTheMessage)

	return buff
}

func decodeStartupResponse(buff []byte) []byte {
	// first byte is the identifier char in this case R
	identifierChar := string(buff[0])
	fmt.Println("Id char: ", identifierChar)

	// the second part is a 4 byte which represent the length of your message
	length := binary.BigEndian.Uint32(buff[1:5])
	fmt.Println("length:", length)

	// this part is again a 4 byte integer which represent the auth method
	// in this case is going to be md5
	authMethod := binary.BigEndian.Uint32(buff[5:9])
	fmt.Println("auth method: ", authMethod)

	// this part is a 4 byte salt to encrypt the postgres credentials
	salt := buff[9:13]
	fmt.Println("salt: ", salt)
	fmt.Println("----------------------")
	return salt
}

func makeAuthMessage(salt []byte) []byte {
	buff := make([]byte, 0, 1024)

	// ASCII identifier for authentication Password
	buff = append(buff, 'p')

	// Formula for postgres Password encryption
	digestedPassword := "md5" + utils.HexMD5(utils.HexMD5(utils.Password+utils.User)+string(salt))

	// total length of the message
	buff = pgio.AppendInt32(buff, int32(4+len(digestedPassword)+1))

	// Attach encrypted Password
	buff = append(buff, digestedPassword...)
	buff = append(buff, 0)

	return buff
}

// In the normal case the backend will send some
// ParameterStatus messages, BackendKeyData, and finally ReadyForQuery.
func decodeAuthMessage(buff []byte) uint32 {
	// first byte is the identifier char in this case R
	identifierChar := string(buff[0])
	fmt.Println("Id char: ", identifierChar)

	// the second part is a 4 bytes which represent the length of your message
	length := binary.BigEndian.Uint32(buff[1:5])
	fmt.Println("length:", length)

	// last byte, if 0 means success
	authResult := binary.BigEndian.Uint32(buff[5:9])
	fmt.Println("auth method: ", authResult)

	// This loop iterate over all the parameters in the message
	// there are quite many of them (11) so I have decided to extract them programmatically
	var index = 9
	for {
		hasParam := utils.HasParameter(buff, index)
		if !hasParam {
			break
		}

		index = index + 1
		parameterStatusLength := binary.BigEndian.Uint32(buff[index : index+4])
		fmt.Println("parameter status length: ", parameterStatusLength)

		nextIndex := index + int(parameterStatusLength)
		param := string(buff[index+1+4 : nextIndex])
		fmt.Println("params: ", param)

		index = nextIndex + 1
	}

	// This message provides secret-key data that the frontend must
	// save if it wants to be able to issue cancel requests later.
	// The frontend should not respond to this message, but should continue listening for a ReadyForQuery message.
	fmt.Println("BackendKeyDate Id:", string(buff[index]))
	fmt.Println("length: ", binary.BigEndian.Uint32(buff[index+1:index+5]))
	fmt.Println("process Id: ", binary.BigEndian.Uint32(buff[333:337]))
	fmt.Println("backend key data: ", binary.BigEndian.Uint32(buff[337:341]))

	// Start-up is completed. The frontend can now issue commands.
	fmt.Println("Ready for query ASCII Id: ", string(buff[341]))
	fmt.Println("length: ", binary.BigEndian.Uint32(buff[342:346]))
	fmt.Println("Status (I: idle): ", string(buff[346]))

	return authResult
}

func makeQueryMessage() []byte {
	buff := make([]byte, 0, 1024)

	// ASCII identifier
	buff = append(buff, 'Q')
	query := "SELECT * FROM users;"

	lengthOfTheMessage := int32(4 + len(query) + 1)
	buff = pgio.AppendInt32(buff, lengthOfTheMessage)

	buff = append(buff, query...)

	buff = append(buff, 0)

	return buff
}

type Field struct {
	Data string
	Type string
}
type Row []Field
type Rows []Row

func getQueryResponse(buff []byte) Rows {
	rows := make(Rows, 0)
	ASCIIId := string(buff[0])
	// This char should be "T"
	fmt.Println("query result id: ", ASCIIId)
	index := 1 // start after the first char
	length := utils.GetUint32Value(buff, &index)
	fmt.Println("length of message: ", length)

	numOfFields := utils.GetUint16Value(buff, &index)
	count := numOfFields
	// Decode table header (column name, type, ...)
	for {
		columnName := utils.GetColumnName(buff[index:])
		index = len(columnName) + index + 1 // At the end of each column name there is a 0 byte
		
		// We skip tableOid and column number, not needed for this demostration
		index = index + utils.Int32ByteLen + utils.Int16ByteLen
		typeOid := utils.GetUint32Value(buff, &index)
		goType := utils.GetGoType(typeOid)
		fmt.Println("type: ", goType)
		
		// We skip typeLength, typeMod, and format,
		// not need those values for this purpose but you can see the specifications
		index = index + utils.Int16ByteLen + utils.Int32ByteLen + utils.Int16ByteLen

		count--
		if count == 0 {
			break
		}
	}

	// Decode each row
	for {
		
		// I am quite sure there is a better way to do this
		// but basically we stop when the next byte has no more DataRow identifiers "D"
		if string(buff[index]) != "D" {
			break
		}

		index = index + 1 // jump the identifier
		length := utils.GetUint32Value(buff, &index)
		fmt.Println("length: ", length)
		numOfFields := utils.GetUint16Value(buff, &index)
		
		// Decode each column in each row
		count := numOfFields
		for {
			fieldLength := utils.GetUint32Value(buff, &index)
			
			// Finally the actual column data for the current row
			data := string(buff[index : index+int(fieldLength)])

			row := make(Row, 0)
			row = append(row, Field{Data: data})
			rows = append(rows, row)

			index = index + int(fieldLength)
			count--
			if count <= 0 {
				break
			}
		}
	}

	return rows
}

func main() {
	conn, err := net.Dial("tcp", utils.Address)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer conn.Close()

	message := makeStartupMessageRaw()
	reply, err := utils.Execute(conn, message)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	salt := decodeStartupResponse(reply)
	authMessage := makeAuthMessage(salt)

	authReply, err := utils.Execute(conn, authMessage)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	authResult := decodeAuthMessage(authReply)
	if authResult != 0 {
		fmt.Println("Authentication failed")
		os.Exit(1)
	}

	queryMessage := makeQueryMessage()
	queryReply, err := utils.Execute(conn, queryMessage)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	response := getQueryResponse(queryReply)
	fmt.Printf("%+v\n", response)
}
