package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

const (
	baseURL     = "http://localhost:8080"
	concurrency = 100
	totalReqs   = 1000
)

type LoginReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func main() {
	// 1. Register a test user
	username := fmt.Sprintf("perf_user_%d", time.Now().Unix())
	password := "Password123"
	registerPayload := map[string]string{
		"username": username,
		"password": password,
		"email":    username + "@example.com",
		"phone":    fmt.Sprintf("139%d", time.Now().Unix()%100000000),
	}
	regBody, _ := json.Marshal(registerPayload)
	resp, err := http.Post(baseURL+"/auth/register", "application/json", bytes.NewBuffer(regBody))
	if err != nil {
		fmt.Printf("Failed to register: %v\n", err)
		return
	}
	resp.Body.Close()
	fmt.Printf("Registered user: %s\n", username)

	// 2. Concurrent Login
	var wg sync.WaitGroup
	start := time.Now()
	successCount := 0
	var mu sync.Mutex

	sem := make(chan struct{}, concurrency)

	for i := 0; i < totalReqs; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			sem <- struct{}{}        // Acquire token
			defer func() { <-sem }() // Release token

			loginReq := LoginReq{
				Username: username,
				Password: password,
			}
			body, _ := json.Marshal(loginReq)
			
			startReq := time.Now()
			resp, err := http.Post(baseURL+"/auth/login", "application/json", bytes.NewBuffer(body))
			if err != nil {
				fmt.Printf("Request failed: %v\n", err)
				return
			}
			defer resp.Body.Close()
			
			if resp.StatusCode == 200 {
				mu.Lock()
				successCount++
				mu.Unlock()
			}
			
			// Optional: print latency for debugging
			// fmt.Printf("Req latency: %v\n", time.Since(startReq))
			_ = startReq
		}()
	}

	wg.Wait()
	duration := time.Since(start)
	
	fmt.Printf("\nPerformance Test Results:\n")
	fmt.Printf("Total Requests: %d\n", totalReqs)
	fmt.Printf("Concurrency: %d\n", concurrency)
	fmt.Printf("Success Count: %d\n", successCount)
	fmt.Printf("Total Duration: %v\n", duration)
	fmt.Printf("QPS: %.2f\n", float64(totalReqs)/duration.Seconds())
	fmt.Printf("Avg Latency: %v\n", duration/time.Duration(totalReqs))
}
