package main

import (
	"RMS/database"
	"RMS/server"
	"fmt"
	"log"
	"os"
)

func main() {
	host := os.Getenv("host")
	port := os.Getenv("port")
	databaseName := os.Getenv("databaseName")
	user := os.Getenv("user")
	password := os.Getenv("password")

	err := database.ConnectAndMigrate(host, port, databaseName, user, password, database.SSLModeDisable)
	if err != nil {
		log.Printf("ConnectAndMigrate: error is:%v", err)
		return
	}
	fmt.Println("connected")
	srv := server.SetupRoutes()
	err = srv.Run(":8080")
	if err != nil {
		log.Printf("could not run the server:%v", err)
		return
	}
}
