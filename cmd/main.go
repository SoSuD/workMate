package main

import (
	"github.com/joho/godotenv"
	"log"
	"workMate/internal/apiserver"
)

//TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
	serverConfig := apiserver.NewConfig()
	apiserver.Start(serverConfig)
}
