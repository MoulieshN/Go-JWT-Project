package main

import (
	"context"
	"os"

	config "github.com/MoulieshN/Go-JWT-Project.git/config"
	server "github.com/MoulieshN/Go-JWT-Project.git/server"
	"github.com/namsral/flag"
)

var (
	env = flag.String("env", ".env", "Path to .env file")
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		port = "8123"
	}

	config.Init(*env)

	server.Init(context.Background(), port)
}
