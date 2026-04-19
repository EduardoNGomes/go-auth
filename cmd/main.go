package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	r "gitbhub.com/eduardongomes/go-auth/internal/routes"

	"gitbhub.com/eduardongomes/go-auth/internal/providers"
	_ "github.com/joho/godotenv/autoload"
)

func main() {

	googleProvider := providers.NewGoogle()
	options := map[providers.Provider]providers.Actions{providers.GOOGLE: googleProvider}
	server, err := r.NewServer(options)

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
