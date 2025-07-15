package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/Ibra-cesar/video-streaming/src/helper"
)
func main() {
	//DB CONNECTION
	ctx := context.Background()
 
	conn, err := helper.DbConnection(ctx)
	if err != nil {
		log.Fatalf("Failed To Set DB Connection: %v", err)
	}

	defer conn.Close(ctx)
	
	fmt.Println("Successfully connected to DataBase")

	//Migration SETUP
	migPath := "file://./src/helper/db/migrations"

	err = helper.Migrator(migPath)
	if err != nil {
		log.Fatalf("Migrations Failed, %v", err)
	}
	fmt.Println("Migrations is success")

	//SERVER SETUP
	routes := http.NewServeMux()
	helper.ServerInitialization(routes)
}
