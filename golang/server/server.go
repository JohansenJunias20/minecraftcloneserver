package main
import (
	"log" // logging 
	"net" // socket programming
	"fmt" // print to terminal
	"strings" // convert to string
	// "encoding/json" // for data
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
// var topicJoin *redis.PubSub
func connectRedis() {
	
	REDIS_HOST := os.Getenv("REDIS_HOST")
	REDIS_PORT := os.Getenv("REDIS_PORT")
    rdb = redis.NewClient(&redis.Options{
		Addr:     REDIS_HOST+":"+REDIS_PORT,
        Password: "", // no password set
        DB:       0,  // use default DB
    })
	fmt.Println("connected to redis server...")

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

// func startSubscribeRedis(){
// 	topicJoin = rdb.Subscribe(ctx, "join")
// 	channelJoin := topicJoin.Channel()
// 	fmt.Println(reflect.TypeOf(channelJoin))
// 	// Itterate any messages sent on the channel
// 	// for msg := range channelJoin {
// 	// 	fmt.Println("msg from rediss: ")
// 	// 	// go appendClient(msg.Payload)
// 	// 	fmt.Println(msg.Payload)
// 	// }

// 	go SubscribeJoin(channelJoin)
// }

// func SubscribeJoin(channelJoin <-chan *redis.Message){
// 	for msg := range channelJoin {
// 		newNickname := msg.Payload
		
// 		var tmp Client
// 		tmp.Name = newNickname
// 		tmp.Client = nil
// 		tmp.LastTimeConnect = time.Now().Unix()
// 		clients = append(clients,tmp)
// 		// fmt.Println("msg from rediss: ")
// 		// go appendClient(msg.Payload)
// 		// fmt.Println(msg.Payload)
// 	}
// }


//print 1 object string only
// func pr1(message string){
// 	fmt.Println(message)
// }
type Position struct{
	X float32
	Y float32
}
type Response struct {
	Channel string
	Position Position
}

var pc net.PacketConn;

func main() {
	//inisialisasi dot env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	// connectRedis()
 	UDP_PORT := os.Getenv("UDP_PORT")
	connectRedis()

	go removeIdleClients() // make new goroutine untuk menghapus client yang idle
	pc, err = net.ListenPacket("udp", ":" + UDP_PORT)
	fmt.Println("listening on port ", UDP_PORT)
	if err != nil {
		log.Fatal(err)
	}
	defer pc.Close()

	for {
		buf := make([]byte, 4096)
		n, addr, err := pc.ReadFrom(buf)
		fmt.Println("incoming message")
		if err != nil {
			continue
		}
		go serve(pc, addr, buf[:n])
	}

}

type Client struct {
	ID int
	Client net.Addr
	LastTimeConnect int64
}

var clients =  []Client { 
    Client {
		ID :1,
		Client: nil,
		LastTimeConnect: 0,
    },
    Client {
		ID: 2,
		Client: nil,
		LastTimeConnect: 0,
    },
    Client {
		ID: 3,
		Client: nil,
		LastTimeConnect: 0,
    },
}
// var timeOut []string
var clientCount = 3
func serve(pc net.PacketConn, addr net.Addr, buf []byte) {
	
	arr := strings.Split(string(buf), "|")
	if(len(arr)>=1){
		return;
	}
	channel := arr[0]
	switch channel {
		case "position":// msg: position|x:30;y:25;z:40

			if(len(arr)!=3){
				return;
			}
			x,errx := strconv.ParseFloat(strings.Split(strings.Split(arr[2],";")[0], ":")[1], 64)
			y,erry := strconv.ParseFloat(strings.Split(strings.Split(arr[2],";")[1], ":")[1], 64)
			z,errz := strconv.ParseFloat(strings.Split(strings.Split(arr[2],";")[2], ":")[1], 64)
			if(errx!=nil || erry != nil || errz != nil){
				return;
			}
			ID, err := strconv.Atoi(string(arr[1]))
			if(err != nil){
				return;
			}
			broadcastPosition(x,y,z,ID)
			break;

		case "join": // msg: join|
			for i := 0; i < clientCount; i++ {
				if(clients[i].Client == nil){
					err := rdb.Publish(ctx, "id_udp_client",clients[i].ID).Err()
					if err != nil {
						panic(err)
					}
					clients[i].Client = addr
					clients[i].LastTimeConnect =  time.Now().Unix()
					break;
				}
			}
	}
	fmt.Println("some one connected")

	description := strings.Split(string(addr.String()),":")
	fmt.Println("address: ",description[0])
	fmt.Println("port: ",description[1])
	
	// broadcast(string(buf))
}
func broadcastPosition(x float64, y float64,z float64, ID int){
	for i := 0; i < clientCount; i++ {
		if(clients[i].ID != ID){
			pc.WriteTo([]byte("{\"ID\":"+strconv.Itoa(ID)+", \"x\":"+FloatToString(x)+", \"y\":"+FloatToString(y)+"\"z\":"+FloatToString(z)+" }"), clients[i].Client)
		}
	}
}

func FloatToString(input_num float64) string {
    // to convert a float number to a string
    return strconv.FormatFloat(input_num, 'f', 6, 64)
}

func removeIdleClients(){
	//loop dipanggil tiap 10 detik sekali
	ticker := time.NewTicker(10 * time.Second)
	quit := make(chan struct{})
	for {
		select {
		 case <- ticker.C:
			for i := 0; i < clientCount; i++ {
				if(clients[i].Client == nil){
					continue
				}
				old := clients[i].LastTimeConnect
				if(time.Now().Unix() - old > int64(DurationTimeOut)){
					popByIndex(clients,i)
				}
			}
		 case <- quit:
			 ticker.Stop()
			 return
		 }
	 }

}
// func broadcast(response string){
// 	fmt.Println(response)
// 	var responseJSON Response
// 	err :=json.Unmarshal([]byte(response), &responseJSON)
// 	if err!=nil{
// 		fmt.Println(err)
// 	}else {
// 		fmt.Println(responseJSON.Channel)
// 	}


// }

func popByIndex(slice []Client, s int) {
	slice[s].Client = nil
	slice[s].LastTimeConnect = 0
}