package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"github.com/lib/pq"
	"github.com/xhyrom/swish/pkg/database"
)

func main() {
	fmt.Println("Hello, World!")

	creds := database.GetCredentials()
	fmt.Printf("Credentials are: %+v\n", creds)

	db := database.NewDatabase(creds)
	queue := db.GetQueue()

	for _, track := range queue {
		fmt.Printf("Video ID: %s", track.VideoId)
		if track.From != nil {
			fmt.Printf(" from %s", *track.From)
		}
		if track.To != nil {
			fmt.Printf(" to %s", *track.To)
		}
		fmt.Println()
	}

	listener := pq.NewListener(creds.ConnectionString(), 10 * time.Second, time.Minute, nil)
	err := listener.Listen("table_changes")
	if err != nil {
		panic(err)
	}

	for range time.Tick(time.Second) {
		select {
		case notification := <-listener.Notify:
			fmt.Println("Received notification: ", notification)

			var prettyJSON bytes.Buffer
			err := json.Indent(&prettyJSON, []byte(notification.Extra), "", "\t")

			if err != nil {
				panic(err)
			}

			fmt.Println(prettyJSON.String())
		}
	}
}