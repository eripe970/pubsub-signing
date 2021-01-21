package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	signing "github.com/eripe970/pubsub-signing"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func main() {

	secret := os.Getenv("SECRET")
	projectId := os.Getenv("PROJECT_ID")
	topic := os.Getenv("TOPIC")

	r := gin.New()
	r.Use(gin.Recovery())

	r.GET("/change-name", func(c *gin.Context) {

		err := PublishEvent(projectId, topic, secret, UserChangedNameEvent{Name: "new-name"})

		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		c.Status(http.StatusOK)
	})

	r.POST("/events", func(c *gin.Context) {
		payload, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		message, err := handleMessage(payload, secret)

		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		log.Print(message)

		c.Status(http.StatusNoContent)

	})

	r.Run()
}

func handleMessage(payload []byte, secret string) (string, error) {
	request, err := signing.ConstructMessage(payload, secret)

	if err != nil {
		return "", err
	}

	decoded, err := base64.StdEncoding.DecodeString(request.Message.Data)

	if err != nil {
		return "", err
	}

	eventType := request.Message.Attributes["event-type"]

	switch eventType {
	case "user.name.changed":
		var event UserChangedNameEvent

		err = json.Unmarshal(decoded, &event)

		if err != nil {
			return "", err
		}

		return fmt.Sprintf("processed event %v", event.Type()), nil
	}

	return "", nil
}
