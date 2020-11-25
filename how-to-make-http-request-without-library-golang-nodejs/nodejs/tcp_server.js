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
        const http = `HTTP/1.1 200 OK
Date: Mon, 27 Jul 2009 12:28:53 GMT
Server: Apache/2.2.14 (Win32)
Last-Modified: Wed, 22 Jul 2009 19:15:56 GMT
Content-Length: ${body.length}
Content-Type: text/html
Connection: Closed

${body}`
        const buff = Buffer.from(http, "utf-8")
        conn.write(buff);
    }

    function onConnClose() {
        console.log('connection from %s closed', remoteAddress);
    }

    function onConnError(err) {
        console.log('Connection %s error: %s', remoteAddress, err.message);
    }
}
