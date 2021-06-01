const dgram = require('dgram');
var client = dgram.createSocket('udp4');
const PORT = 28000;
const HOST = "localhost";
client.on('listening', function () {
    var address = client.address();
    console.log('UDP Server listening on ' + address.address + ":" + address.port);
});

client.on('message', function (message, remote) {

    console.log(remote.address + ':' + remote.port +' - ' + message);

});
var message = "hello world hahaha";
client.send(message, 0, message.length, PORT, HOST, function(err, bytes) {

    if (err) throw err;
    console.log('UDP message sent to ' + HOST +':'+ PORT);

})