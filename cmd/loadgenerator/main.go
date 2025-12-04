package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gaurav/otel-observability/internal/config"
	"github.com/gaurav/otel-observability/internal/logging"
	"go.uber.org/zap"
)

const serviceName = "loadgenerator"

type LoadGenerator struct {
	logger          *logging.Logger
	frontendURL     string
	cartURL         string
	productURL      string
	checkoutURL     string
	requestsPerMin  int
	client          *http.Client
}

func main() {
	cfg, err := config.Load("configs/config.yaml")
	if err != nil {
		fmt.Printf("Failed to load config: %v\n", err)
		os.Exit(1)
	}

	logger, err := logging.NewLogger(logging.LoggerConfig{
		Level:       cfg.Logging.Level,
		Format:      cfg.Logging.Format,
		ServiceName: serviceName,
	})
	if err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()

	lg := &LoadGenerator{
		logger:         logger,
		frontendURL:    "http://localhost:8080",
		cartURL:        "http://localhost:8081",
		productURL:     "http://localhost:8082",
		checkoutURL:    "http://localhost:8083",
		requestsPerMin: 30,
		client:         &http.Client{Timeout: 10 * time.Second},
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger.Info("Load generator starting",
		zap.Int("requests_per_minute", lg.requestsPerMin),
	)

	go lg.generateLoad(ctx)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	logger.Info("Shutting down load generator...")
}

func (lg *LoadGenerator) generateLoad(ctx context.Context) {
	interval := time.Minute / time.Duration(lg.requestsPerMin)
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	requestCount := 0

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			requestCount++
			go lg.simulateUserSession(requestCount)
		}
	}
}

func (lg *LoadGenerator) simulateUserSession(sessionID int) {
	userID := fmt.Sprintf("user-%d", rand.Intn(100))
	
	lg.logger.Info("Starting user session",
		zap.Int("session_id", sessionID),
		zap.String("user_id", userID),
	)

	// 1. Browse homepage
	lg.makeRequest("GET", lg.frontendURL+"/", nil)
	time.Sleep(time.Duration(500+rand.Intn(1000)) * time.Millisecond)

	// 2. View products
	lg.makeRequest("GET", lg.productURL+"/products", nil)
	time.Sleep(time.Duration(300+rand.Intn(700)) * time.Millisecond)

	// 3. View specific product
	productIDs := []string{"OLJCESPC7Z", "66VCHSJNUP", "1YMWWN1N4O", "L9ECAV7KIM", "2ZYFJ3GM2N"}
	productID := productIDs[rand.Intn(len(productIDs))]
	lg.makeRequest("GET", lg.productURL+"/product/"+productID, nil)
	time.Sleep(time.Duration(400+rand.Intn(800)) * time.Millisecond)

	// 4. Add to cart
	addToCartReq := map[string]interface{}{
		"user_id":    userID,
		"product_id": productID,
		"quantity":   rand.Intn(3) + 1,
	}
	lg.makeRequest("POST", lg.cartURL+"/cart/add", addToCartReq)
	time.Sleep(time.Duration(200+rand.Intn(500)) * time.Millisecond)

	// 5. View cart
	lg.makeRequest("GET", lg.cartURL+"/cart/?user_id="+userID, nil)
	time.Sleep(time.Duration(300+rand.Intn(600)) * time.Millisecond)

	// 6. Checkout (70% of users complete checkout)
	if rand.Float64() < 0.7 {
		checkoutReq := map[string]interface{}{
			"user_id": userID,
			"email":   fmt.Sprintf("%s@example.com", userID),
			"address": map[string]string{
				"street":  "123 Main St",
				"city":    "San Francisco",
				"state":   "CA",
				"zip":     "94102",
				"country": "USA",
			},
			"credit_card": map[string]string{
				"number": "4111111111111111",
				"cvv":    "123",
				"expiry": "12/25",
			},
		}
		lg.makeRequest("POST", lg.checkoutURL+"/checkout", checkoutReq)
	}

	lg.logger.Info("User session completed",
		zap.Int("session_id", sessionID),
		zap.String("user_id", userID),
	)
}

func (lg *LoadGenerator) makeRequest(method, url string, body interface{}) {
	var req *http.Request
	var err error

	if body != nil {
		jsonBody, _ := json.Marshal(body)
		req, err = http.NewRequest(method, url, bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
	} else {
		req, err = http.NewRequest(method, url, nil)
	}

	if err != nil {
		lg.logger.Error("Failed to create request", zap.Error(err))
		return
	}

	resp, err := lg.client.Do(req)
	if err != nil {
		lg.logger.Debug("Request failed",
			zap.String("method", method),
			zap.String("url", url),
			zap.Error(err),
		)
		return
	}
	defer resp.Body.Close()

	lg.logger.Debug("Request completed",
		zap.String("method", method),
		zap.String("url", url),
		zap.Int("status", resp.StatusCode),
	)
}

