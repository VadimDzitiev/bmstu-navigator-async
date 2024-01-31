package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

const secretKey = "secret-async-key"

type service struct {
	ID             int64  `json:"id"`
	TransitionTime int64  `json:"transition_time"`
	SecretKey      string `json:"secret_key"`
}

func main() {
	r := gin.Default()

	r.POST("/set_time", func(c *gin.Context) {
		var service service

		if err := c.ShouldBindJSON(&service); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		go func() {
			time.Sleep(4 * time.Second)
			SendStatus(service)
		}()

		c.JSON(http.StatusOK, gin.H{"message": "Status update initiated"})
	})

	r.Run(":8080")
}

func SendStatus(service service) bool {
	// fmt.Println(service.TransitionTime)
	service.TransitionTime = generateRandomStatusRefer()
	// fmt.Println(service.TransitionTime)
	service.SecretKey = secretKey
	response, err := performPUTRequest("http://localhost:8000/api/time/" + fmt.Sprint(service.ID) + "/put/", service)
	if err != nil {
		fmt.Println("Error sending status:", err)
		return false
	}

	if response.StatusCode == http.StatusOK {
		fmt.Println("Status sent successfully for pk:", service.ID)
		return true
	} else {
		fmt.Println("Failed to process PUT request")
		return false
	}
}

func generateRandomStatusRefer() int64 {
	return rand.Int63n(15) + 1
}

func performPUTRequest(url string, data service) (*http.Response, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", secretKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	return resp, nil
}
