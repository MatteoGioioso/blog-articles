const http = require('http');
const port = 9000;

console.log(`Server is running on port ${port}`);
http.createServer(function(req, res){
    console.log("connection!", req.url)
    res.writeHead(200, {'Content-Type': 'application/json'});
    res.write(JSON.stringify({message: "Hello world!!"}));
    res.end();
}).listen(port);
