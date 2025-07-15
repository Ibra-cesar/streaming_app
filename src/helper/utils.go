package helper

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Ibra-cesar/video-streaming/src/internal/routes"
	"github.com/Ibra-cesar/video-streaming/src/middleware"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
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


func DbConnection(ctx context.Context) (*pgx.Conn, error) {
	dbUri := env("DB_URI")
	conn, err := pgx.Connect(ctx, dbUri)
	if err != nil {
		return nil, fmt.Errorf("Failed to connect to database: %w", err)
	}

	return conn, nil
}

func ServerInitialization(mux *http.ServeMux) {
	//middleware CHAIN
	mChain := ChainMiddleware(
		middleware.Loggers,
		middleware.FakeMiddleware,
	)
	routes.RegisterRoutes(mux)
	//Server
	server := http.Server{
		Addr:    ":" + env("PORT"),
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
func Migrator(migPath string) error {
	fmt.Println("Initializing migrations")

	uri := env("DB_URI")

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
func env(name string) string {
	if name == "" {
		log.Fatal("Missing env variable name")
	}
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Failed to Load Env")
	}

	env := os.Getenv(name)

	return  env
}
