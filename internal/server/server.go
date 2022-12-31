package server

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/multimoml/qr-generator/docs"
	"github.com/multimoml/qr-generator/internal/config"
)

var sConfig *config.Config

func Run(_ context.Context) {
	// Load environment variables
	sConfig = config.LoadConfig()

	// Set up router
	router := gin.Default()

	// Endpoints
	q := router.Group("/qr")
	{
		q.GET("/live", Liveness)
		q.GET("/ready", Readiness)
		q.GET("/v1/:id", Generate)

		// Redirect /openapi to /openapi/index.html
		q.GET("/", func(c *gin.Context) {
			c.Redirect(http.StatusMovedPermanently, "/qr/openapi/index.html")
		})
		q.GET("/openapi", func(c *gin.Context) {
			c.Redirect(http.StatusMovedPermanently, "/qr/openapi/index.html")
		})
		q.GET("/openapi/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	// Start HTTP server
	log.Fatal(router.Run(fmt.Sprintf(":%s", sConfig.Port)))
}

// Liveness is a simple endpoint to check if the server is alive
// @Summary Get liveness status of the microservice
// @Description Get liveness status of the microservice
// @Tags Kubernetes
// @Success 200 {string} string
// @Router /qr/live [get]
func Liveness(c *gin.Context) {
	// Check if the config value 'broken' is set to 1
	res, err := http.Get(sConfig.ConfigServer + "/broken")

	// If the config server is not available, treat it as not broken (default to alive)
	if err == nil {
		value, err := io.ReadAll(res.Body)
		if err == nil && string(value) == "1" {
			c.String(http.StatusServiceUnavailable, "dead")
			return
		}
	}

	c.String(http.StatusOK, "alive")
}

// Readiness is a simple endpoint to check if the server is ready
// @Summary Get readiness status of the microservice
// @Description Get readiness status of the microservice
// @Tags Kubernetes
// @Success 200 {string} string
// @Failure 503 {string} string
// @Router /qr/ready [get]
func Readiness(c *gin.Context) {
	qrApi := "https://api.qrserver.com/v1/create-qr-code/?size=10x10&data=1"
	dispatcher := "http://dispatcher:6001"
	iAmReady := true

	// If using dev environment access local dispatcher
	if os.Getenv("ACTIVE_ENV") == "dev" {
		dispatcher = "http://localhost:6001"
	}

	// Call QR API to check if it's ready
	if res, err := http.Get(qrApi); err != nil || res.StatusCode != http.StatusOK {
		iAmReady = false
		log.Println("QR API is not ready: ", err)
	}

	// Call Dispatcher microservice to check if it's ready
	if res, err := http.Get(dispatcher + "/products/ready"); err != nil || res.StatusCode != http.StatusOK {
		iAmReady = false
		log.Println("Dispatcher microservice is not ready: ", err)
	}

	if iAmReady {
		c.String(http.StatusOK, "ready")
	} else {
		c.String(http.StatusServiceUnavailable, "not ready")
	}
}

// Generate generates a QR code
// @Summary Generate a QR code
// @Description Generate a QR code
// @Produce  image/png
// @Param id path string true "Product ID"
// @Success 200 {string} string
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Tags Service
// @Router /qr/v1/{id} [get]
func Generate(c *gin.Context) {
	id := c.Param("id")
	qrApi := "https://api.qrserver.com/v1/create-qr-code/?size=256x256&data="
	dispatcher := "http://dispatcher:6001"

	// If using dev environment access local dispatcher
	if os.Getenv("ACTIVE_ENV") == "dev" {
		dispatcher = "http://localhost:6001"
	}

	// Get product with the given ID from Dispatcher
	res, err := http.Get(dispatcher + "/products/v1/" + id)
	if err != nil {
		log.Println(err)
		c.String(http.StatusNotFound, "Failed to get product: 1")
		return
	}

	// Read the response body
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
		c.String(http.StatusNotFound, "Failed to get product: 2")
		return
	}

	// Call QR API to generate QR code
	res, err = http.Get(qrApi + url.PathEscape(string(body)))
	if err != nil {
		log.Println(err)
		c.String(http.StatusInternalServerError, "Failed to generate QR code: 1")
		return
	}

	c.Header("Content-Type", "image/png")
	c.Status(http.StatusOK)
	if _, err = io.Copy(c.Writer, res.Body); err != nil {
		log.Println(err)
		c.String(http.StatusInternalServerError, "Failed to generate QR code: 2")
	}
}
