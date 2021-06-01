package main

import (
	"log" // logging 
	"net" // socket programming
	"fmt" // print to terminal
	"strings" // convert to string
	"encoding/json" // for data
	"time" // unix for timeout connection
)

const durationTimeOut = 240 int // dalam second

//print 1 object string only
func pr1(message string){
	fmt.Println(message)
}
type Position struct{
	X float32
	Y float32
}
type Response struct {
	Channel string
	Position Position
}

func main() {
	// listen to incoming udp packets
	// var clientsIP []string
	// var clientsPort []int
	go removeIdleClients() // make new goroutine untuk menghapus client yang idle

	pc, err := net.ListenPacket("udp", ":28000")
	if err != nil {
		log.Fatal(err)
	}
	defer pc.Close()

	for {
		buf := make([]byte, 1024)
		n, addr, err := pc.ReadFrom(buf)
		if err != nil {
			continue
		}
		go serve(pc, addr, buf[:n])
	}

}

var clients []net.Addr
var timeOut []string
var clientCount int
func serve(pc net.PacketConn, addr net.Addr, buf []byte) {
	appendClient(addr)

	fmt.Println("some one connected")

	description := strings.Split(string(addr.String()),":")
	fmt.Println("address: ",description[0])
	fmt.Println("port: ",description[1])

	broadcast(string(buf))
}

func appendClient(addr net.Addr){
	for i := 0; i < clientCount; i++ {
		
		if(clients[i].String() == addr.String()){ //sudah ada
			timeOut[i] = string(time.Now().Unix())
			return
		}
	}

	clients = append(clients, addr)
	timeOut =  timeOut(timeOut,time.Now().Unix())
	clientCount++

	return

}
func removeIdleClients(){
	//loop dipanggil tiap 10 detik sekali
	ticker := time.NewTicker(10 * time.Second)
	quit := make(chan struct{})
	for {
		select {
		 case <- ticker.C:
			for i := 0; i < clientCount; i++ {
				if(time.Now().Unix() - int(timeOut[i]) > durationTimeOut ){
					clients=pop(clients,i)
					clientCount--
				}
			}
		 case <- quit:
			 ticker.Stop()
			 return
		 }
	 }

}
func broadcast(response string){
	fmt.Println(response)
	var responseJSON Response
	err :=json.Unmarshal([]byte(response), &responseJSON)
	if err!=nil{
		fmt.Println(err)
	}else {
		fmt.Println(responseJSON.Channel)
	}

	fmt.Println("client count:",clientCount)
	fmt.Println("clients:",clients)
	for i := 0; i < clientCount; i++ {
		switch responseJSON.Channel {
		case "position":
			fmt.Println("broadcasting pos to all player")
			fmt.Println(responseJSON.Position)
		default:
			fmt.Printf("unknown response :")
			pr1(response)
		}
	}
}

func pop(slice []int, s int) []int {
    return append(slice[:s], slice[s+1:]...)
}