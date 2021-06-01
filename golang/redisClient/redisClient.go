package redisClient
import (
    "fmt"
)
// import (
// 	"log" // logging 
// 	"net" // socket programming
// 	"fmt" // print to terminal
// 	"strings" // convert to string
// 	"encoding/json" // for data
// 	"time" // unix for timeout connection
// 	"github.com/go-redis/redis/v8" // redis server
// 	"github.com/joho/godotenv" // dot env
// 	"os" // mengakses dot enviroment variables
// 	"reflect" // mengetahui type dari variable
// 	"strconv" // convert string to int vice versa
// 	"context" // memberi tahu harus di tahan berapa lama (misal goroutine)
// )


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