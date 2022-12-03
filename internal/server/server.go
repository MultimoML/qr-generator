package server

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/julienschmidt/httprouter"

	"github.com/multimoml/qr-generator/internal/config"
	"github.com/multimoml/qr-generator/internal/model"
)

func Run(_ context.Context) {
	// Load environment variables
	config.Environment()

	// Start HTTP server
	router := httprouter.New()

	// Endpoints
	router.GET("/qr", QR)

	log.Fatal(http.ListenAndServe(":6002", router))
}

func QR(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	qrApi := "https://api.qrserver.com/v1/create-qr-code/?size=500x500&data="

	numProducts := 0
	var products []model.Product

	// Call Dispatcher microservice to get number of products in DB
	res, err := http.Get("http://dispatcher:6001/products")
	if err != nil {
		log.Println(err)
	}

	// Decode JSON response
	if err = json.NewDecoder(res.Body).Decode(&products); err != nil {
		log.Println(err)
	}
	numProducts = len(products)

	// Call QR API to generate QR code
	data := fmt.Sprintf("Number of products in database: %d", numProducts)
	log.Println(data)

	res, err = http.Get(qrApi + url.QueryEscape(data))
	if err != nil {
		log.Println(err)
	}

	w.Header().Set("Content-Type", "image/png")
	w.WriteHeader(http.StatusOK)
	if _, err := io.Copy(w, res.Body); err != nil {
		log.Println(err)
	}
}
