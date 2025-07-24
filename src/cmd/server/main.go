package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/Ibra-cesar/video-streaming/src/helper"
	"github.com/Ibra-cesar/video-streaming/src/internal/query_repo"
)

func main() {
	//DB CONNECTION
	ctx := context.Background()

	app, err := helper.App(ctx)
	if err != nil {
		log.Fatal("Failed to initialize Application, %w", err)
	}

	defer app.Pool.Close()

	err = helper.Migrator()
	if err != nil {
		log.Fatalf("Migrations Failed, %v", err)
	}
	fmt.Println("Migrations is success")

	repo := query_repo.New(app.Pool)

	userList, err := repo.GetAllPlayers(ctx)
	if err != nil {
		log.Fatal("Failed to get Users")
	}
	helper.PrintUser(userList)
	//SERVER SETUP
	routes := http.NewServeMux()
	helper.ServerInitialization(routes, app.AuthHandler, app.RefreshHandler)
}
