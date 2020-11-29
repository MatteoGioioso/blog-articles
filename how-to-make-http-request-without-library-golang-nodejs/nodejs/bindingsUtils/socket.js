'use strict';

const stream = require('stream');
const util = require('util');
const {
    TCP,
    TCPConnectWrap,
    constants: TCPConstants
} = process.binding('tcp_wrap');
const { WriteWrap } = process.binding('stream_wrap');
const {EventEmitter} = require('events')

const kLastWriteQueueSize = Symbol('lastWriteQueueSize');

const debug = util.debuglog('net');

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

function onStreamRead(nread, buf) {
    console.log("On Stream Read called")
}

const kBytesRead = Symbol('kBytesRead');
const kBytesWritten = Symbol('kBytesWritten');

function Socket(options) {
    if (!(this instanceof Socket)) return new Socket(options);

    this.connecting = false;
    this._hadError = false;
    this._handle = null;
    this._parent = null;
    this._host = null;
    this[kLastWriteQueueSize] = 0;

    options = util._extend({}, options);
    options.readable = options.readable || false;
    options.writable = options.writable || false;
    const { allowHalfOpen } = options;

    // Prevent the "no-half-open enforcer" from being inherited from `Duplex`.
    options.allowHalfOpen = true;
    // For backwards compat do not emit close on destroy.
    options.emitClose = false;
    stream.Duplex.call(this, options);

    // Default to *not* allowing half open sockets.
    this.allowHalfOpen = Boolean(allowHalfOpen);

    this._handle = options.handle; // private
    // this[async_id_symbol] = getNewAsyncId(this._handle);

    // shut down the socket when we're finished with it.
    this.on('end', onReadableStreamEnd);

    initSocketHandle(this);

    // if we have a handle, then start the flow of data into the
    // buffer.  if not, then this will happen when we connect
    if (this._handle && options.readable !== false) {
        this.read(0);
    }

    // Used after `.destroy()`
    this[kBytesRead] = 0;
    this[kBytesWritten] = 0;
}
util.inherits(Socket, stream.Duplex);

Object.defineProperty(Socket.prototype, '_connecting', {
    get: function() {
        return this.connecting;
    }
});

Object.defineProperty(Socket.prototype, 'pending', {
    get() {
        return !this._handle || this.connecting;
    },
    configurable: true
});


Object.defineProperty(Socket.prototype, 'readyState', {
    get: function() {
        if (this.connecting) {
            return 'opening';
        } else if (this.readable && this.writable) {
            return 'open';
        } else if (this.readable && !this.writable) {
            return 'readOnly';
        } else if (!this.readable && this.writable) {
            return 'writeOnly';
        } else {
            return 'closed';
        }
    }
});


Object.defineProperty(Socket.prototype, 'bufferSize', {
    get: function() {
        if (this._handle) {
            return this[kLastWriteQueueSize] + this.writableLength;
        }
    }
});

// Just call handle.readStart until we have enough in the buffer
Socket.prototype._read = function(n) {
    debug('_read');
    debug('Socket._read readStart');
    this._handle.reading = true;
    const err = this._handle.readStart();
    if (err) {
        throw new Error("read error: " + err.message)
    }
};


Socket.prototype.end = function(data, encoding, callback) {
    stream.Duplex.prototype.end.call(this, data, encoding, callback);
    return this;
};

// Called when the 'end' event is emitted.
function onReadableStreamEnd() {
    maybeDestroy(this);
}

// Call whenever we set writable=false or readable=false
function maybeDestroy(socket) {
    if (!socket.readable &&
        !socket.writable &&
        !socket.destroyed &&
        !socket.connecting &&
        !socket.writableLength) {
        socket.destroy();
    }
}

// ============== WRITE =====================
Socket.prototype._write = function(data, encoding, cb) {
    const req = new WriteWrap();
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
    if (this.write !== Socket.prototype.write)
        this.write = Socket.prototype.write;

    if (!this._handle) {
        this._handle = new TCP(TCPConstants.SOCKET);
        initSocketHandle(this);
    }

    if (cb !== null) {
        cb()
    }

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

    return this;
};


function afterConnect(status, handle, req, readable, writable) {
    const self = handle['owner_symbol'];
    self.emit('connect')
    self.emit('ready')
}

Socket.prototype.ref = function() {
    if (!this._handle) {
        this.once('connect', this.ref);
        return this;
    }

    if (typeof this._handle.ref === 'function') {
        this._handle.ref();
    }

    return this;
};


Socket.prototype.unref = function() {
    if (!this._handle) {
        this.once('connect', this.unref);
        return this;
    }

    if (typeof this._handle.unref === 'function') {
        this._handle.unref();
    }

    return this;
};

module.exports = {
    Socket
};
