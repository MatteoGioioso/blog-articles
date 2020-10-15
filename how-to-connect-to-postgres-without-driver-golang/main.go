package main

import (
	"encoding/binary"
	"fmt"
	"github.com/jackc/pgio"
	"io"
	"net"
	"os"
	"pg-go/utils"
)

func makeStartupMessageRaw() []byte {
	buff := make([]byte, 0, 1024)

	// Allocate space for the length which will be calculated at the end of the encoding
	buff = append(buff, 0, 0, 0, 0)

	// Attach protocol version translated as uint 32 (3.0)
	buff = pgio.AppendUint32(buff, 196608)

	// Attach params, each key and value are separated by a 0 byte
	buff = append(buff, "user"...)
	buff = append(buff, 0)
	buff = append(buff, utils.User...)
	buff = append(buff, 0)
	buff = append(buff, "database"...)
	buff = append(buff, 0)
	buff = append(buff, utils.Database...)
	buff = append(buff, 0)
	buff = append(buff, "client_encoding"...)
	buff = append(buff, 0)
	buff = append(buff, "UTF8"...)
	buff = append(buff, 0)
	buff = append(buff, 0)

	// Calculate and append at the beginning of the buffer the total length of the message
	lengthOfTheMessage := int32(len(buff[0:]))
	binary.BigEndian.PutUint32(buff[0:], uint32(lengthOfTheMessage))

	return buff
}

func decodeStartupResponse(buff []byte) []byte {
	// first byte is the identifier char in this case R
	index := 0
	identifierChar := utils.GetASCIIIdentifier(buff, &index)
	fmt.Println("Id char: ", identifierChar)

	// the second part is a 4 byte which represent the length of your message
	length := utils.GetUint32Value(buff, &index)
	fmt.Println("length:", length)

	// this part is again a 4 byte integer which represent the auth method
	// in this case is going to be md5
	authMethod := utils.GetUint32Value(buff, &index)
	fmt.Println("auth method: ", authMethod)

	// this part is a 4 byte salt to encrypt the postgres credentials
	salt := buff[index : index+4]
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
	lengthOfMessage := int32(utils.Int32ByteLen + len(digestedPassword) + 1)
	buff = pgio.AppendInt32(buff, lengthOfMessage)

	// Attach encrypted Password
	buff = append(buff, digestedPassword...)
	buff = append(buff, 0)

	return buff
}

// In the normal case the backend will send some
// ParameterStatus messages, BackendKeyData, and finally ReadyForQuery.
func decodeAuthMessage(buff []byte) uint32 {
	index := 0
	// first byte is the identifier char in this case R
	identifierChar := utils.GetASCIIIdentifier(buff, &index)
	fmt.Println("Id char: ", identifierChar)

	// the second part is a 4 bytes which represent the length of your message
	length := utils.GetUint32Value(buff, &index)
	fmt.Println("length:", length)

	// last byte, if 0 means success
	authResult := utils.GetUint32Value(buff, &index)
	fmt.Println("auth result: ", authResult)

	// This loop iterate over all the parameters in the message
	// there are quite many of them (11) so I have decided to extract them programmatically
	for {
		identifierChar := utils.GetASCIIIdentifier(buff, &index)
		if identifierChar != "S" {
			break
		}

		length := utils.GetUint32Value(buff, &index)
		parameterStatusLength := length - 4 // subtract the length
		param := utils.GetStringValue(buff, parameterStatusLength, &index)
		fmt.Println("param: ", param)
	}
	// We need to go back to the ASCII identifier
	index = index - 1
	// This message provides secret-key data that the frontend must
	// save if it wants to be able to issue cancel requests later.
	// The frontend should not respond to this message, but should continue listening for a ReadyForQuery message.
	fmt.Println("BackendKeyDate Id:", utils.GetASCIIIdentifier(buff, &index))
	fmt.Println("length: ", utils.GetUint32Value(buff, &index))
	fmt.Println("process Id: ", utils.GetUint32Value(buff, &index))
	fmt.Println("backend key data: ", utils.GetUint32Value(buff, &index))

	// Start-up is completed. The frontend can now issue commands.
	readyForQueryId := utils.GetASCIIIdentifier(buff, &index)
	fmt.Println("Ready for query ASCII Id: ", readyForQueryId)
	rfqLength := utils.GetUint32Value(buff, &index)
	fmt.Println("rfqLength: ", rfqLength)
	status := utils.GetASCIIIdentifier(buff, &index)
	fmt.Println("Status (I: idle): ", status)

	return authResult
}

func makeQueryMessage() []byte {
	buff := make([]byte, 0, 1024)

	// ASCII identifier
	buff = append(buff, 'Q')
	query := "SELECT generate_series(1,10) AS id, md5(random()::text) AS descr;"

	lengthOfTheMessage := int32(utils.Int32ByteLen + len(query) + 1)
	buff = pgio.AppendInt32(buff, lengthOfTheMessage)

	buff = append(buff, query...)

	buff = append(buff, 0)

	return buff
}

func getQueryResponse(buff []byte) {
	types := make([]string, 0)
	names := make([]string, 0)
	tablePrinter := utils.TablePrinter{}
	fmt.Println(string(buff))
	index := 0
	ASCIIId := utils.GetASCIIIdentifier(buff, &index)
	fmt.Println("query result id: ", ASCIIId)
	length := utils.GetUint32Value(buff, &index)
	fmt.Println("length of message: ", length)

	numOfFields := utils.GetUint16Value(buff, &index)

	count := numOfFields
	// Decode table header (column name, type, ...)
	for {
		columnName := utils.GetStringValueWithNullTermination(buff[index:], &index)
		names = append(names, columnName)

		// We skip tableOid and column number, not needed for this demonstration
		index = index + utils.Int32ByteLen + utils.Int16ByteLen

		// Type in postgres have an OID you can run this query to check to which type correspond
		// SELECT oid,typname FROM pg_type WHERE oid='<found oid>'
		// For simplicity we just map few types to avoid making a new call
		typeOid := utils.GetUint32Value(buff, &index)
		typeName := utils.GetTypeName(typeOid)
		types = append(types, typeName)

		// We skip typeLength, typeMod, and format (text of binary)
		// not need those values for this purpose but you can see the specifications
		index = index + utils.Int16ByteLen + utils.Int32ByteLen + utils.Int16ByteLen

		count--
		if count == 0 {
			break
		}
	}

	// Decode each row
	rowCount := 0
	var rows [][]string
	for {
		rowCount++
		// I am quite sure there is a better way to do this
		// but basically we stop when the next byte has no more DataRow identifiers "D"
		if string(buff[index]) != "D" {
			break
		}

		index = index + 1 // jump the identifier
		_ = utils.GetUint32Value(buff, &index)
		numOfFields := utils.GetUint16Value(buff, &index)
		
		// Decode each column in each row
		count := numOfFields
		var row []string
		for {
			fieldLength := utils.GetUint32Value(buff, &index)
			// Finally the actual column data for the current row
			data := utils.GetStringValue(buff, fieldLength, &index)
			row = append(row, data)
			
			count--
			if count <= 0 {
				// Append the completed row into the rows array
				rows = append(rows, row)
				break
			}
		}
	}
	
	// psql style table printer
	tablePrinter.PrintTable(names, types, rows)
	
	commandComplete := utils.GetASCIIIdentifier(buff, &index)
	fmt.Println("Command complete: ", commandComplete)
	commandCompleteLength := utils.GetUint32Value(buff, &index)
	fmt.Println("length: ", commandCompleteLength)

	value := utils.GetStringValue(buff, commandCompleteLength-4, &index)
	fmt.Println(value)
	fmt.Println(rowCount)
}

func makeCloseMessage() []byte {
	buff := make([]byte, 0, 5)
	buff = append(buff, "X"...)
	buff = pgio.AppendInt32(buff, int32(utils.Int32ByteLen))
	return buff
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
	queryReply, err := utils.ReadAllBuffer(conn, queryMessage)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	getQueryResponse(queryReply)

	closeMessage := makeCloseMessage()
	if _, err := utils.Execute(conn, closeMessage); err != nil {
		if err == io.EOF {
			fmt.Println("Connection terminated")
			os.Exit(0)
		}
		fmt.Println(err)
		os.Exit(1)
	}
}
