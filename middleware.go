package main

import (
	"bufio"
	"github.com/gin-gonic/gin"
	"os"
)

func cors(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "*")
	c.Header("Access-Control-Allow-Headers", "*")
	c.Next()
}

var API_KEYS map[string]bool

func authorization(c *gin.Context) {
	if API_KEYS == nil {
		API_KEYS = make(map[string]bool)
		if _, err := os.Stat("config/api_keys.txt"); err == nil {
			file, _ := os.Open("config/api_keys.txt")
			defer file.Close()
			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				key := scanner.Text()
				if key != "" {
					API_KEYS["Bearer "+key] = true
				}
			}
		}
	}
	if len(API_KEYS) != 0 && !API_KEYS[c.Request.Header.Get("Authorization")] {
		if c.Request.Header.Get("Authorization") == "" {
			c.JSON(401, gin.H{"error": "No API key provided."})
		} else {
			c.JSON(401, gin.H{"error": "Invalid API key " + c.Request.Header.Get("Authorization")})
		}
		c.Abort()
		return
	}
	c.Next()
}
