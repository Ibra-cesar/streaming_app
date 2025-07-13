package main

import (
	"fmt"
	"log"
	"os"
	"github.com/joho/godotenv"
)

func main(){
	fmt.Println("Hello Project")

	err := godotenv.Load()
	if err != nil{
		log.Fatal("Failed to load .env")
	}

	port := os.Getenv("PORT")
	if port == ""{
		log.Fatal("Port is missing!")
	}

	fmt.Println("Port is running on: ", port)
}
