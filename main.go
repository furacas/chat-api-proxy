package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.Use(cors)

	router.GET("/ping", ping)

	router.OPTIONS("/v1/chat/completions", optionsHandler)
	router.POST("/v1/chat/completions", authorization, completionsHandler)

	router.Run("0.0.0.0:8080")
}
