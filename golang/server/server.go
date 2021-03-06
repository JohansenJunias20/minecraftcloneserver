package main

import (
	"fmt"     // print to terminal
	"log"     // logging
	"net"     // socket programming
	"strings" // convert to string

	// "encoding/json" // for data
	// "time" // unix for timeout connection
	"context" // memberi tahu harus di tahan berapa lama (misal goroutine)
	"os"      // mengakses dot enviroment variables
	"reflect" // mengetahui type dari variable
	"strconv" // convert string to int vice versa

	"github.com/go-redis/redis/v8" // redis server
	"github.com/joho/godotenv"     // dot env
)

var DurationTimeOut = 240 // dalam second

var rdb *redis.Client

var ctx = context.Background()
var topicJoin *redis.PubSub
var topicLeave *redis.PubSub

func connectRedis() {

	REDIS_HOST := os.Getenv("REDIS_HOST")
	REDIS_PORT := os.Getenv("REDIS_PORT")
	rdb = redis.NewClient(&redis.Options{
		Addr:     REDIS_HOST + ":" + REDIS_PORT,
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	fmt.Println("connected to redis server...")
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

func startSubscribeRedis() {
	topicJoin = rdb.Subscribe(ctx, "join")
	channelJoin := topicJoin.Channel()
	topicLeave = rdb.Subscribe(ctx, "leave")
	channelLeave := topicLeave.Channel()
	fmt.Println(reflect.TypeOf(channelJoin))
	// Itterate any messages sent on the channel
	// for msg := range channelJoin {
	// 	fmt.Println("msg from rediss: ")
	// 	// go appendClient(msg.Payload)
	// 	fmt.Println(msg.Payload)
	// }

	go SubscribeJoin(channelJoin)
	go SubscribeLeave(channelLeave)
}

//subscribe some one who disconnected
func SubscribeLeave(channelLeave <-chan *redis.Message) {
	for msg := range channelLeave {
		fmt.Println("recieve leave signal from redis")
		fmt.Println(msg.Payload)
		ID := msg.Payload

		id, err := strconv.Atoi(ID)
		if err != nil {
			fmt.Println(err)
		}

		for i := 0; i < len(clients); i++ {
			if clients[i].ID == id {
				clients = popByIndex(clients, i)
			}
		}
		fmt.Println("panjang client: ")
		fmt.Println(len(clients))
		// fmt.Println("msg from rediss: ")
		// go appendClient(msg.Payload)
		// fmt.Println(msg.Payload)
	}
}
func SubscribeJoin(channelJoin <-chan *redis.Message) {
	for msg := range channelJoin {
		fmt.Println("recieve join signal from redis")
		fmt.Println(msg.Payload)
		newID := msg.Payload

		var tmp Client
		id, err := strconv.Atoi(newID)
		if err != nil {
			fmt.Println(err)
		}
		tmp.ID = id

		fmt.Println("append client to clients:")
		fmt.Println(tmp)
		tmp.Client = nil
		clients = append(clients, tmp)
		// fmt.Println("msg from rediss: ")
		// go appendClient(msg.Payload)
		// fmt.Println(msg.Payload)
	}
}

//print 1 object string only
// func pr1(message string){
// 	fmt.Println(message)
// }
type Position struct {
	X float32
	Y float32
}
type Response struct {
	Channel  string
	Position Position
}

var pc net.PacketConn

func main() {
	//inisialisasi dot env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	// connectRedis()
	UDP_PORT := os.Getenv("UDP_PORT")
	connectRedis()

	// go removeIdleClients() // make new goroutine untuk menghapus client yang idle
	pc, err = net.ListenPacket("udp", ":"+UDP_PORT)
	fmt.Println("listening on port ", UDP_PORT)
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

type Client struct {
	ID     int
	Client net.Addr
	// LastTimeConnect int64
}

var clients []Client

type TimeOutCounter struct {
	Addr    net.Addr
	Counter int
	// LastTimeConnect int64
}

var TimeOut []TimeOutCounter

// var clientCount = 3
func serve(pc net.PacketConn, addr net.Addr, buf []byte) {

	// fmt.Println("recieving msg:", string(buf))
	if strings.Contains(string(buf), "|") {

		arr := strings.Split(string(buf), "|")
		if len(arr) <= 1 {
			return
		}

		channel := arr[0]
		if channel == "init" {
			// for i := 0; i < len(TimeOut); i++ {
			// 	if(TimeOut[i].Addr == addr){
			// 		found = true
			// 		TimeOut[i].Counter++
			// 	}
			// }
			// if(!found){
			// 	var tmp TimeOutCounter
			// 	tmp.Addr = addr
			// 	tmp.TimeOutCounter
			// 	TimeOut =  append(TimeOut,tmp)
			// }
			// var ID int
			id, err := strconv.Atoi(arr[2])
			fmt.Println("id: ", id)
			for i := 0; i < len(clients); i++ {
				if err != nil {
					break
				}
				if clients[i].ID == id {
					fmt.Println("id found!, sending verified signal to WS server")
					clients[i].Client = addr
					//send ID to redis
					err := rdb.Publish(ctx, "verified", id).Err()
					if err != nil {
						fmt.Println(err)
					}
					return
				}
			}
			return
		}
	} else {
		broadcast(buf)
	}
}
func FloatToString(input_num float64) string {
	// to convert a float number to a string
	return strconv.FormatFloat(input_num, 'f', 6, 64)
}
func popByIndex(slice []Client, s int) []Client {
	return append(slice[:s], slice[s+1:]...)
}
func broadcast(buf []byte) {
	for i := 0; i < len(clients); i++ {
		pc.WriteTo(buf, clients[i].Client)
	}
}
