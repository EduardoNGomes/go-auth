package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"gitbhub.com/eduardongomes/go-auth/internal/cache"
	r "gitbhub.com/eduardongomes/go-auth/internal/routes"

	"gitbhub.com/eduardongomes/go-auth/internal/providers"
	_ "github.com/joho/godotenv/autoload"
)

func main() {

	redisCache := cache.RedisConect()

	defer redisCache.Close()

	if err := validateEnvs(); err != nil {
		log.Fatal(err)
	}

	options, err := providers.NewOAuthOptions()
	if err != nil {
		log.Fatal(err)
	}

	server, err := r.NewServer(options, redisCache)
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

func validateEnvs() error {
	if os.Getenv("SECRET") == "" {
		return fmt.Errorf("Missing env SECRET")
	}

	if os.Getenv("PORT") == "" {
		return fmt.Errorf("Missing env PORT")
	}

	if os.Getenv("REDIRECT_URL") == "" {
		return fmt.Errorf("Missing env REDIRECT_URL")
	}

	if os.Getenv("REDIS_ADDR") == "" {
		return fmt.Errorf("Missing env REDIRECT_URL")
	}

	return nil
}
