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
    // {
    //     ID: 1,
    //     used: false,
    //     ws: undefined
    // },
    // {
    //     ID: 2,
    //     used: false,
    //     ws: undefined
    // },
    // {
    //     ID: 3,
    //     used: false,
    //     ws: undefined
    // }
]

const trees = [
    
    {
        x:2,
        y:0,
        z:4
    }
]

wss.on('connection', function connection(ws) {
    console.log("made connection")
    ws.on('message', function incoming(message) {
        console.log("reciving msg")
        console.log(message)
        message = JSON.parse(message);
        switch (message.channel) {
            case "join":
                //kasih tau lokasi player2 sebelumnya ke player yg baru join,
                players.forEach((player, index) => {
                    ws.send(JSON.stringify({
                        channel: "join",
                        ID: index
                    }))
                })
                console.log(`recieved message from join channel`)
                console.log(`selecting player ID...`)
                players.push(
                    {
                        ID: players.length,
                        ws: undefined
                    });
                const selectedPlayer = players[players.length - 1]//last index
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
                    broadcastToOthers("join", selectedPlayer.ID);
                });
                break;

            default:
                break;
        }
        console.log('received: %s', message);
    });
    ws.on("disconnect", function () {
        console.log("disconnected")
    });

    ws.on('close', function close() {
        const indexDP = players.findIndex(obj => obj.ws == ws);
        console.log(`player count ${players.length}`);
        console.log(`${indexDP} disconnected`);
        if (indexDP != -1) {
            const ID = players[indexDP].ID;
            console.log(`${ID} disconnected from the server`)
            players.splice(indexDP, 1);
            publisher.publish("leave", ID, function () {

                console.log("done publish redis!")
            });
        }
    });
    // ws.send('something');
});

function broadcastToOthers(channel, index) {
    console.log("broadcasting")
    var msg = {
        channel,
        ID: index
    }
    players.forEach((element, ind) => {
        if (ind == index) { } //tidak perlu broadcast ke dia lagi
        else
            element.ws.send(JSON.stringify(msg));
    });
}




