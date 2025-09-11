package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.POST("/stub/payments", func(c *gin.Context) {
		log.Println("💸 External system recieved payment request")
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	log.Println("Stub external system running on :9000")
	r.Run(":9000")
}
