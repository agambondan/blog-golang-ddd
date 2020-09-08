package main

import (
	sv "Repository-Pattern/api/server"
	"Repository-Pattern/helper"
	"fmt"
	"github.com/joho/godotenv"
	"os"
)

var server = sv.Server{}

func main() {
	err := godotenv.Load()
	helper.FailOnError(err, "error from load env")
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		os.Getenv("Hostname"), os.Getenv("PortPostgres"),
		os.Getenv("Username"), os.Getenv("PasswordDB"),
		os.Getenv("DatabaseName"))
	server.Initialize(os.Getenv("Database"), psqlInfo)
	server.Run(os.Getenv("PORT"))
}
