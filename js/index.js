const WebSocket = require('ws');
// var redis = require("redis");
require('dotenv').config()
const wss = new WebSocket.Server({ port: process.env.WS_PORT });
// var publisher = redis.createClient({ host: process.env.REDIS_HOST, port: process.env.REDIS_PORT });

// publisher.publish("join", "testing", function () {
//     // selectedPlayer.ws = ws
//     console.log("done publish redis!")
// });

subscriber.on("message", function (channel, message) {
    console.log("Message: " + message + " on channel: " + channel + " is arrive!");
});

const players = [
    {
        ID: "NoobMaster69",
        used: false,
        ws: undefined
    },
    {
        ID: "Kata Ilham",
        used: false,
        ws: undefined
    },
    {
        ID: "Welost",
        used: false,
        ws: undefined
    }]

wss.on('connection', function connection(ws) {
    ws.on('message', function incoming(message) {
        message = JSON.parse(message);
        switch (message.channel) {
            case "join":
                const selectedPlayer = players.find(obj => obj.used == false)
                selectedPlayer.ws = ws
                ws.send(JSON.stringify({
                    channel: "ID",
                    ID: selectedPlayer.ID
                }));
             
                break;

            default:
                break;
        }
        console.log('received: %s', message);
    });
    ws.on("disconnect", function () {
        const disconnectedPlayer = players.find(obj => obj.ws == ws);
        publisher.publish("leave", selectedPlayer.name, function () {
            disconnectedPlayer.used = false;
            disconnectedPlayer.ws = undefined;

            console.log("done publish redis!")
        });
    });
    // ws.send('something');
});




wss.on("disconnect", (ws) => {
    players
})