'use strict';

const stream = require('stream');
const util = require('util');
const {
    TCP,
    TCPConnectWrap,
    constants: TCPConstants
    // This is a hack to access the bindings from outside nodejs
    // THIS WILL ONLY WORK WITH NODE < v10.x
} = process.binding('tcp_wrap');
const { WriteWrap } = process.binding('stream_wrap');

function Socket() {
    this._handle = null;
    stream.Duplex.call(this, {});
    // shut down the socket when we're finished with it.
    this.on('end', onReadableStreamEnd);
}
util.inherits(Socket, stream.Duplex);

Socket.prototype.end = function(data, encoding, callback) {
    // This will end the stream
    stream.Duplex.prototype.end.call(this, data, encoding, callback);
    return this;
};

// ============== READ ======================
// Just call handle.readStart until we have enough in the buffer
Socket.prototype._read = function(n) {
    // Probably call int LibuvStreamWrap::ReadStart()
    const err = this._handle.readStart();
    if (err) {
        throw new Error("read error: " + err.message)
    }
};

function onStreamRead(nread, buf) {
    const self = this.self
    if (buf){
        self.emit('data', buf)
    } else {
        self.emit('end')
    }
}

// Called when the 'end' event is emitted.
function onReadableStreamEnd() {
    // This will close the socket
    this._handle.close(() => {
        this.emit('close');
    });
}

// ============== WRITE =====================
Socket.prototype.write = function(data, cb) {
    // stream_wrap.cc
    const req = new WriteWrap();
    req.handle = this._handle;
    req.oncomplete = afterWrite;
    req.async = false;

    // 3. WRITE to the socket
    // writeBuffer is probably int StreamBase::WriteBuffer(const FunctionCallbackInfo<Value>& args)
    // in stream_base.cc
    const err = req.handle.writeBuffer(req, data);
    if (err){
        throw new Error('Write error' + err.message)
    }
    cb()
};

function afterWrite(status, handle, err) {
    console.log("after write called")
}

// ============== CONNECTION ==========================

Socket.prototype.connect = function(port, host, cb) {
    // 1. INITIALIZE the TCP socket
    // TCP and TCPConnectWrap are from tcp_wrap.cc
    this._handle = new TCP(TCPConstants.SOCKET);
    this._handle.self = this;
    this._handle.onread = onStreamRead;

    const req = new TCPConnectWrap();
    req.oncomplete = afterConnect;
    req.address = host;
    req.port = port;

    // 2. CONNECT to the server
    const err = this._handle.connect(req, host, port);
    if (err) {
        throw new Error("connect error: " + err.message)
    }

    cb()

    return this;
};


function afterConnect(status, handle, req, readable, writable) {
    if (readable && writable){
        const self = handle.self;
        // Emit the connect event
        self.emit('connect')
        self.emit('ready')
    } else {
        throw new Error("Socket is not duplex")
    }
}


module.exports = {
    Socket
};
