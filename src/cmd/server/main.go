package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/Ibra-cesar/video-streaming/src/helper"
	"github.com/Ibra-cesar/video-streaming/src/internal/query_repo"
	"github.com/google/uuid"
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
	migPath := "file://./src/helper/db/migrations/"

	err = helper.Migrator(migPath)
	if err != nil {
		log.Fatalf("Migrations Failed, %v", err)
	}
	fmt.Println("Migrations is success")

  repo := query_repo.New(conn)

	newUserId := uuid.New()
	newUsr, err := repo.InsertUser(ctx, query_repo.InsertUserParams{
		ID: newUserId,
		Name: "Ibrahim",
		Email: "db_test@exam.ple",
		PasswordHash: "IbrahimHash",
		Salt: "Ibrahim Salt",
	})
	if err != nil {
		log.Fatal("Failed to insert new User")
	}
	fmt.Printf("New User: %v", newUsr)

	userList, err := repo.GetAllPlayers(ctx)
	if err != nil {
		log.Fatal("Failed to get Users")
	}
	fmt.Printf("Users: %v", userList)
	//SERVER SETUP
	routes := http.NewServeMux()
	helper.ServerInitialization(routes)
}
