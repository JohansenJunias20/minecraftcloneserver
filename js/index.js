const WebSocket = require('ws');
var redis = require("redis");
require('dotenv').config()
const wss = new WebSocket.Server({ port: process.env.WS_PORT });
console.log(`Web Socket ready to serve on port ${process.env.WS_PORT}`)
const publisher = redis.createClient({ host: process.env.REDIS_HOST, port: process.env.REDIS_PORT });
const subscriber = redis.createClient({ host: process.env.REDIS_HOST, port: process.env.REDIS_PORT });

// redisClient.publish("join", "testing", function () {
//     // selectedPlayer.ws = ws
//     console.log("done publish redis!")
// });

subscriber.on("message", function (channel, message) {
    console.log("Message: " + message + " on channel: " + channel + " is arrive!");
    if (channel == "verified") {
        const ID = message;
        var player = players.find(player => player.ID == ID)
        player.ws.send(JSON.stringify({ channel: "verified", status: true }));
    }
});
subscriber.subscribe("verified");
const players = [
    {
        ID: 1,
        used: false,
        ws: undefined
    },
    {
        ID: 2,
        used: false,
        ws: undefined
    },
    {
        ID: 3,
        used: false,
        ws: undefined
    }]

wss.on('connection', function connection(ws) {
    console.log("made connection")
    ws.on('message', function incoming(message) {
        console.log(message)
        message = JSON.parse(message);
        switch (message.channel) {
            case "join":
                console.log(`recieved message from join channel`)
                console.log(`selecting player ID...`)
                const selectedPlayer = players.find(obj => obj.used == false)
                selectedPlayer.ws = ws
                selectedPlayer.used = true
                console.log(`player ID selected`)
                var ID = selectedPlayer.ID;
                console.log(`reply client with ID`)
                ws.send(JSON.stringify({
                    channel: "ID",
                    ID
                }));

                console.log(`publish ID to redis in the join channel`)
                publisher.publish("join", selectedPlayer.ID, function () {
                    // selectedPlayer.ws = ws
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