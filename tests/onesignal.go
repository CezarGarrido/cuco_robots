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
// CreateNotifications ewdwd
// wfdwefwefw
// efwffrefrefe
func CreateNotifications(client *onesignal.Client) string {
	fmt.Println("### CreateNotifications ###")
	//playerID := "2ea61d88-5b9e-472a-a097-e41f407f6ab0" // valid
	//playerID := "a5f1f696-550a-44ba-9ae7-9ad5f0f5414d" // invalid
	//playerID := "bfb6f92a-ca46-4407-883a-db0e340c6947"
	//playerID:="bfb6f92a-ca46-4407-883a-db0e340c6947"
	playerID:="44e41ef0-7991-487f-9ad4-5c33c877a392"
	
	notificationReq := &onesignal.NotificationRequest{
		AppID:           appID,
		Headings:        map[string]string{"en": "Uma nova Publicação"},
		Contents:        map[string]string{"en": "Processo 0805260-53.2016.8.12.0002"},
		ADMGroupMessage: map[string]string{"en": "You have $[notif_count] new messages"},
		AndroidGroup:    "projud_publicacoes",
		//$[notif_count]
		SmallIcon:           "icone_teste",
		ADMSmallIcon:        "icone_teste",
		//CollapseID:"teste",
		AndroidGroupMessage: map[string]string{"en": "Você tem $[notif_count] novas publicações"},
		/*Tags: []interface{}{
			map[string]interface{}{
				"key": "collapse_id", "relation": "=", "value": "male",
			},
		},*/
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
// Start dcwewcw
// cewcewcw
// ewefwfwrfefwce
func main() {
	//appID = "7a23693a-e65d-4893-a68f-eae8acc287b8"
	appID = "be605398-9907-4849-92ee-8d4584c587ba"
	//appID = "3e80a56a-ba22-4fce-b5af-35056c0b4979"
	client := onesignal.NewClient(nil)
	//11872906087d4a7a30fa9bc711cb19743b05d43faeddd984c9cb799e7aad35fc
	//client.AppKey = "NWVkNzFjZjktOWE4Yy00NzUxLWJmMzgtOWNkNjA1YTczYWE0"
	client.AppKey = "Yjc5YjlhMjMtOGU0MC00ZGYwLWJmNTctYzZkZGMwYTc0ZTE1"
	//client.AppKey = "NTI1ODdhODktNGNjYS00MTkzLWI3MWQtNGM2NGIyMDAyYmVj"
	//client.UserKey = "11872906087d4a7a30fa9bc711cb19743b05d43faeddd984c9cb799e7aad35fc"
	CreateNotifications(client)
	// Output:
	// "Arslan"
	// 123456
	// true
}
