package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	err := router.Run()
	if err != nil {
		log.Fatal("unable to run server:", err)
	}
}
 