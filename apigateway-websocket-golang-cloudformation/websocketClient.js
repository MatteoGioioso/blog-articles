export class WebsocketAPI {
  constructor(webSocketEndpoint, token) {
    const websocketFullUrl = `${webSocketEndpoint}?Auth=${token}`;
    this.socket = new WebsocketAPI(websocketFullUrl);
    this.webSocketOnClose();
    this.webSocketOnError();

    // Close the web socket in case of refresh or close browser
    window.onbeforeunload = () => {
      if (this.socket) {
        this.socket.onclose = function() {}; // disable onclose handler first
        this.socket.close();
        console.log('closing webSocket');
      }
    };
  }

  webSocketOnOpen(callback) {
    this.socket.addEventListener('open', () => {
      console.log('connected');
      if (callback) {
        callback();
      }
    });
  }

  webSocketOnError(callback) {
    this.socket.addEventListener('error', (e) => {
      console.log(e);
      this.socket.close();
      if (callback) {
        callback(e);
      }
    });
  }

  closeWebsocket() {
    if (this.socket) {
      this.socket.close();
    }
  }

  webSocketOnMessage(callback) {
    this.socket.addEventListener('message', (e) => {
      try {
        const data = JSON.parse(e.data);
        console.log(data);
        return callback(null, data);
      } catch (e) {
        callback(e, {});
      }
    });
  }

  webSocketOnClose(callback) {
    this.socket.addEventListener('close', (e) => {
      console.log('disconnected: ', e.code, e.reason);
      if (callback) {
        callback({ code: e.code, reason: e.reason });
      }
    });
  }
}

const myUrl = "wss://api-id.execute-api.region.amazonaws.com/stage"
const client = new WebsocketAPI(myUrl, "123")
client.webSocketOnOpen(() => {
  console.log("connectedgi")
})
client.webSocketOnError((e) => {
  console.log(e.message)
})
client.webSocketOnMessage((e, data) => {
  if (e){
    console.log(e.message)
  }
  console.log(data)
})
