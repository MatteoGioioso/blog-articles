const {
    TCP,
    TCPConnectWrap,
    constants: TCPConstants
    // To use this we need to have v10
    // otherwise the new API is internalBinding
    // which is not accessible in userland
} = process.binding('tcp_wrap');

const handle = new TCP(TCPConstants.SOCKET)
handle.onconnection = function (data) {
    console.log("connected!", data)
}

// Add local address binding
// https://docs.oracle.com/cd/E19455-01/806-1017/sockets-47146/index.html
// https://stackoverflow.com/questions/39314086/what-does-it-mean-to-bind-a-socket-to-any-address-other-than-localhost/39314221
const req = new TCPConnectWrap();
req.address = "127.0.0.1";
req.port = 9000;
// req.localAddress = localAddress;
// req.localPort = localPort;
const err = handle.connect(req, "127.0.0.1", 9000)
console.log(err)
