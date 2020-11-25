const {
    TCP,
    constants: TCPConstants
} = internalBinding('tcp_wrap');


function Socket(options) {}

Socket.prototype.connect = function(...args) {
    let normalized;
    const options = normalized[0];
    const cb = normalized[1];

    if (this.write !== Socket.prototype.write)
        this.write = Socket.prototype.write;

    const { path } = options;
    const pipe = !!path;
    debug('pipe', pipe, path); // pipe is false from the DEBUG

    if (!this._handle) {
        this._handle = new TCP(TCPConstants.SOCKET);
        initSocketHandle(this);
    }

    if (cb !== null) {
        this.once('connect', cb);
    }

    this._unrefTimer();

    this.connecting = true;
    this.writable = true;

    lookupAndConnect(this, options);

    return this;
}


function lookupAndConnect(self, options) {
    const { localAddress, localPort } = options;
    const host = options.host || 'localhost';
    let { port } = options;

    // If host is an IP, skip performing a lookup
    debug('connect: find host', host);
    // debug('connect: dns options', dnsopts);
    self._host = host;
    const lookup = options.lookup || dns.lookup;
    defaultTriggerAsyncIdScope(self[async_id_symbol], function() {
        lookup(host, dnsopts, function emitLookup(err, ip, addressType) {
            self.emit('lookup', err, ip, addressType, host);

            // It's possible we were destroyed while looking this up.
            // XXX it would be great if we could cancel the promise returned by
            // the look up.
            if (!self.connecting) return;

            if (err) {
                // net.createConnection() creates a net.Socket object and immediately
                // calls net.Socket.connect() on it (that's us). There are no event
                // listeners registered yet so defer the error event to the next tick.
                process.nextTick(connectErrorNT, self, err);
            } else if (addressType !== 4 && addressType !== 6) {
                err = new ERR_INVALID_ADDRESS_FAMILY(addressType,
                    options.host,
                    options.port);
                process.nextTick(connectErrorNT, self, err);
            } else {
                self._unrefTimer();
                defaultTriggerAsyncIdScope(
                    self[async_id_symbol],
                    internalConnect,
                    self, ip, port, addressType, localAddress, localPort
                );
            }
        });
    });
}
