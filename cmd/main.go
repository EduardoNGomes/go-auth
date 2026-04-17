package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	r "gitbhub.com/eduardongomes/go-auth/internal/routes"

	g "gitbhub.com/eduardongomes/go-auth/internal/google"
	_ "github.com/joho/godotenv/autoload"
)

func main() {

	googleInterface := g.NewGoogle()
	server, err := r.NewServer(googleInterface)

	if err != nil {
		log.Fatal(err)
	}

	routes, err := r.NewRoutes(server)

	if err != nil {
		log.Fatal(err)
	}

	port := os.Getenv("PORT")

	fmt.Printf("Server is running o port %s\n", port)

	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), routes); err != nil {
		log.Fatalf("Fail to run port 5000 %v", err)
	}

}
