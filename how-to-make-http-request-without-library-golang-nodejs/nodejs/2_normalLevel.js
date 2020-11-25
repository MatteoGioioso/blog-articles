const http = require('http');
const host = "http://localhost:9000/hello-world"

http.get(host, (resp) => {
    let data = '';

    // A chunk of data has been recieved.
    resp.on('data', (chunk) => {
        data += chunk;
    });

    // The whole response has been received. Print out the result.
    resp.on('end', () => {
        console.log("DATA: ", JSON.parse(data))
    });
})
