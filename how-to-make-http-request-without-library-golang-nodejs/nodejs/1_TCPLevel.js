const net = require('net');

const port = 9000;
const host = "127.0.0.1"

const socket = new net.Socket()
socket.connect(port, host, (e) => {
    console.log("EVENT: ", e)
})

socket.on("connect", (e) => {
    console.log("CONNECTED")
    const req = `GET / HTTP/1.1\r\nHost: localhost:9000\r\nConnection: close\r\n\r\n`
    const buff = Buffer.from(req, "utf-8")
    socket.write(buff, (err) => {
        console.log("ERROR?:", err)
    })

    // const data;

    socket.on("data", chunk => {
        parseHTTP(chunk)
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
    const httpMessage = buff.toString().split("\r\n")
    const statusLine = httpMessage.shift();
    const statusLineComponents = statusLine.split(" ");
    console.log(statusLineComponents)
    for (const line of httpMessage) {
        httpMessage.shift()
        console.log(httpMessage)
        console.log(line === "")
        if (line === ""){
                break
        }
    }

    const chunkSize = httpMessage.shift();
    console.log(chunkSize)
    const body = httpMessage.shift();
    console.log(body)
}
