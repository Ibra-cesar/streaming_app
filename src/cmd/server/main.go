package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Ibra-cesar/video-streaming/src/helper"
	"github.com/Ibra-cesar/video-streaming/src/internal/routes"
	"github.com/Ibra-cesar/video-streaming/src/middleware"
	"github.com/joho/godotenv"
)

func getPort() string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Failed to load .env")
	}

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("Port is Missing")
	}

	return port
}
func main() {
	//new MULTIPLEXER
	router := http.NewServeMux()

	mChain := helper.ChainMiddleware(
		middleware.Loggers,
		middleware.FakeMiddleware,
	)

	routes.RegisterRoutes(router)

	server := http.Server{
		Addr:    ":" + getPort(),
		Handler: mChain(router),
	}
	//serve the server
	fmt.Println("Server is running on: ", server.Addr)
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal("Failed to serve server")
	}
}
