const net = require('net');
const server = net.createServer();

server.on('connection', handleConnection);
server.listen(9000, function() {
    console.log('server listening to %j', server.address());
});

function handleConnection(conn) {
    const remoteAddress = conn.remoteAddress + ':' + conn.remotePort;
    console.log('new client connection from %s', remoteAddress);

    conn.on('data', onConnData);
    conn.once('close', onConnClose);
    conn.on('error', onConnError);

    function onConnData(d) {
        console.log('connection data from %s: %j', remoteAddress, d.toString());
        const body = JSON.stringify({message: "Hello world"})
        conn.write("HTTP/1.1 200 OK\r\n");
        conn.write(`Content-Length: ${body.length}\r\n`);
        conn.write("Content-Type: application/json\r\n");
        conn.write("Connection: Closed\r\n");
        conn.write("\r\n");
        conn.write(body + "\r\n");
        conn.emit("end")
    }

    function onConnClose() {
        console.log('connection from %s closed', remoteAddress);
    }

    function onConnError(err) {
        console.log('Connection %s error: %s', remoteAddress, err.message);
    }
}
