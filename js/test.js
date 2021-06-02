var redis = require("redis");

var subscriber = redis.createClient({ host: "18.141.9.99", port: 6379 });
console.log("initilized")
subscriber.on("message", function (channel, message) {
    console.log("Message: " + message + " on channel: " + channel + " is arrive!");
});

subscriber.subscribe("join");

