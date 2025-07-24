package helper

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Ibra-cesar/video-streaming/src/internal/handlers"
	"github.com/Ibra-cesar/video-streaming/src/internal/routes"
	"github.com/Ibra-cesar/video-streaming/src/middleware"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"     // Add this line
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)
// Middleware Chaining Helpers
type Middleware func(http.Handler) http.Handler

func ChainMiddleware(ms ...Middleware) Middleware {
	return func(next http.Handler) http.Handler {
	for i := len(ms) - 1; i >= 0; i-- {
			x := ms[i]
			next = x(next)
		}
		return next
	}
}

//DB CONNECTION
func DbConnection(ctx context.Context) (*pgxpool.Pool, error) {
	config := DefaultConfig()

	pool, err := ConnPool(ctx, config)
		if err != nil {
		return nil, fmt.Errorf("Failed to connect to database: %w", err)
	}
	return  pool, nil
}

//server
func ServerInitialization(mux *http.ServeMux, authHandlers *handlers.AuthConnServices, refreshHandler *handlers.RefreshHandlers) {
	//middleware CHAIN
	mChain := ChainMiddleware(
		middleware.Loggers,
		middleware.FakeMiddleware,
	)
	routes.RegisterPubRoutes(mux, authHandlers, refreshHandler)
	routes.RegisterPrivRoutes(mux)
	//Server
	server := http.Server{
		Addr:    ":" + Env("PORT"),
		Handler: mChain(mux),
	}
	//serve the server
	fmt.Println("Server is running on: ", server.Addr)
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal("Failed to serve server")
	}
}

//DB MIGRATOR
func Migrator() error {
	fmt.Println("Initializing migrations")

	migPath := "file://./src/helper/db/migrations/"
	uri := Env("DB_URI")

	m, err := migrate.New(
		migPath,
		uri,
	)

	if err != nil{
		return fmt.Errorf("Error while migrate, %v", err)
	}

	defer func(){
		srcErr, dbErr := m.Close()
		if srcErr != nil{
			log.Printf("Source Error while closing, %v", srcErr)
		}
		if dbErr != nil {
			log.Printf("DataBase Error while closing, %v", dbErr)
		}
	}()

	//Backup clean state
	//if err := m.Force(1); err != nil { // IMPORTANT: Use the correct dirty version number
  //      return fmt.Errorf("Failed to force migration version: %v", err)
  //  }
  //  fmt.Println("Successfully forced migration version 1 to clean state.")

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("Err while migrating db, %v", err)
	}else if err == migrate.ErrNoChange{
		fmt.Println("No migration changes is applied")
	}else{
		fmt.Println("Successfully migrate db")
	}

	return nil
}

//ENV
func Env(name string) string {
	if name == "" {
		log.Fatal("Missing env variable name")
	}
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Failed to Load Env")
	}
	env := os.Getenv(name)
	return env
}

