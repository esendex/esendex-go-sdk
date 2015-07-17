package xesende_test

import (
	"log"
	"time"

	"hawx.me/code/xesende"
)

func sendMessage(c *xesende.Client, to, body string) (id string, err error) {
	account := c.Account("EX00000")

	messages, err := account.Send(xesende.Messages{
		{To: "4498499", Body: "Hey"},
	})

	if err != nil {
		return "", err
	}

	return messages.Messages[0].ID, nil
}

func getStatus(c *xesende.Client, id string) (status string, err error) {
	message, err := c.Message(id)
	if err != nil {
		return "", err
	}

	return message.Status, nil
}

func Example() {
	client := xesende.New("user@example.com", "pass")

	messageID, err := sendMessage(client, "538734", "Hey")
	if err != nil {
		log.Fatal(err)
	}

	for i := range []int{0, 1, 2, 3, 4, 5} {
		status, err := getStatus(client, messageID)
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("message status after %d seconds: %s\n", i, status)
		time.Sleep(time.Second)
	}
}
