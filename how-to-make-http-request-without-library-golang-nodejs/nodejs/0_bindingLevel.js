const net = require('./bindingsUtils/socket')

const port = 9000;
const host = "127.0.0.1"

const socket = new net.Socket()
socket.connect(port, host, () => {
    console.log("on connect")
})

socket.on("connect", () => {
    console.log("CONNECTED")
    const req = `GET / HTTP/1.1\r\nHost: localhost:9000\r\nConnection: close\r\n\r\n`
    const buff = Buffer.from(req, "utf-8")
    socket.write(buff, (err) => {
        console.log("ERROR?:", err)
    })

    socket.on("data", chunk => {
        const response = parseHTTP(chunk);
        console.log(response)
    })

    socket.on("end", (e) => {
        console.log("END: ", e)
        socket.end(() => {
            console.log("CONNECTION ENDED")
        })
    })
})

// HTTP-Version SP Status-Code SP Reason-Phrase CRL
// *(( general-header        ; Section 4.5
// | response-header        ; Section 6.2
// | entity-header ) CRLF)  ; Section 7.1
// CRLF
// [ message-body ]          ; Section 7.2
function parseHTTP(buff) {
    const response = {}
    const httpMessage = buff.toString().split("\r\n")
    const statusLine = httpMessage.shift();
    const statusLineComponents = statusLine.split(" ");
    response.statusCode = statusLineComponents[1]

    for (const _ of httpMessage) {
        const line = httpMessage.shift();

        // CRLF before the body
        if (line === ""){
            break
        }
        const headerKeyValue = line.split(" ")
        const key = headerKeyValue[0].replace(":", "")
        const value = headerKeyValue[1]
        response[key] = value
    }

    const chunkSize = httpMessage.shift();
    response.size = parseInt(chunkSize, 16)
    const body = httpMessage.shift();
    response.body = JSON.parse(body)
    return response
}
