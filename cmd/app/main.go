package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"main/internal/handler"
	"main/internal/service"
	"main/internal/storage"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func main() {
	_ = godotenv.Load()
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalln("ERROR:Не получилось создать логгер")
	}
	defer logger.Sync()
	args := os.Args
	if len(args) < 2 {
		fmt.Println("Не хватает аргументвов - Шаблон: go run . postgres/dbelg")
		logger.Fatal("Не ввели аргумент Базы Данных")
	}
	typeOfDB := args[1]
	db, err := storage.MakeDB(typeOfDB)
	if err != nil {
		logger.Fatal("Не получислоь запустить базу данных", zap.Error(err))
	}
	if closer, ok := db.(interface{ Close() error }); ok {
		defer closer.Close()
	}

	srvc := service.NewService(db)
	hndlr := handler.NewHandler(srvc)
	mux := http.NewServeMux()

	mux.HandleFunc("GET /{shortURL}", hndlr.GetHandler)
	mux.HandleFunc("POST /url", hndlr.PostHandler)
	server := http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  time.Minute,
	}
	logger.Info("Начинаем запускать сервер на порту 8080")
	serverErr := make(chan error, 1)
	go func() {
		err := server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErr <- err
		}
		close(serverErr)
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErr:
		if err != nil {
			logger.Fatal("Не смог запуститься сервер на порту 8080", zap.Error(err))
		}
	case <-stop:
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			logger.Fatal("Не получилось остановить сервер", zap.Error(err))
		}
	}

}
