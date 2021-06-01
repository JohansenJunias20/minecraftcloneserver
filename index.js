const dgram = require('dgram');
const server = dgram.createSocket('udp4');

server.on('error', (err) => {
    console.log(`server error:\n${err.stack}`);
    server.close();
});

server.on('message', (msg, rinfo) => {
    console.log(`server got: ${msg} from ${rinfo.address}:${rinfo.port}`);
    var dummyPos = 0;
    var json = {
        channel:"positions",
        json:
    }
    setInterval(() => {
        server.send(dummyPos, rinfo.port, rinfo.address);
        dummyPos++;
    }, 1000);
});

server.on('listening', () => {
    const address = server.address();
    console.log(`server listening ${address.address}:${address.port}`);
});

server.bind(28000);
// Prints: server listening 0.0.0.0:41234