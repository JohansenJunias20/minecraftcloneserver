const WebSocket = require('ws');

const wss = new WebSocket.Server({ port: 5000 });
var PlayerName
wss.on('connection', function connection(ws) {
    ws.on('message', function incoming(message) {
        console.log('received: %s', message);
    });

    ws.send('something');
});

wss.on("disconnect", (ws) => {

})