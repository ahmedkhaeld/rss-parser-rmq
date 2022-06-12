package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/streadway/amqp"
	"log"
	"net/http"
	"os"
)

var channelAmqp *amqp.Channel

func init() {
	// set up a TCP connection to the RabbitMQ server
	amqpConnection, err := amqp.Dial(os.Getenv(
		"RABBITMQ_URI"))
	if err != nil {
		log.Fatal(err)
	}
	channelAmqp, _ = amqpConnection.Channel()
}

type Request struct {
	URL string `json:"url"`
}

func ParserHandler(c *gin.Context) {
	var request Request
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	data, _ := json.Marshal(request)

	err := channelAmqp.Publish(
		"",
		os.Getenv("RABBITMQ_QUEUE"),
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        []byte(data),
		})
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while publishing to, RabbitMQ"})
		return
	}
	c.JSON(http.StatusOK, map[string]string{
		"message": "success"})
}
func main() {
	router := gin.Default()
	router.POST("/parse", ParserHandler)

	router.Run(":6000")
}
