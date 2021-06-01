package main
import (
	"log" // logging 
	"net" // socket programming
	"fmt" // print to terminal
	"strings" // convert to string
	"encoding/json" // for data
	"time" // unix for timeout connection
	"github.com/go-redis/redis/v8" // redis server
	"github.com/joho/godotenv" // dot env
	"os" // mengakses dot enviroment variables
	// "reflect" // mengetahui type dari variable
	"strconv" // convert string to int vice versa
	"context" // memberi tahu harus di tahan berapa lama (misal goroutine)
)

var DurationTimeOut = 240  // dalam second

var rdb *redis.Client

var ctx = context.Background()
var topic *redis.PubSub
func connectRedis() {
	
	REDIS_HOST := os.Getenv("REDIS_HOST")
	REDIS_PORT := os.Getenv("REDIS_PORT")
    rdb = redis.NewClient(&redis.Options{
		Addr:     REDIS_HOST+":"+REDIS_PORT,
        Password: "", // no password set
        DB:       0,  // use default DB
    })

	go startSubscribeRedis()

    // err := rdb.Set(ctx, "key", "value", 0).Err()
    // if err != nil {
    //     panic(err)
    // }

    // val, err := rdb.Get(ctx, "key").Result()
    // if err != nil {
    //     panic(err)
    // }
    // fmt.Println("key", val)


}

func startSubscribeRedis(){
	topic = rdb.Subscribe(ctx, "new_users")
	channel := topic.Channel()
	// Itterate any messages sent on the channel
	for msg := range channel {
		fmt.Println("msg from redis: ",msg)
	
	}
}




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

	//inisialisasi dot env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	connectRedis()
 	UDP_PORT := os.Getenv("UDP_PORT")


	go removeIdleClients() // make new goroutine untuk menghapus client yang idle

	pc, err := net.ListenPacket("udp", ":" + UDP_PORT)
	if err != nil {
		log.Fatal(err)
	}
	defer pc.Close()

	for {
		buf := make([]byte, 4096)
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
	timeOut =  append(timeOut,string(time.Now().Unix()))
	clientCount++


}
func removeIdleClients(){
	//loop dipanggil tiap 10 detik sekali
	ticker := time.NewTicker(10 * time.Second)
	quit := make(chan struct{})
	for {
		select {
		 case <- ticker.C:
			for i := 0; i < clientCount; i++ {
				old, err := strconv.Atoi(timeOut[i])
				if(err!= nil){
					panic(err)
				}
				if(int64(time.Now().Unix()) - int64(old) > int64(DurationTimeOut) ){
					clients=popByIndex(clients,i)
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

func popByIndex(slice []net.Addr, s int) []net.Addr {
	
    return append(slice[:s], slice[s+1:]...)
}