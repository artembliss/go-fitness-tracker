package main

import (
	"fmt"
	"log"
	"os"

	"github.com/artembliss/go-fitness-tracker/logger/sl"
	"github.com/artembliss/go-fitness-tracker/storage/postgre"
	"github.com/joho/godotenv"
)

func main() {
	storage, err := postgre.New()
	if err != nil{
		sl.Err(err)
		fmt.Println(err)
		os.Exit(1)
	}
	_ = storage
	fmt.Println("successfully connected!")
}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Ошибка загрузки .env: %s", err)
	}
}