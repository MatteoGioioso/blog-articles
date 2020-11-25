const net = require('net');

const socket = new net.Socket({ readable: true, writable: true})
socket.connect({host: "10.16.0.2", port: 80}, () => {

})
