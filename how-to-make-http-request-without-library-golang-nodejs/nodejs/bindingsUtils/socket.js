'use strict';

const stream = require('stream');
const util = require('util');
const {
    TCP,
    TCPConnectWrap,
    constants: TCPConstants
} = process.binding('tcp_wrap');

// called when creating new Socket, or when re-using a closed Socket
function initSocketHandle(self) {
    self._sockname = null;

    // Handle creation may be deferred to bind() or connect() time.
    if (self._handle) {
        // Original API use a Symbol from the bindings
        // for simplicity we use just a string
        self._handle['owner_symbol'] = self;
        self._handle.onread = onStreamRead;
    }
}

function Socket() {
    if (!(this instanceof Socket)) return new Socket();

    this._handle = null;
    stream.Duplex.call(this, {});
    // shut down the socket when we're finished with it.
    this.on('end', onReadableStreamEnd);

    initSocketHandle(this);
}
util.inherits(Socket, stream.Duplex);


Socket.prototype.end = function(data, encoding, callback) {
    stream.Duplex.prototype.end.call(this, data, encoding, callback);
    return this;
};

// Called when the 'end' event is emitted.
function onReadableStreamEnd() {
    console.log("Readable stream end")
}

// ============== READ ======================
// Just call handle.readStart until we have enough in the buffer
Socket.prototype._read = function(n) {
    this._handle.reading = true;

    // Probably call int LibuvStreamWrap::ReadStart()
    const err = this._handle.readStart();
    if (err) {
        throw new Error("read error: " + err.message)
    }
};

function onStreamRead(nread, buf) {
    const self = this['owner_symbol']
    if (buf){
        self.emit('data', buf)
    }
}

// ============== WRITE =====================
Socket.prototype.write = function(data, cb) {
    const req = {}
    req.handle = this._handle;
    req.oncomplete = afterWrite;
    req.async = false;

    // Handle only utf8
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
    this._handle = new TCP(TCPConstants.SOCKET);
    initSocketHandle(this);

    this.connecting = false;
    this.writable = true;

    const req = new TCPConnectWrap();
    req.oncomplete = afterConnect;
    req.address = host;
    req.port = port;

    const err = this._handle.connect(req, host, port);

    if (err) {
        console.log(err)
    }

    cb()

    return this;
};


function afterConnect(status, handle, req, readable, writable) {
    if (readable && writable){
        const self = handle['owner_symbol'];
        self.emit('connect')
        self.emit('ready')
    } else {
        throw new Error("Socket is not duplex")
    }
}


module.exports = {
    Socket
};
