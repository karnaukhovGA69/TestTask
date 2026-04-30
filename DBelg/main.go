package DBelg

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"go.uber.org/zap"
)

var longToShort map[string]string
var shortToLong map[string]string
var mutex = &sync.RWMutex{}

func main() {
	mux := http.NewServeMux()
	longToShort = make(map[string]string)
	shortToLong = make(map[string]string)
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalln("ERROR:Не получилось создать логгер")
		return
	}
	defer logger.Sync()
	fileShortToLong, err := os.ReadFile("shortToLong.json")
	if err != nil {
		logger.Warn("Не получилось прочитать shortToLong.json", zap.Error(err))
	} else {
		err = json.Unmarshal(fileShortToLong, &shortToLong)
		if err != nil {
			logger.Warn("Не получилось распарсить shortToLong.json", zap.Error(err))
		}
	}
	fileLongToShort, err := os.ReadFile("longToShort.json")
	if err != nil {
		logger.Warn("Не получилось прочитать longToShort.json", zap.Error(err))
	} else {
		err = json.Unmarshal(fileLongToShort, &longToShort)
		if err != nil {
			logger.Warn("Не получилось распарсить longToShort.json", zap.Error(err))
		}
	}
	mux.HandleFunc("/", handler)
	server := &http.Server{
		Addr:         ":7000",
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  time.Minute,
	}
	if err := server.ListenAndServe(); err != nil {
		logger.Fatal("Сервер DBelg не запустился: ", zap.Error(err))
	}
}

func handler(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		GetHandler(rw, r)
	case http.MethodPost:
		PostHandler(rw, r)
	default:
		http.Error(rw, "Этот метод не поддерживается", http.StatusMethodNotAllowed)
	}
}

func GetHandler(rw http.ResponseWriter, r *http.Request) {

}

func PostHandler(rw http.ResponseWriter, r *http.Request) {

}
