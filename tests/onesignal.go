package main

import (
	"fmt"
	"log"

	"github.com/tbalthazar/onesignal-go"
)

var (
	appID   string
	appKey  string
	userKey string
)

func CreateNotifications(client *onesignal.Client) string {
	fmt.Println("### CreateNotifications ###")
	//playerID := "2ea61d88-5b9e-472a-a097-e41f407f6ab0" // valid
	 playerID := "a5f1f696-550a-44ba-9ae7-9ad5f0f5414d" // invalid
	notificationReq := &onesignal.NotificationRequest{
		AppID:            appID,
		Headings:         map[string]string{"en": "Projud"},
		Contents:         map[string]string{"en": "English message"},
		//Tags:             map[string]string{"key":"collapse_id","relation":"=","value":"teste"},
		IsIOS:            true,
		IncludePlayerIDs: []string{playerID},
	}

	createRes, res, err := client.Notifications.Create(notificationReq)
	if err != nil {
		fmt.Printf("--- res:%+v, err:%+v\n", res)
		log.Fatal(err)
	}
	fmt.Printf("--- createRes:%+v\n", createRes)
	fmt.Println()

	return createRes.ID
}
func main() {
	//appID = "7a23693a-e65d-4893-a68f-eae8acc287b8"
	appID = "be605398-9907-4849-92ee-8d4584c587ba"
	client := onesignal.NewClient(nil)
	//11872906087d4a7a30fa9bc711cb19743b05d43faeddd984c9cb799e7aad35fc
	//client.AppKey = "NWVkNzFjZjktOWE4Yy00NzUxLWJmMzgtOWNkNjA1YTczYWE0"
	client.AppKey = "Yjc5YjlhMjMtOGU0MC00ZGYwLWJmNTctYzZkZGMwYTc0ZTE1"
	//client.UserKey = "11872906087d4a7a30fa9bc711cb19743b05d43faeddd984c9cb799e7aad35fc"
	CreateNotifications(client)
}
