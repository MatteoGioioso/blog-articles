const net = require('net');

const port = 9000;
const host = "127.0.0.1"

const socket = new net.Socket()
socket.connect(port, host, (e) => {
    console.log("EVENT: ", e)
})

socket.on("connect", (e) => {
    console.log("CONNECTED")
    const buff = Buffer.from("hello world", "utf-8")
    socket.write(buff, (err) => {
        console.log("ERROR?:", err)
    })
})

socket.on("data", data => {
    console.log("SOME DATA: ", data.toString())
})

socket.on("end", (e) => {
    console.log("END: ", e)
    socket.end(e => {
        console.log("CONNECTION ENDED: ", e)
    })
})
