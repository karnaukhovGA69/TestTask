package main

import (
	"fmt"
	"log"
	"os"

	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalln("ERROR:Не получилось создать логгер")
	}
	args := os.Args
	if len(args) < 2 {
		fmt.Println("Не хватает аргументвов - Шаблон: go run . postgres/DBelg")
		logger.Fatal("Не ввели аргумент Базы Данных")
		os.Exit(1)
	}

	//database := MakeDB(args[1])

}
