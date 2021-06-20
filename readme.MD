# Server UDP & TCP untuk project Grafika Komputer
## Flowchart memapankan koneksi pertama kali
![alt text](https://i.postimg.cc/BQbCzgwV/mc-clone-server-Page-1-1.png)

Warna hijau artinya tereksekusi di NODE JS.<br />
Warna merah tereksekusi di Redis.<br />
Warna biru tereksekusi di Golang.<br />
Warna putih tereksekusi di C# (Client).<br />
WS artinya websocket (NODE JS).<br />

## Deployment
Kami menggunakan [vultr.com](https://vultr.com) untuk host golang UDP server dan NODE JS TCP server.<br />
Kami juga menggunakan [AWS](https://aws.amazon.com) untuk host redis servernya.

