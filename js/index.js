const WebSocket = require('ws');
var redis = require("redis");

const wss = new WebSocket.Server({ port: 5000 });
var publisher = redis.createClient({ host: "18.141.9.99", port: 6379 });

const players = [
    {
        name: "NoobMaster69",
        used: false,
        ws: undefined
    },
    {
        name: "Kata Ilham",
        used: false,
        ws: undefined
    },
    {
        name: "Welost",
        used: false,
        ws: undefined
    }]

wss.on('connection', function connection(ws) {
    ws.on('message', function incoming(message) {
        message = JSON.parse(message);
        switch (message.channel) {
            case "join":
                const selectedPlayer = players.find(obj => obj.used == false)
                publisher.publish("join", selectedPlayer.name, function () {
                    selectedPlayer.ws = ws
                    console.log("done publish redis!")
                });
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