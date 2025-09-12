package application

import (
	"sync"

	"github.com/mahdi-cpp/iris-tools/image_loader"
	"github.com/mahdi-cpp/upload-service/internal/config"
)

type AppManager struct {
	mu sync.RWMutex
	//rdb                  *redis.Client
	IconImageLoader      *image_loader.ImageLoader
	OriginalImageLoader  *image_loader.ImageLoader
	ThumbnailImageLoader *image_loader.ImageLoader
}

func NewAppManager() (*AppManager, error) {

	manager := &AppManager{
		//rdb: redis.NewClient(&redis.Options{
		//	Addr: "localhost:50001",
		//	DB:   0,
		//}),
	}

	manager.IconImageLoader = image_loader.NewImageLoader(5000, config.UploadDir, 0)
	manager.OriginalImageLoader = image_loader.NewImageLoader(100, config.UploadDir, 0)
	manager.ThumbnailImageLoader = image_loader.NewImageLoader(5000, config.UploadDir, 0)

	//// Check the connection to Redis.
	//_, err := manager.rdb.Ping(ctx).Result()
	//if err != nil {
	//	log.Fatalf("Could not connect to Redis: %v", err)
	//}
	//fmt.Println("Successfully connected to Redis!")
	//
	//// Create a new PubSub client and subscribe to a channel.
	//pubsub := manager.Subscribe(ctx, "mychannel")
	//
	//// Make sure to close the subscription when the program exits.
	//defer pubsub.Close()
	//
	//// Wait for confirmation that the subscription is successful.
	//_, err = pubsub.Receive(ctx)
	//if err != nil {
	//	log.Fatalf("Failed to subscribe: %v", err)
	//}
	//
	//fmt.Println("Subscribed to 'mychannel'. Listening for messages...")
	//
	//// Listen for incoming messages in a loop.
	//ch := pubsub.Channel()
	//for msg := range ch {
	//	fmt.Printf("Received message from channel %s: %s\n", msg.Channel, msg.Payload)
	//}

	return manager, nil
}
